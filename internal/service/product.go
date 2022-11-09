package service

import (
	"context"

	"github.com/eskermese/template-go/internal/core"
)

type Product struct {
	repo ProductStorage
}

func NewProduct(repo ProductStorage) *Product {
	return &Product{
		repo: repo,
	}
}

func (s *Product) GetAll(ctx context.Context) ([]core.Product, error) {
	return s.repo.GetAll(ctx)
}

func (s *Product) Create(ctx context.Context, inp core.CreateProductInput) (core.Product, error) {
	product := core.Product{
		Name:  inp.Name,
		Price: inp.Price,
	}

	if err := s.repo.Create(ctx, &product); err != nil {
		return core.Product{}, err
	}

	return product, nil
}

func (s *Product) Delete(ctx context.Context, id int) error {
	return s.repo.Delete(ctx, id)
}

func (s *Product) Update(ctx context.Context, id int, inp core.UpdateProductInput) error {
	return s.repo.Update(ctx, core.Product{
		ID:    id,
		Name:  inp.Name,
		Price: inp.Price,
	})
}

func (s *Product) GetByID(ctx context.Context, id int) (core.Product, error) {
	return s.repo.GetByID(ctx, id)
}
