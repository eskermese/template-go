package storage

import "github.com/eskermese/template-go/pkg/database/postgresql"

type Storage struct {
	Product *Product
}

func New(db postgresql.Client) *Storage {
	return &Storage{
		Product: NewProduct(db),
	}
}
