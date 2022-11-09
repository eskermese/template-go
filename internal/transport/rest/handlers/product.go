package restHandler

import (
	"net/http"

	"github.com/eskermese/template-go/internal/core"
	"github.com/gin-gonic/gin"
)

func (h *Handler) initProductRoutes(api *gin.RouterGroup) {
	products := api.Group("/products")
	{
		products.GET("", h.getAllProducts)
		products.POST("", h.createProduct)
		products.PUT("/:id", h.updateProduct)
		products.DELETE("/:id", h.deleteProduct)
	}
}

// @Summary Get Products
// @Tags products
// @Description Get products
// @ModuleID getAllProducts
// @Accept  json
// @Produce  json
// @Success 200 {object} []core.Product
// @Failure 400 {object} core.Errors
// @Router /products [get]
func (h *Handler) getAllProducts(c *gin.Context) {
	products, err := h.productService.GetAll(c.Request.Context())
	if err != nil {
		h.newResponse(c, err)

		return
	}

	c.JSON(http.StatusOK, products)
}

// @Summary Create Product
// @Tags products
// @Description Create products
// @ModuleID createProduct
// @Accept  json
// @Produce  json
// @Param input body core.CreateProductInput true "create product"
// @Success 201 {object} core.Product
// @Failure 400 {object} core.Errors
// @Router /products [post]
func (h *Handler) createProduct(c *gin.Context) {
	var inp core.CreateProductInput
	if err := c.ShouldBindJSON(&inp); err != nil {
		h.newResponse(c, err)

		return
	}

	if validationError := inp.Validate(); validationError != nil {
		h.newResponse(c, validationError)

		return
	}

	product, err := h.productService.Create(c.Request.Context(), inp)
	if err != nil {
		h.newResponse(c, err)

		return
	}

	c.JSON(http.StatusOK, product)
}

// @Summary Update Product
// @Tags products
// @Description Update products
// @ModuleID updateProduct
// @Accept  json
// @Produce  json
// @Param id path string true "product id"
// @Param input body core.UpdateProductInput true "update product"
// @Success 204
// @Failure 400 {object} core.Errors
// @Failure 404 {object} response
// @Router /products/{id} [put]
func (h *Handler) updateProduct(c *gin.Context) {
	var (
		inp core.UpdateProductInput
		err error
	)

	if err = c.ShouldBindJSON(&inp); err != nil {
		h.newResponse(c, err)

		return
	}

	if err = inp.Validate(); err != nil {
		h.newResponse(c, err)

		return
	}

	productID, err := paramInt(c, "id")
	if err != nil {
		h.newResponse(c, err)

		return
	}

	if _, err = h.productService.GetByID(c.Request.Context(), productID); err != nil {
		h.newResponse(c, err)

		return
	}

	if err = h.productService.Update(c.Request.Context(), productID, inp); err != nil {
		h.newResponse(c, err)

		return
	}

	c.Status(http.StatusNoContent)
}

// @Summary Delete Product
// @Tags products
// @Description Delete product
// @ModuleID deleteProduct
// @Accept  json
// @Produce  json
// @Param id path string true "product id"
// @Success 204
// @Failure 404 {object} response
// @Router /products/{id} [delete]
func (h *Handler) deleteProduct(c *gin.Context) {
	productID, err := paramInt(c, "id")
	if err != nil {
		h.newResponse(c, err)

		return
	}

	if err = h.productService.Delete(c.Request.Context(), productID); err != nil {
		h.newResponse(c, err)

		return
	}

	c.Status(http.StatusNotFound)
}
