package http

import (
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"

	"store-management/internal/core/product"
	"store-management/pkg/appcontext"
	"store-management/pkg/apperrors"
	"store-management/pkg/paging"
)

type ProductHandler struct {
	useCase product.UseCase
}

func NewProductHandler(useCase product.UseCase) *ProductHandler {
	return &ProductHandler{
		useCase: useCase,
	}
}

// CreateCategory creates a new product category
// @Summary Create category
// @Description Create a new product category (admin only)
// @Tags categories
// @Accept json
// @Produce json
// @Param request body CreateCategoryRequest true "Category details"
// @Success 201 {object} product.Category
// @Failure 400 {object} apperrors.Error
// @Failure 401 {object} apperrors.Error
// @Security BearerAuth
// @Router /categories [post]
func (h *ProductHandler) CreateCategory(c echo.Context) error {
	var req CreateCategoryRequest
	if err := c.Bind(&req); err != nil {
		return err
	}

	if err := c.Validate(req); err != nil {
		return err
	}

	category := product.Category{
		Name:        req.Name,
		Description: req.Description,
	}

	created, err := h.useCase.CreateCategory(c.Request().Context(), category)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusCreated, created)
}

// UpdateCategory updates an existing product category
// @Summary Update category
// @Description Update an existing product category (admin only)
// @Tags categories
// @Accept json
// @Produce json
// @Param id path int true "Category ID"
// @Param request body UpdateCategoryRequest true "Category details"
// @Success 200 {object} string
// @Failure 400 {object} apperrors.Error
// @Failure 401 {object} apperrors.Error
// @Security BearerAuth
// @Router /categories/{id} [put]
func (h *ProductHandler) UpdateCategory(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return apperrors.NewError(http.StatusBadRequest, "INVALID_CATEGORY_ID", "invalid category id")
	}

	var req UpdateCategoryRequest
	if err := c.Bind(&req); err != nil {
		return err
	}

	if err := c.Validate(req); err != nil {
		return err
	}

	category := product.Category{
		ID:          id,
		Name:        req.Name,
		Description: req.Description,
	}

	if err := h.useCase.UpdateCategory(c.Request().Context(), category); err != nil {
		return err
	}

	return c.JSON(http.StatusOK, map[string]string{"message": "category updated successfully"})
}

// DeleteCategory deletes a product category
// @Summary Delete category
// @Description Delete a product category (admin only)
// @Tags categories
// @Accept json
// @Produce json
// @Param id path int true "Category ID"
// @Success 200 {object} string
// @Failure 400 {object} apperrors.Error
// @Failure 401 {object} apperrors.Error
// @Security BearerAuth
// @Router /categories/{id} [delete]
func (h *ProductHandler) DeleteCategory(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return apperrors.NewError(http.StatusBadRequest, "INVALID_CATEGORY_ID", "invalid category id")
	}

	if err := h.useCase.DeleteCategory(c.Request().Context(), id); err != nil {
		return err
	}

	return c.JSON(http.StatusOK, map[string]string{"message": "category deleted successfully"})
}

// GetCategory retrieves a product category by ID
// @Summary Get category
// @Description Get a product category by ID
// @Tags categories
// @Accept json
// @Produce json
// @Param id path int true "Category ID"
// @Success 200 {object} product.Category
// @Failure 400 {object} apperrors.Error
// @Router /categories/{id} [get]
func (h *ProductHandler) GetCategory(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return apperrors.NewError(http.StatusBadRequest, "INVALID_CATEGORY_ID", "invalid category id")
	}

	category, err := h.useCase.GetCategory(c.Request().Context(), id)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, category)
}

// ListCategories retrieves all product categories
// @Summary List categories
// @Description Get all product categories
// @Tags categories
// @Accept json
// @Produce json
// @Success 200 {array} product.Category
// @Router /categories [get]
func (h *ProductHandler) ListCategories(c echo.Context) error {
	categories, err := h.useCase.ListCategories(c.Request().Context())
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, categories)
}

