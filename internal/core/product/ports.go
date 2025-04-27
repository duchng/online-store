package product

import (
	"context"

	"store-management/pkg/paging"
)

type ProductPersistencePort interface {
	// Category operations
	CreateCategory(ctx context.Context, category Category) (Category, error)
	UpdateCategory(ctx context.Context, category Category) error
	DeleteCategory(ctx context.Context, id int) error
	GetCategory(ctx context.Context, id int) (Category, error)
	ListCategories(ctx context.Context) ([]Category, error)

	// Product operations
	CreateProduct(ctx context.Context, product Product) (Product, error)
	UpdateProduct(ctx context.Context, product Product) error
	CreateProductCategories(ctx context.Context, productId int, categoryIds []int) error
	DeleteAllProductCategories(ctx context.Context, productId int) error
	DeleteProduct(ctx context.Context, id int) error
	GetProduct(ctx context.Context, id int) (Product, error)
	ListProductsWithPagination(ctx context.Context, productFilters ProductFilters) (paging.Page[Product], error)
	ListProductsByCategory(ctx context.Context, categoryId int) ([]Product, error)

	// Review operations
	CreateReview(ctx context.Context, review Review) (Review, error)
	ListReviews(ctx context.Context, productName string) ([]Review, error)
	DeleteReview(ctx context.Context, id int) error
}
