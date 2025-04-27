package http

import (
	"store-management/internal/core/product"
	"store-management/pkg/paging"
)

// ProductPage represents a paginated list of products
type ProductPage struct {
	Data     []product.Product `json:"data"`
	Metadata paging.MetaData   `json:"metadata"`
}

type CreateCategoryRequest struct {
	Name        string `json:"name" validate:"required"`
	Description string `json:"description" validate:"required"`
}

type UpdateCategoryRequest struct {
	Name        string `json:"name" validate:"required"`
	Description string `json:"description" validate:"required"`
}

type CreateProductRequest struct {
	Name          string  `json:"name" validate:"required"`
	Description   string  `json:"description" validate:"required"`
	Price         float64 `json:"price" validate:"required,gt=0"`
	StockQuantity int     `json:"stockQuantity" validate:"required,gte=0"`
	CategoryIds   []int   `json:"categoryIds"`
}

type UpdateProductRequest struct {
	Name          string  `json:"name" validate:"required"`
	Description   string  `json:"description" validate:"required"`
	Price         float64 `json:"price" validate:"required,gt=0"`
	StockQuantity int     `json:"stockQuantity" validate:"required,gte=0"`
}

type CreateReviewRequest struct {
	Rating  int    `json:"rating" validate:"required,min=1,max=5"`
	Comment string `json:"comment" validate:"required"`
}
