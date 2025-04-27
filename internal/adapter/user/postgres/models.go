package postgres

import (
	"time"

	"github.com/uptrace/bun"

	productDB "store-management/internal/adapter/product/postgres"
	userCore "store-management/internal/core/user"
)

type User struct {
	bun.BaseModel `bun:"table:users"`

	ID           int `bun:",pk,nullzero"`
	Username     string
	Email        string
	PasswordHash string
	FullName     string
	Role         userCore.UserRole
	CreatedAt    time.Time `bun:",nullzero"`
	UpdatedAt    time.Time `bun:",nullzero"`

	WishList []productDB.Product `bun:"m2m:wishlist,join:User=Product"`
}

type Wishlist struct {
	bun.BaseModel `bun:"table:wishlist"`

	UserID    int `bun:",pk"`
	ProductID int `bun:",pk"`
	AddedAt   time.Time

	User    *User              `bun:"rel:belongs-to,join:user_id=id"`
	Product *productDB.Product `bun:"rel:belongs-to,join:product_id=id"`
}
