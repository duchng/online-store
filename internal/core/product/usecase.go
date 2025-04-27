package product

import (
	"context"
	"fmt"

	"store-management/pkg/atomicity"
	"store-management/pkg/paging"
)

type UseCase interface {
	// Category operations
	CreateCategory(ctx context.Context, category Category) (Category, error)
	UpdateCategory(ctx context.Context, category Category) error
	DeleteCategory(ctx context.Context, id int) error
	GetCategory(ctx context.Context, id int) (Category, error)
	ListCategories(ctx context.Context) ([]Category, error)

	// Product operations
	CreateProduct(ctx context.Context, product Product, categoryIds []int) (Product, error)
	UpdateProduct(ctx context.Context, product Product) error
	DeleteProduct(ctx context.Context, id int) error
	GetProduct(ctx context.Context, id int) (Product, error)
	ListProductsWithPagination(ctx context.Context, productFilters ProductFilters) (paging.Page[Product], error)
	ListProductsByCategory(ctx context.Context, categoryId int) ([]Product, error)

	// Review operations
	CreateReview(ctx context.Context, review Review) (Review, error)
	ListReviews(ctx context.Context, productName string) ([]Review, error)
	DeleteReview(ctx context.Context, id int) error
}

type usecase struct {
	persistencePort ProductPersistencePort
	atomicExecutor  atomicity.AtomicExecutor
}

func NewUseCase(persistencePort ProductPersistencePort, atomicExecutor atomicity.AtomicExecutor) UseCase {
	return &usecase{
		persistencePort: persistencePort,
		atomicExecutor:  atomicExecutor,
	}
}

func (u *usecase) CreateCategory(ctx context.Context, category Category) (Category, error) {
	created, err := u.persistencePort.CreateCategory(ctx, category)
	if err != nil {
		return Category{}, fmt.Errorf("usecase.CreateCategory: %w", err)
	}
	return created, nil
}

func (u *usecase) UpdateCategory(ctx context.Context, category Category) error {
	if err := u.persistencePort.UpdateCategory(ctx, category); err != nil {
		return fmt.Errorf("usecase.UpdateCategory: %w", err)
	}
	return nil
}

func (u *usecase) DeleteCategory(ctx context.Context, id int) error {
	if err := u.persistencePort.DeleteCategory(ctx, id); err != nil {
		return fmt.Errorf("usecase.DeleteCategory: %w", err)
	}
	return nil
}

func (u *usecase) GetCategory(ctx context.Context, id int) (Category, error) {
	category, err := u.persistencePort.GetCategory(ctx, id)
	if err != nil {
		return Category{}, fmt.Errorf("usecase.GetCategory: %w", err)
	}
	return category, nil
}

func (u *usecase) ListCategories(ctx context.Context) ([]Category, error) {
	categories, err := u.persistencePort.ListCategories(ctx)
	if err != nil {
		return nil, fmt.Errorf("usecase.ListCategories: %w", err)
	}
	return categories, nil
}

func (u *usecase) CreateProduct(ctx context.Context, product Product, categoryIds []int) (Product, error) {
	var res Product
	atomicFunc := func(ctx context.Context) error {
		created, err := u.persistencePort.CreateProduct(ctx, product)
		if err != nil {
			return err
		}
		res = created
		return u.persistencePort.CreateProductCategories(ctx, created.ID, categoryIds)
	}
	err := u.atomicExecutor.Execute(ctx, atomicFunc)
	if err != nil {
		return Product{}, fmt.Errorf("usecase.CreateProduct: %w", err)
	}
	return res, nil
}

func (u *usecase) UpdateProduct(ctx context.Context, product Product) error {
	if err := u.persistencePort.UpdateProduct(ctx, product); err != nil {
		return fmt.Errorf("usecase.UpdateProduct: %w", err)
	}
	return nil
}

func (u *usecase) DeleteProduct(ctx context.Context, id int) error {
	if err := u.persistencePort.DeleteProduct(ctx, id); err != nil {
		return fmt.Errorf("usecase.DeleteProduct: %w", err)
	}
	return nil
}

func (u *usecase) GetProduct(ctx context.Context, id int) (Product, error) {
	product, err := u.persistencePort.GetProduct(ctx, id)
	if err != nil {
		return Product{}, fmt.Errorf("usecase.GetProduct: %w", err)
	}
	return product, nil
}

func (u *usecase) ListProductsWithPagination(ctx context.Context, productFilter ProductFilters) (paging.Page[Product], error) {
	products, err := u.persistencePort.ListProductsWithPagination(ctx, productFilter)
	if err != nil {
		return paging.Page[Product]{}, fmt.Errorf("usecase.ListProducts: %w", err)
	}
	return products, nil
}

func (u *usecase) ListProductsByCategory(ctx context.Context, categoryId int) ([]Product, error) {
	products, err := u.persistencePort.ListProductsByCategory(ctx, categoryId)
	if err != nil {
		return nil, fmt.Errorf("usecase.ListProductsByCategory: %w", err)
	}
	return products, nil
}

func (u *usecase) CreateReview(ctx context.Context, review Review) (Review, error) {
	created, err := u.persistencePort.CreateReview(ctx, review)
	if err != nil {
		return Review{}, fmt.Errorf("usecase.CreateReview: %w", err)
	}
	return created, nil
}

func (u *usecase) ListReviews(ctx context.Context, productName string) ([]Review, error) {
	reviews, err := u.persistencePort.ListReviews(ctx, productName)
	if err != nil {
		return nil, fmt.Errorf("usecase.ListReviews: %w", err)
	}
	return reviews, nil
}

func (u *usecase) DeleteReview(ctx context.Context, id int) error {
	err := u.persistencePort.DeleteReview(ctx, id)
	if err != nil {
		return fmt.Errorf("usecase.DeleteReview: %w", err)
	}
	return nil
}

var _ UseCase = (*usecase)(nil)
