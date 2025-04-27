package postgres

import (
	"context"
	"fmt"
	"slices"

	"github.com/uptrace/bun"

	productCore "store-management/internal/core/product"
	"store-management/pkg/apperrors"
	"store-management/pkg/database"
	"store-management/pkg/paging"
)

type ProductPostgresAdapter struct {
	GetDbFunc database.GetDbFunc
}

func (a *ProductPostgresAdapter) CreateProductCategories(ctx context.Context, productId int, categoryIds []int) error {
	if len(categoryIds) == 0 {
		return nil
	}
	db := a.GetDbFunc(ctx)
	productCategories := make([]ProductCategory, 0, len(categoryIds))
	for _, category := range categoryIds {
		productCategory := ProductCategory{
			ProductID:  productId,
			CategoryID: category,
		}
		productCategories = append(productCategories, productCategory)
	}
	_, err := db.NewInsert().Model(&productCategories).Exec(ctx)
	if err != nil {
		return fmt.Errorf("ProductPostgresAdapter.CreateProductCategories: %w", apperrors.FromError(err))
	}
	return nil
}

func (a *ProductPostgresAdapter) DeleteAllProductCategories(ctx context.Context, productId int) error {
	db := a.GetDbFunc(ctx)
	_, err := db.NewDelete().Model((*ProductCategory)(nil)).Where("product_id = ?", productId).Exec(ctx)
	if err != nil {
		return fmt.Errorf("ProductPostgresAdapter.DeleteAllProductCategories: %w", apperrors.FromError(err))
	}
	return nil
}

func NewProductPostgresAdapter(getDb database.GetDbFunc) *ProductPostgresAdapter {
	return &ProductPostgresAdapter{
		GetDbFunc: getDb,
	}
}

// Category operations
func (a *ProductPostgresAdapter) CreateCategory(ctx context.Context, category productCore.Category) (productCore.Category, error) {
	db := a.GetDbFunc(ctx)
	toCreate := MapCategoryDomainToDb(category)
	created := Category{}
	_, err := db.NewInsert().Model(&toCreate).Returning("*").Exec(ctx, &created)
	if err != nil {
		return productCore.Category{}, fmt.Errorf("ProductPostgresAdapter.CreateCategory: %w", apperrors.FromError(err))
	}
	return MapCategoryDbToDomain(created), nil
}

func (a *ProductPostgresAdapter) UpdateCategory(ctx context.Context, category productCore.Category) error {
	db := a.GetDbFunc(ctx)
	toUpdate := MapCategoryDomainToDb(category)
	_, err := db.NewUpdate().Model(&toUpdate).WherePK().Exec(ctx)
	if err != nil {
		return fmt.Errorf("ProductPostgresAdapter.UpdateCategory: %w", apperrors.FromError(err))
	}
	return nil
}

func (a *ProductPostgresAdapter) DeleteCategory(ctx context.Context, id int) error {
	db := a.GetDbFunc(ctx)
	_, err := db.NewDelete().Model((*Category)(nil)).Where("id = ?", id).Exec(ctx)
	if err != nil {
		return fmt.Errorf("ProductPostgresAdapter.DeleteCategory: %w", apperrors.FromError(err))
	}
	return nil
}

func (a *ProductPostgresAdapter) GetCategory(ctx context.Context, id int) (productCore.Category, error) {
	db := a.GetDbFunc(ctx)
	category := Category{}
	err := db.NewSelect().Model(&category).Where("id = ?", id).Scan(ctx)
	if err != nil {
		return productCore.Category{}, fmt.Errorf("ProductPostgresAdapter.GetCategory: %w", apperrors.FromError(err))
	}
	return MapCategoryDbToDomain(category), nil
}

type categoryCount struct {
	CategoryID int `bun:"category_id"`
	Count      int `bun:"count"`
}

