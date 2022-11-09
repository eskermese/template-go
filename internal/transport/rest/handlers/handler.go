package restHandler

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/eskermese/template-go/docs"
	"github.com/eskermese/template-go/internal/config"
	"github.com/eskermese/template-go/internal/core"
	"github.com/eskermese/template-go/pkg/logger"
	"github.com/gin-gonic/gin"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

type ProductService interface {
	GetAll(ctx context.Context) ([]core.Product, error)
	Create(ctx context.Context, inp core.CreateProductInput) (core.Product, error)
	Delete(ctx context.Context, id int) error
	Update(ctx context.Context, id int, inp core.UpdateProductInput) error
	GetByID(ctx context.Context, id int) (core.Product, error)
}

type Deps struct {
	ProductService ProductService
	Logger         logger.Logger
}

type Handler struct {
	productService ProductService
	logger         logger.Logger
}

func New(deps Deps) *Handler {
	return &Handler{
		productService: deps.ProductService,
		logger:         deps.Logger,
	}
}

func (h *Handler) InitRouter(cfg *config.Config) *gin.Engine {
	// Init gin handler
	router := gin.Default()

	router.Use(
		gin.Recovery(),
		gin.Logger(),
	)

	docs.SwaggerInfo.Host = fmt.Sprintf("%s:%d", cfg.ClientHTTP.Host, cfg.ClientHTTP.Port)
	if cfg.Environment != config.EnvLocal {
		docs.SwaggerInfo.Host = cfg.ClientHTTP.Host
	}

	if cfg.Environment != config.Prod {
		router.GET("/swagger/*any", ginSwagger.WrapHandler(
			swaggerfiles.Handler,
			ginSwagger.URL(fmt.Sprintf("%s://%s/%s",
				cfg.ClientHTTP.Schema,
				docs.SwaggerInfo.Host,
				"swagger/doc.json",
			)),
		))
	}

	// Init router
	router.GET("/ping", func(c *gin.Context) {
		c.String(http.StatusOK, "pong")
	})

	h.initAPI(router)

	return router
}

func (h *Handler) initAPI(router *gin.Engine) {
	api := router.Group("/api")
	{
		h.initProductRoutes(api)
	}
}

type response struct {
	Detail string `json:"detail"`
}

func (h *Handler) newResponse(c *gin.Context, err error) {
	var ve core.Errors
	if errors.As(err, &ve) {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"errors": ve})

		return
	}

	var statusCode int

	switch {
	case errors.Is(err, core.ErrNotFound):
		statusCode = http.StatusNotFound
	case errors.Is(err, core.ErrProductNameDuplicate):
		statusCode = http.StatusConflict
	case errors.Is(err, core.ErrInvalidType):
		statusCode = http.StatusBadRequest
	default:
		h.logger.Error("internal server error", logger.Error(err))

		statusCode = http.StatusInternalServerError
	}

	c.AbortWithStatusJSON(statusCode, response{err.Error()})
}

func paramInt(c *gin.Context, param string) (int, error) {
	var errs core.Errors

	val := c.Param(param)
	if val == "" {
		errs = append(errs, core.FieldError{Field: param, Message: "empty id param"})

		return 0, errs
	}

	id, err := strconv.Atoi(val)
	if err != nil {
		errs = append(errs, core.FieldError{Field: param, Message: "invalid id param"})

		return 0, errs
	}

	return id, nil
}
