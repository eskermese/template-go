package storage

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v4"

	"github.com/eskermese/template-go/internal/core"
	"github.com/eskermese/template-go/pkg/database/postgresql"
)

var ErrProductNameUnique = "ERROR: duplicate key value violates unique constraint \"products_name_key\" (SQLSTATE 23505)"

type Product struct {
	db postgresql.Client
}

func NewProduct(db postgresql.Client) *Product {
	return &Product{
		db: db,
	}
}

func (r *Product) GetAll(ctx context.Context) ([]core.Product, error) {
	q := "SELECT id, name, price FROM products"

	products := make([]core.Product, 0)

	rows, err := r.db.Query(ctx, q)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		var product core.Product

		err = rows.Scan(&product.ID, &product.Name, &product.Price)
		if err != nil {
			return nil, err
		}

		products = append(products, product)
	}

	return products, nil
}

func (r *Product) Create(ctx context.Context, product *core.Product) error {
	q := "INSERT INTO products (name, price) VALUES ($1, $2) RETURNING id"

	err := r.db.QueryRow(ctx, q, product.Name, product.Price).Scan(&product.ID)
	if err != nil {
		if err.Error() == ErrProductNameUnique {
			return core.ErrProductNameDuplicate
		}

		return err
	}

	return nil
}

func (r *Product) Delete(ctx context.Context, id int) error {
	q := "DELETE FROM products WHERE id=$1"

	res, err := r.db.Exec(ctx, q, id)
	if err != nil {
		return err
	}

	if res.RowsAffected() == 0 {
		return core.ErrNotFound
	}

	return nil
}

func (r *Product) GetByID(ctx context.Context, id int) (core.Product, error) {
	q := "SELECT id, name, price FROM products WHERE id=$1"

	var product core.Product

	if err := r.db.QueryRow(ctx, q, id).Scan(&product.ID, &product.Name, &product.Price); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return core.Product{}, core.ErrNotFound
		}

		return core.Product{}, err
	}

	return product, nil
}

func (r *Product) Update(ctx context.Context, product core.Product) error {
	q := "UPDATE products SET name=$1, price=$2 WHERE id=$3"

	if _, err := r.db.Exec(ctx, q, product.Name, product.Price, product.ID); err != nil {
		if err.Error() == ErrProductNameUnique {
			return core.ErrProductNameDuplicate
		}

		return err
	}

	return nil
}
