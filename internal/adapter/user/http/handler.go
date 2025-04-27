package http

import (
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"

	"store-management/internal/core/user"
	"store-management/pkg/appcontext"
	"store-management/pkg/apperrors"
)

type UserHandler struct {
	useCase user.UseCase
}

func NewUserHandler(useCase user.UseCase) *UserHandler {
	return &UserHandler{
		useCase: useCase,
	}
}

// SignIn handles user authentication
// @Summary Sign in user
// @Description Authenticate a user and return access token
// @Tags users
// @Accept json
// @Produce json
// @Param request body SignInRequest true "Sign in credentials"
// @Success 200 {object} user.Token
// @Failure 400 {object} apperrors.Error
// @Failure 401 {object} apperrors.Error
// @Router /signin [post]
func (h *UserHandler) SignIn(c echo.Context) error {
	var req SignInRequest
	if err := c.Bind(&req); err != nil {
		return apperrors.NewError(http.StatusBadRequest, "INVALID_REQUEST", "invalid request body")
	}

	if err := c.Validate(req); err != nil {
		return apperrors.NewError(http.StatusBadRequest, "VALIDATION_ERROR", err.Error())
	}

	token, err := h.useCase.SignIn(c.Request().Context(), req.UserName, req.Password)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, token)
}

// SignUp handles user registration
// @Summary Register new user
// @Description Create a new user account
// @Tags users
// @Accept json
// @Produce json
// @Param request body SignUpRequest true "User registration details"
// @Success 201 {object} user.User
// @Failure 400 {object} apperrors.Error
// @Router /signup [post]
func (h *UserHandler) SignUp(c echo.Context) error {
	var req SignUpRequest
	if err := c.Bind(&req); err != nil {
		return apperrors.NewError(http.StatusBadRequest, "INVALID_REQUEST", "invalid request body")
	}

	if err := c.Validate(req); err != nil {
		return apperrors.NewError(http.StatusBadRequest, "VALIDATION_ERROR", err.Error())
	}

	newUser := user.User{
		Username: req.UserName,
		Email:    req.Email,
		FullName: req.FullName,
		Role:     user.UserRole(req.Role),
	}

	createdUser, err := h.useCase.SignUp(c.Request().Context(), newUser, req.Password)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusCreated, createdUser)
}

// AddToWishList adds a product to user's wishlist
// @Summary Add product to wishlist
// @Description Add a product to the authenticated user's wishlist
// @Tags wishlist
// @Accept json
// @Produce json
// @Param productId path int true "Product ID"
// @Success 200 {object} string
// @Failure 400 {object} apperrors.Error
// @Failure 401 {object} apperrors.Error
// @Security BearerAuth
// @Router /wishlist/{productId} [post]
func (h *UserHandler) AddToWishList(c echo.Context) error {
	userId, err := appcontext.ContextGetUserId(c)
	if err != nil {
		return err
	}

	productId, err := strconv.Atoi(c.Param("productId"))
	if err != nil {
		return apperrors.NewError(http.StatusBadRequest, "INVALID_PRODUCT_ID", "invalid product id")
	}

	if err := h.useCase.AddToWishList(c.Request().Context(), userId, productId); err != nil {
		return err
	}

	return c.JSON(http.StatusOK, map[string]string{"message": "product added to wishlist"})
}

// RemoveFromWishList removes a product from user's wishlist
// @Summary Remove product from wishlist
// @Description Remove a product from the authenticated user's wishlist
// @Tags wishlist
// @Accept json
// @Produce json
// @Param productId path int true "Product ID"
// @Success 200 {object} string
// @Failure 400 {object} apperrors.Error
// @Failure 401 {object} apperrors.Error
// @Security BearerAuth
// @Router /wishlist/{productId} [delete]
func (h *UserHandler) RemoveFromWishList(c echo.Context) error {
	userId, err := appcontext.ContextGetUserId(c)
	if err != nil {
		return err
	}

	productId, err := strconv.Atoi(c.Param("productId"))
	if err != nil {
		return apperrors.NewError(http.StatusBadRequest, "INVALID_PRODUCT_ID", "invalid product id")
	}

	if err := h.useCase.RemoveFromWishList(c.Request().Context(), userId, productId); err != nil {
		return err
	}

	return c.JSON(http.StatusOK, map[string]string{"message": "product removed from wishlist"})
}

