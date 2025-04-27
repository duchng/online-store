package product

import (
	"time"

	"store-management/pkg/paging"
)

// Product represents a product in the system
type Product struct {
	ID            int           `json:"id"`
	Name          string        `json:"name"`
	Description   string        `json:"description"`
	Price         float64       `json:"price"`
	StockQuantity int           `json:"stockQuantity"`
	Status        ProductStatus `json:"status"`
	CreatedAt     time.Time     `json:"createdAt"`
	UpdatedAt     time.Time     `json:"updatedAt"`
}

type ProductFilters struct {
	paging.Paging

	Statuses []string `query:"statuses"`
	Name     string   `query:"name"`
}

type ProductStatus string

const (
	ProductStatusInStock    ProductStatus = "IN_STOCK"
	ProductStatusOutOfStock ProductStatus = "OUT_OF_STOCK"
)

// Category represents a product category
type Category struct {
	ID          int       `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`

	Total int `json:"total"`
}

// ProductCategory represents the many-to-many relationship between products and categories
type ProductCategory struct {
	ProductID  int `json:"productId"`
	CategoryID int `json:"categoryId"`
}

// Review represents a product review
type Review struct {
	ID        int       `json:"id"`
	ProductID int       `json:"productId"`
	UserID    int       `json:"userId"`
	Rating    int       `json:"rating"`
	Comment   string    `json:"comment"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}
