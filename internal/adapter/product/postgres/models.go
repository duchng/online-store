package postgres

import (
	"time"

	"github.com/uptrace/bun"

	"store-management/internal/core/product"
)

type Product struct {
	bun.BaseModel `bun:"table:products"`

	ID            int `bun:",pk,nullzero"`
	Name          string
	Description   string
	Price         float64
	StockQuantity int
	Status        product.ProductStatus
	CreatedAt     time.Time `bun:",nullzero"`
	UpdatedAt     time.Time `bun:",nullzero"`
}

type Category struct {
	bun.BaseModel `bun:"table:categories"`

	ID          int `bun:",pk,nullzero"`
	Name        string
	Description string
	CreatedAt   time.Time `bun:",nullzero"`
	UpdatedAt   time.Time `bun:",nullzero"`
}

type ProductCategory struct {
	bun.BaseModel `bun:"table:product_categories"`

	ProductID  int `bun:",pk"`
	CategoryID int `bun:",pk"`
}

type Review struct {
	bun.BaseModel `bun:"table:reviews"`

	ID        int `bun:",pk,nullzero"`
	ProductID int
	UserID    int
	Rating    int
	Comment   string
	CreatedAt time.Time `bun:",nullzero"`
	UpdatedAt time.Time `bun:",nullzero"`
}