// GetWishList retrieves user's wishlist
// @Summary Get user's wishlist
// @Description Get the authenticated user's wishlist
// @Tags wishlist
// @Accept json
// @Produce json
// @Success 200 {array} product.Product
// @Failure 400 {object} apperrors.Error
// @Failure 401 {object} apperrors.Error
// @Security BearerAuth
// @Router /wishlist [get]
func (h *UserHandler) GetWishList(c echo.Context) error {
	userId, err := appcontext.ContextGetUserId(c)
	if err != nil {
		return err
	}

	products, err := h.useCase.GetWishList(c.Request().Context(), userId)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, products)
}

// GetProfile retrieves the authenticated user's profile
// @Summary Get user profile
// @Description Get the authenticated user's profile information
// @Tags users
// @Accept json
// @Produce json
// @Success 200 {object} user.User
// @Failure 401 {object} apperrors.Error
// @Security BearerAuth
// @Router /profile [get]
func (h *UserHandler) GetProfile(c echo.Context) error {
	userId, err := appcontext.ContextGetUserId(c)
	if err != nil {
		return err
	}

	user, err := h.useCase.GetProfile(c.Request().Context(), userId)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, user)
}

// ChangePassword updates the authenticated user's password
// @Summary Change password
// @Description Change the authenticated user's password
// @Tags users
// @Accept json
// @Produce json
// @Param request body ChangePasswordRequest true "Password change details"
// @Success 200 {object} string
// @Failure 400 {object} apperrors.Error
// @Failure 401 {object} apperrors.Error
// @Security BearerAuth
// @Router /change-password [post]
func (h *UserHandler) ChangePassword(c echo.Context) error {
	userId, err := appcontext.ContextGetUserId(c)
	if err != nil {
		return err
	}

	var req ChangePasswordRequest
	if err := c.Bind(&req); err != nil {
		return err
	}

	if err := c.Validate(req); err != nil {
		return err
	}

	if err := h.useCase.ChangePassword(
		c.Request().Context(), userId, req.CurrentPassword, req.NewPassword,
	); err != nil {
		return err
	}

	return c.JSON(http.StatusOK, map[string]string{"message": "password changed successfully"})
}

// ListUsers retrieves all users (admin only)
// @Summary List all users
// @Description Get a list of all users (requires admin role)
// @Tags users
// @Accept json
// @Produce json
// @Param search query string false "Search by username or email"
// @Param roles query []string false "Filter by roles (user, admin)"
// @Success 200 {array} user.User
// @Failure 401 {object} apperrors.Error
// @Failure 403 {object} apperrors.Error
// @Security BearerAuth
// @Router /admin/users [get]
func (h *UserHandler) ListUsers(c echo.Context) error {
	filters := user.UserFilters{
		Search: c.QueryParam("search"),
		Roles:  c.QueryParams()["roles"],
	}

	users, err := h.useCase.ListUsers(c.Request().Context(), filters)
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, users)
}

// UpdateUserRole updates a user's role (admin only)
// @Summary Update user role
// @Description Update a user's role (requires admin role)
// @Tags users
// @Accept json
// @Produce json
// @Param id path int true "User ID"
// @Param request body UpdateUserRoleRequest true "New role details"
// @Success 200 {object} string
// @Failure 400 {object} apperrors.Error
// @Failure 401 {object} apperrors.Error
// @Failure 403 {object} apperrors.Error
// @Security BearerAuth
// @Router /admin/users/{id}/role [put]
func (h *UserHandler) UpdateUserRole(c echo.Context) error {
	userId, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return apperrors.NewError(http.StatusBadRequest, "INVALID_USER_ID", "invalid user id")
	}

	var req UpdateUserRoleRequest
	if err := c.Bind(&req); err != nil {
		return err
	}

	if err := c.Validate(req); err != nil {
		return err
	}

	if err := h.useCase.UpdateUserRole(c.Request().Context(), userId, user.UserRole(req.Role)); err != nil {
		return err
	}

	return c.JSON(http.StatusOK, map[string]string{"message": "user role updated successfully"})
}
