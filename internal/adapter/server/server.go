package server

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	echoMiddlewares "github.com/labstack/echo/v4/middleware"
	"github.com/samber/do/v2"
	echoSwagger "github.com/swaggo/echo-swagger"

	_ "store-management/docs" // swagger generated docs
	"store-management/internal/adapter/middlewares"
	productHttp "store-management/internal/adapter/product/http"
	userHttp "store-management/internal/adapter/user/http"
	"store-management/internal/config"
	"store-management/pkg/jwttoken"
	"store-management/pkg/shutdown"
)

// @title Store Management API
// @version 1.0
// @description This is a store management server.
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.url http://www.swagger.io/support
// @contact.email support@swagger.io

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @host localhost:8080
// @BasePath /
// @schemes http

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization

func ServeRest(
	injector do.Injector,
) error {
	logger := do.MustInvoke[*slog.Logger](injector)
	tasks := do.MustInvoke[*shutdown.Tasks](injector)
	cfg := do.MustInvoke[config.AppConfig](injector)
	signParser := do.MustInvoke[jwttoken.SignParser](injector)

	// Get handlers
	userHandler := do.MustInvoke[*userHttp.UserHandler](injector)
	productHandler := do.MustInvoke[*productHttp.ProductHandler](injector)
	e := echo.New()
	e.Validator = &CustomValidator{Validator: validator.New()}
	e.HTTPErrorHandler = middlewares.CustomHTTPErrorHandler
	e.Pre(
		echoMiddlewares.RemoveTrailingSlash(),
		echoMiddlewares.RequestID(),
		echoMiddlewares.Recover(),
		echoMiddlewares.Secure(),
		// echoMiddlewares.CSRF(),
		middlewares.AuthenticationMiddleware(signParser),
	)

	e.Use(echoMiddlewares.Logger())

	// Swagger docs
	e.GET("/swagger/*", echoSwagger.WrapHandler)

	// Public routes
	e.POST("/signin", userHandler.SignIn)
	e.POST("/signup", userHandler.SignUp)

	// Protected user routes
	userGroup := e.Group("")
	userGroup.Use(middlewares.RequireOneOfRoles("user", "admin"))
	{
		userGroup.GET("/profile", userHandler.GetProfile)
		userGroup.POST("/change-password", userHandler.ChangePassword)
	}

	// Admin user management routes
	adminUsersGroup := e.Group("/admin/users", middlewares.RequireOneOfRoles("admin"))
	{
		adminUsersGroup.GET("", userHandler.ListUsers)
		adminUsersGroup.PUT("/:id/role", userHandler.UpdateUserRole)
	}

	// Protected wishlist routes
	wishlistGroup := e.Group("/wishlist")
	wishlistGroup.Use(middlewares.RequireOneOfRoles("user", "admin"))
	{
		wishlistGroup.POST("/:productId", userHandler.AddToWishList)
		wishlistGroup.DELETE("/:productId", userHandler.RemoveFromWishList)
		wishlistGroup.GET("", userHandler.GetWishList)
	}

	// Categories endpoints
	categoriesGroup := e.Group("/categories")
	{
		// Public endpoints
		categoriesGroup.GET("", productHandler.ListCategories)
		categoriesGroup.GET("/:id", productHandler.GetCategory)
		categoriesGroup.GET("/:categoryId/products", productHandler.ListProductsByCategory)
	}

	// Products endpoints
	productsGroup := e.Group("/products")
	{
		// Public endpoints
		productsGroup.GET("", productHandler.ListProducts)
		productsGroup.GET("/:id", productHandler.GetProduct)

		// Protected review endpoints
		productsGroup.POST("/:id/reviews", productHandler.CreateReview, middlewares.RequireOneOfRoles("user", "admin"))
	}

	// Reviews endpoints
	adminReviewsGroup := e.Group("/admin/reviews", middlewares.RequireOneOfRoles("admin"))
	{
		adminReviewsGroup.GET("", productHandler.ListReviews)
		adminReviewsGroup.DELETE("/:id", productHandler.DeleteReview)
	}

	// Admin stats WebSocket endpoint
	adminGroup := e.Group("/admin", middlewares.RequireOneOfRoles("admin"))
	e.GET("/admin/activity-stats/ws", userHandler.HandleActivityStatsWS, middlewares.RequireOneOfRoles("admin"))

	// Admin-only endpoints
	adminCategoryGroup := adminGroup.Group("/categories")
	{
		adminCategoryGroup.POST("", productHandler.CreateCategory)
		adminCategoryGroup.PUT("/:id", productHandler.UpdateCategory)
		adminCategoryGroup.DELETE("/:id", productHandler.DeleteCategory)
	}
	adminProductsGroup := adminGroup.Group("/products")
	{
		adminProductsGroup.POST("", productHandler.CreateProduct)
		adminProductsGroup.PUT("/:id", productHandler.UpdateProduct)
		adminProductsGroup.DELETE("/:id", productHandler.DeleteProduct)
	}

	tasks.AddShutdownTask(
		func(ctx context.Context) error {
			return e.Shutdown(ctx)
		},
	)

	logger.Info("REST server started", slog.String("address", fmt.Sprintf(":%d", cfg.ServerHttpPort)))
	return e.Start(fmt.Sprintf(":%d", cfg.ServerHttpPort))
}

type CustomValidator struct {
	Validator *validator.Validate
}

func (cv *CustomValidator) Validate(i interface{}) error {
	if err := cv.Validator.Struct(i); err != nil {
		// Optionally, you could return the error to give each route more control over the status code
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	return nil
}
