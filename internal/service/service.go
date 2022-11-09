package service

import (
	"context"

	"github.com/eskermese/template-go/internal/core"
)

type ProductStorage interface {
	GetAll(ctx context.Context) ([]core.Product, error)
	Create(ctx context.Context, product *core.Product) error
	Delete(ctx context.Context, id int) error
	Update(ctx context.Context, product core.Product) error
	GetByID(ctx context.Context, id int) (core.Product, error)
}

type Deps struct {
	ProductStorage ProductStorage
}

type Service struct {
	Product *Product
}

func New(deps Deps) *Service {
	return &Service{
		Product: NewProduct(deps.ProductStorage),
	}
}