func (a *ProductPostgresAdapter) ListCategories(ctx context.Context) ([]productCore.Category, error) {
	db := a.GetDbFunc(ctx)
	var categories []Category
	err := db.NewSelect().Model(&categories).Order("id ASC").Scan(ctx)
	if err != nil {
		return nil, fmt.Errorf("ProductPostgresAdapter.ListCategories: %w", apperrors.FromError(err))
	}
	var categoryCount []categoryCount
	err = db.NewSelect().
		Table("product_categories").
		Column("category_id").
		ColumnExpr("count(*)").
		Group("category_id").
		Scan(ctx, &categoryCount)
	if err != nil {
		return nil, fmt.Errorf("ProductPostgresAdapter.ListCategories: %w", apperrors.FromError(err))
	}
	mapCategoryCount := make(map[int]int, len(categoryCount))
	for _, count := range categoryCount {
		mapCategoryCount[count.CategoryID] = count.Count
	}

	result := make([]productCore.Category, len(categories))
	for i, cat := range categories {
		categoryModel := MapCategoryDbToDomain(cat)
		categoryModel.Total = mapCategoryCount[cat.ID]
		result[i] = categoryModel
	}
	return result, nil
}

func (a *ProductPostgresAdapter) CreateProduct(ctx context.Context, product productCore.Product) (productCore.Product, error) {
	db := a.GetDbFunc(ctx)
	toCreate := MapProductDomainToDb(product)
	created := Product{}
	_, err := db.NewInsert().Model(&toCreate).Returning("*").Exec(ctx, &created)
	if err != nil {
		return productCore.Product{}, fmt.Errorf("ProductPostgresAdapter.CreateProduct: %w", apperrors.FromError(err))
	}

	return MapProductDbToDomain(created), nil
}

func (a *ProductPostgresAdapter) UpdateProduct(ctx context.Context, product productCore.Product) error {
	db := a.GetDbFunc(ctx)

	// Update product
	toUpdate := MapProductDomainToDb(product)
	_, err := db.NewUpdate().Model(&toUpdate).WherePK().Exec(ctx)
	if err != nil {
		return fmt.Errorf("ProductPostgresAdapter.UpdateProduct: %w", apperrors.FromError(err))
	}

	// Delete existing product-category relationships
	_, err = db.NewDelete().Model((*ProductCategory)(nil)).Where("product_id = ?", product.ID).Exec(ctx)
	if err != nil {
		return fmt.Errorf("ProductPostgresAdapter.UpdateProduct (delete categories): %w", apperrors.FromError(err))
	}

	return nil
}

func (a *ProductPostgresAdapter) DeleteProduct(ctx context.Context, id int) error {
	db := a.GetDbFunc(ctx)

	// Delete product-category relationships
	_, err := db.NewDelete().Model((*ProductCategory)(nil)).Where("product_id = ?", id).Exec(ctx)
	if err != nil {
		return fmt.Errorf("ProductPostgresAdapter.DeleteProduct (categories): %w", apperrors.FromError(err))
	}

	// Delete product
	_, err = db.NewDelete().Model((*Product)(nil)).Where("id = ?", id).Exec(ctx)
	if err != nil {
		return fmt.Errorf("ProductPostgresAdapter.DeleteProduct: %w", apperrors.FromError(err))
	}

	return nil
}

func (a *ProductPostgresAdapter) GetProduct(ctx context.Context, id int) (productCore.Product, error) {
	db := a.GetDbFunc(ctx)
	product := Product{}
	err := db.NewSelect().Model(&product).Where("id = ?", id).Scan(ctx)
	if err != nil {
		return productCore.Product{}, fmt.Errorf("ProductPostgresAdapter.GetProduct: %w", apperrors.FromError(err))
	}
	return MapProductDbToDomain(product), nil
}

