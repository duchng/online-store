package user

import (
	"time"

	"github.com/golang-jwt/jwt/v5"

	"store-management/pkg/jwttoken"
)

// User represents a user in the system
type User struct {
	ID           int       `json:"id"`
	Username     string    `json:"username"`
	Email        string    `json:"email"`
	PasswordHash string    `json:"-"`
	FullName     string    `json:"fullName"`
	Role         UserRole  `json:"role"`
	CreatedAt    time.Time `json:"createdAt"`
	UpdatedAt    time.Time `json:"updatedAt"`
}

type UserRole string

const (
	UserRoleUser  UserRole = "user"
	UserRoleAdmin UserRole = "admin"
)

func (e UserRole) String() string {
	return string(e)
}

func (u User) GetSalt() []byte {
	return []byte(u.Username)
}

func (u User) AccessTokenClaims() jwttoken.Claims {
	return jwttoken.Claims{
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    jwttoken.Issuer,
			Subject:   u.Username,
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(jwttoken.DefaultAccessTokenLifetime)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
		FullName: u.FullName,
		UserId:   u.ID,
		Email:    u.Email,
		Role:     u.Role.String(),
	}
}

// Wishlist represents a user's wishlist item
type Wishlist struct {
	UserID    int       `json:"userId"`
	ProductID int       `json:"productId"`
	AddedAt   time.Time `json:"addedAt"`
}

type Token struct {
	AccessToken string `json:"accessToken"`
}

type UserFilters struct {
	Search string   `query:"search"`
	Roles  []string `query:"roles"`
}

type ActivityKey string

const (
	ActivitykeyWishListItems  ActivityKey = "wishListItems"
	ActivityKeyProductReviews ActivityKey = "productReviews"
	KeyStatUpdate                         = "statUpdate"
)

type UpdatedStat struct {
	Key   ActivityKey `json:"key"`
	Value int         `json:"value"`
}
