package core

type Product struct {
	ID    int    `json:"id" csv:"-"`
	Name  string `json:"name" csv:"PRODUCT NAME"`
	Price int    `json:"price" csv:"PRICE"`
}

type CreateProductInput struct {
	Name  string `json:"name" validate:"required"`
	Price int    `json:"price" validate:"number,min=0"`
}

func (c CreateProductInput) Validate() error {
	return validateStruct(c)
}

type UpdateProductInput struct {
	Name  string `json:"name" validate:"required"`
	Price int    `json:"price" validate:"number,min=0"`
}

func (c UpdateProductInput) Validate() error {
	return validateStruct(c)
}