func (a *ProductPostgresAdapter) ListProductsWithPagination(ctx context.Context, productFilters productCore.ProductFilters) (paging.Page[productCore.Product], error) {
	db := a.GetDbFunc(ctx)
	var products []Product
	query := db.NewSelect().Model(&products)
	AppyProductFilter(query, productFilters)
	if productFilters.Cursor >= 0 {
		query.Order("id ASC")
		query.Where("id > ?", productFilters.Cursor)
	} else {
		query.Order("id DESC")
		query.Where("id < ?", -productFilters.Cursor)
	}
	query.Limit(productFilters.Size + 1)

	err := query.Scan(ctx)
	if err != nil {
		return paging.Page[productCore.Product]{}, fmt.Errorf(
			"ProductPostgresAdapter.ListProducts: %w", apperrors.FromError(err),
		)
	}
	metadata := paging.MetaData{
		PageSize: productFilters.Size,
	}
	if len(products) > productFilters.Size { // extra record, so there must be a next or previous page
		if productFilters.Cursor >= 0 {
			metadata.HasNext = true
		} else {
			metadata.HasPrevious = true
		}
		products = products[:productFilters.Size]
	}
	if productFilters.Cursor < 0 {
		metadata.HasNext = true
		slices.Reverse(products)
	} else {
		metadata.HasPrevious = true
	}

	result := make([]productCore.Product, len(products))
	for i, p := range products {
		result[i] = MapProductDbToDomain(p)
	}
	return paging.Page[productCore.Product]{
		Data:     result,
		Metadata: metadata,
	}, nil
}

func AppyProductFilter(query *bun.SelectQuery, productFilters productCore.ProductFilters) *bun.SelectQuery {
	if productFilters.Name != "" {
		query.Where("name ILIKE ?", "%"+productFilters.Name+"%")
	}

	if len(productFilters.Statuses) > 0 {
		query.Where("status IN (?)", bun.In(productFilters.Statuses))
	}

	return query
}

func (a *ProductPostgresAdapter) ListProductsByCategory(ctx context.Context, categoryId int) ([]productCore.Product, error) {
	db := a.GetDbFunc(ctx)
	var products []Product

	err := db.NewSelect().
		Model(&products).
		Join("JOIN product_categories pc ON pc.product_id = product.id").
		Where("pc.category_id = ?", categoryId).
		Order("product.id ASC").
		Scan(ctx)
	if err != nil {
		return nil, fmt.Errorf("ProductPostgresAdapter.ListProductsByCategory: %w", apperrors.FromError(err))
	}

	result := make([]productCore.Product, len(products))
	for i, p := range products {
		result[i] = MapProductDbToDomain(p)
	}
	return result, nil
}

func (a *ProductPostgresAdapter) CreateReview(ctx context.Context, review productCore.Review) (productCore.Review, error) {
	db := a.GetDbFunc(ctx)
	toCreate := MapReviewDomainToDb(review)
	created := Review{}
	_, err := db.NewInsert().Model(&toCreate).Returning("*").Exec(ctx, &created)
	if err != nil {
		return productCore.Review{}, fmt.Errorf("ProductPostgresAdapter.CreateReview: %w", apperrors.FromError(err))
	}
	return MapReviewDbToDomain(created), nil
}

func (a *ProductPostgresAdapter) ListReviews(ctx context.Context, productName string) ([]productCore.Review, error) {
	db := a.GetDbFunc(ctx)
	var reviews []Review
	query := db.NewSelect().
		ColumnExpr("reviews.*").
		Table("reviews")

	if productName != "" {
		query.Join("JOIN products p ON reviews.product_id = p.id").
			Where("p.name ILIKE ?", "%"+productName+"%")
	}

	err := query.Order("reviews.created_at DESC").Scan(ctx, &reviews)
	if err != nil {
		return nil, fmt.Errorf("ProductPostgresAdapter.ListReviews: %w", apperrors.FromError(err))
	}

	result := make([]productCore.Review, len(reviews))
	for i, review := range reviews {
		result[i] = MapReviewDbToDomain(review)
	}
	return result, nil
}

func (a *ProductPostgresAdapter) DeleteReview(ctx context.Context, id int) error {
	db := a.GetDbFunc(ctx)
	_, err := db.NewDelete().Model((*Review)(nil)).Where("id = ?", id).Exec(ctx)
	if err != nil {
		return fmt.Errorf("ProductPostgresAdapter.DeleteReview: %w", apperrors.FromError(err))
	}
	return nil
}

var _ productCore.ProductPersistencePort = (*ProductPostgresAdapter)(nil)