// CreateProduct creates a new product
// @Summary Create product
// @Description Create a new product (admin only)
// @Tags products
// @Accept json
// @Produce json
// @Param request body CreateProductRequest true "Product details"
// @Success 201 {object} product.Product
// @Failure 400 {object} apperrors.Error
// @Failure 401 {object} apperrors.Error
// @Security BearerAuth
// @Router /products [post]
func (h *ProductHandler) CreateProduct(c echo.Context) error {
	var req CreateProductRequest
	if err := c.Bind(&req); err != nil {
		return err
	}

	if err := c.Validate(req); err != nil {
		return err
	}
	status := product.ProductStatusInStock
	if req.StockQuantity == 0 {
		status = product.ProductStatusOutOfStock
	}
	newProduct := product.Product{
		Name:          req.Name,
		Description:   req.Description,
		Price:         req.Price,
		StockQuantity: req.StockQuantity,
		Status:        status,
	}

	created, err := h.useCase.CreateProduct(c.Request().Context(), newProduct, req.CategoryIds)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusCreated, created)
}

// UpdateProduct updates an existing product
// @Summary Update product
// @Description Update an existing product (admin only)
// @Tags products
// @Accept json
// @Produce json
// @Param id path int true "Product ID"
// @Param request body UpdateProductRequest true "Product details"
// @Success 200 {object} string
// @Failure 400 {object} apperrors.Error
// @Failure 401 {object} apperrors.Error
// @Security BearerAuth
// @Router /products/{id} [put]
func (h *ProductHandler) UpdateProduct(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return apperrors.NewError(http.StatusBadRequest, "INVALID_PRODUCT_ID", "invalid product id")
	}

	var req UpdateProductRequest
	if err := c.Bind(&req); err != nil {
		return err
	}

	if err := c.Validate(req); err != nil {
		return err
	}

	productToUpdate := product.Product{
		ID:            id,
		Name:          req.Name,
		Description:   req.Description,
		Price:         req.Price,
		StockQuantity: req.StockQuantity,
	}

	if err := h.useCase.UpdateProduct(c.Request().Context(), productToUpdate); err != nil {
		return err
	}

	return c.JSON(http.StatusOK, map[string]string{"message": "product updated successfully"})
}

// DeleteProduct deletes a product
// @Summary Delete product
// @Description Delete a product (admin only)
// @Tags products
// @Accept json
// @Produce json
// @Param id path int true "Product ID"
// @Success 200 {object} string
// @Failure 400 {object} apperrors.Error
// @Failure 401 {object} apperrors.Error
// @Security BearerAuth
// @Router /products/{id} [delete]
func (h *ProductHandler) DeleteProduct(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return apperrors.NewError(http.StatusBadRequest, "INVALID_PRODUCT_ID", "invalid product id")
	}

	if err := h.useCase.DeleteProduct(c.Request().Context(), id); err != nil {
		return err
	}

	return c.JSON(http.StatusOK, map[string]string{"message": "product deleted successfully"})
}

// GetProduct retrieves a product by ID
// @Summary Get product
// @Description Get a product by ID
// @Tags products
// @Accept json
// @Produce json
// @Param id path int true "Product ID"
// @Success 200 {object} product.Product
// @Failure 400 {object} apperrors.Error
// @Router /products/{id} [get]
func (h *ProductHandler) GetProduct(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return apperrors.NewError(http.StatusBadRequest, "INVALID_PRODUCT_ID", "invalid product id")
	}

	product, err := h.useCase.GetProduct(c.Request().Context(), id)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, product)
}

