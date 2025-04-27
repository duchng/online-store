package user

import (
	"context"
	"iter"

	"store-management/internal/core/product"
)

type UsePersistencePort interface {
	// Create saves a new user to the data store and returns the created user
	Create(ctx context.Context, user User) (User, error)
	// GetByID retrieves a user by their ID from the data store and returns the user
	GetByID(ctx context.Context, id int) (User, error)
	// GetByIDWithLock retrieves a user by their ID with a lock to prevent concurrent modifications
	GetByIDWithLock(ctx context.Context, id int) (User, error)
	// GetByUserName retrieves a user by their username from the data store and returns the user
	GetByUserName(ctx context.Context, username string) (User, error)
	// ListUsers retrieves filtered users from the data store
	ListUsers(ctx context.Context, filters UserFilters) ([]User, error)
	// UpdateRole updates a user's role
	UpdateRole(ctx context.Context, userId int, role UserRole) error
	// AddProductToWishList adds a product to the user's wishlist
	AddProductToWishList(ctx context.Context, userId int, productId int) error
	// RemoveProductFromWishList removes a product from the user's wishlist
	RemoveProductFromWishList(ctx context.Context, userId int, productId int) error
	// GetWishList retrieves the user's wishlist and returns a slice of products
	GetWishList(ctx context.Context, userId int) ([]product.Product, error)
	// UpdatePassword updates the user's password
	UpdatePassword(ctx context.Context, userId int, newPasswordHash string) error
	// GetActivityStats retrieves activity statistics for the user, support realtime update
	GetActivityStats(ctx context.Context) (map[ActivityKey]int, iter.Seq[UpdatedStat], error)
}
