package main

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os/signal"
	"syscall"

	"github.com/eskermese/template-go/pkg/workers"

	"github.com/eskermese/template-go/pkg/logger"

	"github.com/eskermese/template-go/pkg/database/postgresql"

	_ "github.com/eskermese/template-go/docs"
	"github.com/eskermese/template-go/internal/transport/rest"
	restHandler "github.com/eskermese/template-go/internal/transport/rest/handlers"

	"github.com/eskermese/template-go/internal/config"
	"github.com/eskermese/template-go/internal/service"
	"github.com/eskermese/template-go/internal/storage"
)

// @title Product app REST-API
// @version 1.0
// @description Application for adding/getting products

// @host localhost:8000
// @BasePath /api/

// Run initializes whole application.
func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGTERM, syscall.SIGINT)
	defer cancel()

	log := logger.New("debug", "template-go")
	defer func() {
		if err := logger.Cleanup(log); err != nil {
			log.Error("error cleanup logs", logger.Error(err))
		}
	}()

	cfg, err := config.New()
	if err != nil {
		log.Fatal("error when initializing the config", logger.Error(err))
	}

	db, err := postgresql.NewClient(context.TODO(), postgresql.StorageConfig{
		ConnStr:     cfg.Database.DSURL,
		MaxAttempts: 5,
	})
	if err != nil {
		log.Fatal("connection to postgres error", logger.Error(err))
	}
	defer db.Close()

	storages := storage.New(db)
	services := service.New(service.Deps{
		ProductStorage: storages.Product,
	})
	handlers := restHandler.New(restHandler.Deps{
		ProductService: services.Product,
		Logger:         log,
	})

	g, gCtx := workers.GroupWithContext(ctx)

	g.Go(func() error {
		httpSrv := rest.NewServer(gCtx, cfg, handlers.InitRouter(cfg), log)
		if err := httpSrv.Run(); err != nil {
			return fmt.Errorf("http server run: %w", err)
		}

		return nil
	})

	if err := g.Wait(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		log.Error("error group wait", logger.Error(err))
	}
}
