package postgres

import (
	"context"
	"fmt"
	"iter"
	"time"

	"github.com/uptrace/bun"

	productDB "store-management/internal/adapter/product/postgres"
	"store-management/internal/core/product"
	userCore "store-management/internal/core/user"
	"store-management/pkg/apperrors"
	"store-management/pkg/database"
)

type UserPostgresAdapter struct {
	GetDbFunc database.GetDbFunc
}

func (u *UserPostgresAdapter) UpdatePassword(ctx context.Context, userId int, newPasswordHash string) error {
	db := u.GetDbFunc(ctx)
	_, err := db.NewUpdate().
		Model((*User)(nil)).
		Set("password_hash = ?", newPasswordHash).
		Where("id = ?", userId).
		Exec(ctx)
	if err != nil {
		return fmt.Errorf("UserPostgresAdapter.UpdatePassword: %w", apperrors.FromError(err))
	}
	return nil
}

func NewUserPostgresAdapter(getDb database.GetDbFunc) *UserPostgresAdapter {
	getDb(context.Background()).(*bun.DB).RegisterModel((*Wishlist)(nil))
	return &UserPostgresAdapter{
		GetDbFunc: getDb,
	}
}

func (u *UserPostgresAdapter) AddProductToWishList(ctx context.Context, userId int, productId int) error {
	db := u.GetDbFunc(ctx)
	wishlist := Wishlist{
		UserID:    userId,
		ProductID: productId,
		AddedAt:   time.Now(),
	}
	_, err := db.NewInsert().Model(&wishlist).Exec(ctx)
	if err != nil {
		return fmt.Errorf("UserPostgresAdapter.AddProductToWishList: %w", apperrors.FromError(err))
	}
	return nil
}

func (u *UserPostgresAdapter) RemoveProductFromWishList(ctx context.Context, userId int, productId int) error {
	db := u.GetDbFunc(ctx)
	_, err := db.NewDelete().Model((*Wishlist)(nil)).
		Where("user_id = ? AND product_id = ?", userId, productId).
		Exec(ctx)
	if err != nil {
		return fmt.Errorf("UserPostgresAdapter.RemoveProductFromWishList: %w", apperrors.FromError(err))
	}
	return nil
}

func (u *UserPostgresAdapter) GetWishList(ctx context.Context, userId int) ([]product.Product, error) {
	db := u.GetDbFunc(ctx)
	var user User
	err := db.NewSelect().
		Model(&user).
		Where("id = ?", userId).
		Relation("WishList").
		Scan(ctx)
	if err != nil {
		return nil, fmt.Errorf("UserPostgresAdapter.GetWishList: %w", apperrors.FromError(err))
	}

	products := make([]product.Product, len(user.WishList))
	for i, dbProduct := range user.WishList {
		products[i] = productDB.MapProductDbToDomain(dbProduct)
	}
	return products, nil
}

func (u *UserPostgresAdapter) Create(ctx context.Context, user userCore.User) (userCore.User, error) {
	db := u.GetDbFunc(ctx)
	toCreate := MapUserDomainToDb(user)
	created := User{}
	_, err := db.NewInsert().Model(&toCreate).Returning("*").Exec(ctx, &created)
	if err != nil {
		return user, fmt.Errorf("UserPostgresAdapter.Create: %w", apperrors.FromError(err))
	}
	return MapUserDbToDomain(created), nil
}

func (u *UserPostgresAdapter) GetByID(ctx context.Context, id int) (userCore.User, error) {
	db := u.GetDbFunc(ctx)
	user := User{}
	err := db.NewSelect().Model(&user).Where("id = ?", id).Scan(ctx)
	if err != nil {
		return userCore.User{}, fmt.Errorf("UserPostgresAdapter.GetByID: %w", apperrors.FromError(err))
	}
	return MapUserDbToDomain(user), nil
}

func (u *UserPostgresAdapter) GetByIDWithLock(ctx context.Context, id int) (userCore.User, error) {
	db := u.GetDbFunc(ctx)
	user := User{}
	err := db.NewSelect().Model(&user).Where("id = ?", id).For("UPDATE").Scan(ctx)
	if err != nil {
		return userCore.User{}, fmt.Errorf("UserPostgresAdapter.GetByIDWithLock: %w", apperrors.FromError(err))
	}
	return MapUserDbToDomain(user), nil
}

func (u *UserPostgresAdapter) GetByUserName(ctx context.Context, username string) (userCore.User, error) {
	db := u.GetDbFunc(ctx)
	user := User{}
	err := db.NewSelect().Model(&user).Where("username = ?", username).Scan(ctx)
	if err != nil {
		return userCore.User{}, fmt.Errorf("UserPostgresAdapter.GetByUserName: %w", apperrors.FromError(err))
	}
	return MapUserDbToDomain(user), nil
}

func (u *UserPostgresAdapter) ListUsers(ctx context.Context, filters userCore.UserFilters) ([]userCore.User, error) {
	db := u.GetDbFunc(ctx)
	var users []User
	query := db.NewSelect().Model(&users).Order("created_at DESC")

	// Apply search filter
	if filters.Search != "" {
		query.WhereGroup(
			" AND ", func(q *bun.SelectQuery) *bun.SelectQuery {
				return q.Where("username ILIKE ?", "%"+filters.Search+"%").
					WhereOr("email ILIKE ?", "%"+filters.Search+"%")
			},
		)
	}

	// Apply role filter
	if len(filters.Roles) > 0 {
		query.Where("role IN (?)", bun.In(filters.Roles))
	}

	err := query.Scan(ctx)
	if err != nil {
		return nil, fmt.Errorf("UserPostgresAdapter.ListUsers: %w", err)
	}

	result := make([]userCore.User, len(users))
	for i, user := range users {
		result[i] = MapUserDbToDomain(user)
	}
	return result, nil
}

func (u *UserPostgresAdapter) UpdateRole(ctx context.Context, userId int, role userCore.UserRole) error {
	db := u.GetDbFunc(ctx)
	_, err := db.NewUpdate().
		Model((*User)(nil)).
		Set("role = ?", role).
		Where("id = ?", userId).
		Exec(ctx)
	if err != nil {
		return fmt.Errorf("UserPostgresAdapter.UpdateRole: %w", apperrors.FromError(err))
	}
	return nil
}

func (u *UserPostgresAdapter) GetActivityStats(ctx context.Context) (map[userCore.ActivityKey]int, iter.Seq[userCore.UpdatedStat], error) {
	count := 0
	res := make(map[userCore.ActivityKey]int)
	err := u.GetDbFunc(ctx).NewSelect().ColumnExpr("COUNT(*) as count").Table("reviews").Scan(ctx, &count)
	if err != nil {
		return nil, nil, fmt.Errorf("UserPostgresAdapter.GetActivityStats: %w", apperrors.FromError(err))
	}
	res[userCore.ActivityKeyProductReviews] = count
	err = u.GetDbFunc(ctx).NewSelect().ColumnExpr("COUNT(*) as count").Table("wishlist").Scan(ctx, &count)
	if err != nil {
		return nil, nil, fmt.Errorf("UserPostgresAdapter.GetActivityStats: %w", apperrors.FromError(err))
	}
	res[userCore.ActivitykeyWishListItems] = count
	return res, nil, nil
}

var _ userCore.UsePersistencePort = (*UserPostgresAdapter)(nil)