// ListProducts retrieves all products with filtering and pagination
// @Summary List products
// @Description Get all products with optional filtering and pagination
// @Tags products
// @Accept json
// @Produce json
// @Param name query string false "Filter by product name"
// @Param statuses query []string false "Filter by product statuses (IN_STOCK, OUT_OF_STOCK)"
// @Param size query int false "Page size (default: 20, max: 200)"
// @Param cursor query int false "Cursor for keyset pagination"
// @Param sort query []string false "Sort orders (e.g., name ASC, price DESC)"
// @Success 200 {object} ProductPage
// @Failure 400 {object} apperrors.Error
// @Router /products [get]
func (h *ProductHandler) ListProducts(c echo.Context) error {
	filter, p, err := paging.ParseRequestWithKeysetPagination[product.ProductFilters](c)
	if err != nil {
		return apperrors.NewError(http.StatusBadRequest, "INVALID_REQUEST", "invalid request")
	}
	filter.Size = p.Size
	filter.Cursor = p.Cursor
	filter.Sort = p.Sort
	products, err := h.useCase.ListProductsWithPagination(c.Request().Context(), filter)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, products)
}

// ListProductsByCategory retrieves all products in a category
// @Summary List products by category
// @Description Get all products in a specific category
// @Tags products
// @Accept json
// @Produce json
// @Param categoryId path int true "Category ID"
// @Success 200 {array} product.Product
// @Failure 400 {object} apperrors.Error
// @Router /categories/{categoryId}/products [get]
func (h *ProductHandler) ListProductsByCategory(c echo.Context) error {
	categoryId, err := strconv.Atoi(c.Param("categoryId"))
	if err != nil {
		return apperrors.NewError(http.StatusBadRequest, "INVALID_CATEGORY_ID", "invalid category id")
	}

	products, err := h.useCase.ListProductsByCategory(c.Request().Context(), categoryId)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, products)
}

// CreateReview creates a new product review
// @Summary Create review
// @Description Create a new product review (authenticated users only)
// @Tags reviews
// @Accept json
// @Produce json
// @Param productId path int true "Product ID"
// @Param request body CreateReviewRequest true "Review details"
// @Success 201 {object} product.Review
// @Failure 400 {object} apperrors.Error
// @Failure 401 {object} apperrors.Error
// @Security BearerAuth
// @Router /products/{productId}/reviews [post]
func (h *ProductHandler) CreateReview(c echo.Context) error {
	productId, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return apperrors.NewError(http.StatusBadRequest, "INVALID_PRODUCT_ID", "invalid product id")
	}

	var req CreateReviewRequest
	if err := c.Bind(&req); err != nil {
		return err
	}

	if err := c.Validate(req); err != nil {
		return err
	}

	userId, err := appcontext.ContextGetUserId(c)
	if err != nil {
		return err
	}

	review := product.Review{
		ProductID: productId,
		UserID:    userId,
		Rating:    req.Rating,
		Comment:   req.Comment,
	}

	created, err := h.useCase.CreateReview(c.Request().Context(), review)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusCreated, created)
}

// ListReviews retrieves all product reviews with optional product name filter (admin only)
// @Summary List reviews
// @Description Get all product reviews with optional product name filter (admin only)
// @Tags reviews
// @Accept json
// @Produce json
// @Param productName query string false "Filter by product name"
// @Success 200 {array} product.Review
// @Failure 401 {object} apperrors.Error
// @Failure 403 {object} apperrors.Error
// @Security BearerAuth
// @Router /admin/reviews [get]
func (h *ProductHandler) ListReviews(c echo.Context) error {
	productName := c.QueryParam("productName")
	reviews, err := h.useCase.ListReviews(c.Request().Context(), productName)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, reviews)
}

// DeleteReview deletes a review by ID (admin only)
// @Summary Delete review
// @Description Delete a review by ID (admin only)
// @Tags reviews
// @Accept json
// @Produce json
// @Param id path int true "Review ID"
// @Success 200 {object} string
// @Failure 400 {object} apperrors.Error
// @Failure 401 {object} apperrors.Error
// @Failure 403 {object} apperrors.Error
// @Security BearerAuth
// @Router /admin/reviews/{id} [delete]
func (h *ProductHandler) DeleteReview(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return apperrors.NewError(http.StatusBadRequest, "INVALID_REVIEW_ID", "invalid review id")
	}

	if err := h.useCase.DeleteReview(c.Request().Context(), id); err != nil {
		return err
	}

	return c.JSON(http.StatusOK, map[string]string{"message": "review deleted successfully"})
}
