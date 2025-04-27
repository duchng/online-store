package redis

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"iter"
	"time"

	"github.com/redis/go-redis/v9"

	"store-management/internal/adapter/user/postgres"
	userCore "store-management/internal/core/user"
	"store-management/pkg/apperrors"
)

type UserRedisAdapter struct {
	*postgres.UserPostgresAdapter
	rdb *redis.Client
}

func NewUserRedisAdapter(userPostgresAdapter *postgres.UserPostgresAdapter, rdb *redis.Client) *UserRedisAdapter {
	return &UserRedisAdapter{
		UserPostgresAdapter: userPostgresAdapter,
		rdb:                 rdb,
	}
}

func (u *UserRedisAdapter) GetByID(ctx context.Context, id int) (userCore.User, error) {
	key := fmt.Sprintf("user:%d", id)
	cacheTime := time.Minute
	user := userCore.User{}
	bytes, err := u.rdb.Get(ctx, key).Result()
	if err == nil {
		_ = json.Unmarshal([]byte(bytes), &user)
		return user, nil
	}
	if !errors.Is(err, redis.Nil) {
		return userCore.User{}, fmt.Errorf("UserRedisAdapter.GetByID: %w", apperrors.FromError(err))
	}
	user, err = u.UserPostgresAdapter.GetByID(ctx, id)
	if err != nil {
		return userCore.User{}, err
	}
	userJSON, err := json.Marshal(user)
	if err != nil {
		return userCore.User{}, fmt.Errorf("UserRedisAdapter.GetByID: %w", apperrors.FromError(err))
	}
	err = u.rdb.Set(ctx, key, userJSON, cacheTime).Err()
	if err != nil {
		return userCore.User{}, fmt.Errorf("UserRedisAdapter.GetByID: %w", apperrors.FromError(err))
	}
	return user, nil
}

func (u *UserRedisAdapter) UpdateRole(ctx context.Context, userId int, role userCore.UserRole) error {
	key := fmt.Sprintf("user:%d", userId)
	// invalidate the user
	err := u.rdb.Del(ctx, key).Err()
	if err != nil {
		return fmt.Errorf("UserRedisAdapter.UpdateRole: %w", apperrors.FromError(err))
	}
	return u.UserPostgresAdapter.UpdateRole(ctx, userId, role)
}

func (u *UserRedisAdapter) RemoveProductFromWishList(ctx context.Context, userId int, productId int) error {
	if err := u.UserPostgresAdapter.RemoveProductFromWishList(ctx, userId, productId); err != nil {
		return err
	}
	val, err := u.rdb.Decr(ctx, string(userCore.ActivitykeyWishListItems)).Uint64()
	if err != nil {
		return fmt.Errorf("UserRedisAdapter.RemoveProductFromWishList: %w", apperrors.FromError(err))
	}
	updateState, _ := json.Marshal(
		userCore.UpdatedStat{
			Key:   userCore.ActivitykeyWishListItems,
			Value: int(val),
		},
	)
	// notify stat update
	err = u.rdb.Publish(ctx, userCore.KeyStatUpdate, updateState).Err()
	if err != nil {
		return fmt.Errorf("UserRedisAdapter.RemoveProductFromWishList: %w", apperrors.FromError(err))
	}
	return nil
}

func (u *UserRedisAdapter) AddProductToWishList(ctx context.Context, userId int, productId int) error {
	if err := u.UserPostgresAdapter.AddProductToWishList(ctx, userId, productId); err != nil {
		return err
	}
	val, err := u.rdb.Incr(ctx, string(userCore.ActivitykeyWishListItems)).Uint64()
	if err != nil {
		return fmt.Errorf("UserRedisAdapter.AddProductToWishList: %w", apperrors.FromError(err))
	}
	updateState, _ := json.Marshal(
		userCore.UpdatedStat{
			Key:   userCore.ActivitykeyWishListItems,
			Value: int(val),
		},
	)
	// notify stat update
	err = u.rdb.Publish(ctx, userCore.KeyStatUpdate, updateState).Err()
	if err != nil {
		return fmt.Errorf("UserRedisAdapter.AddProductToWishList: %w", apperrors.FromError(err))
	}
	return nil
}

func (u *UserRedisAdapter) GetActivityStats(ctx context.Context) (map[userCore.ActivityKey]int, iter.Seq[userCore.UpdatedStat], error) {
	refetch := false
	wishListStat, err := u.rdb.Get(ctx, string(userCore.ActivitykeyWishListItems)).Int()
	if err != nil && !errors.Is(err, redis.Nil) {
		return nil, nil, fmt.Errorf("UserRedisAdapter.GetActivityStats: %w", apperrors.FromError(err))
	}
	if errors.Is(err, redis.Nil) {
		refetch = true
	}
	reviewsStat, err := u.rdb.Get(ctx, string(userCore.ActivityKeyProductReviews)).Int()
	if err != nil && !errors.Is(err, redis.Nil) {
		return nil, nil, fmt.Errorf("UserRedisAdapter.GetActivityStats: %w", apperrors.FromError(err))
	}
	if errors.Is(err, redis.Nil) {
		refetch = true
	}
	if refetch {
		stats, _, err := u.UserPostgresAdapter.GetActivityStats(ctx)
		if err != nil {
			return nil, nil, err
		}
		for k, v := range stats {
			err = u.rdb.Set(ctx, string(k), v, time.Hour).Err()
			if err != nil {
				return nil, nil, fmt.Errorf("UserRedisAdapter.GetActivityStats: %w", apperrors.FromError(err))
			}
		}
		wishListStat = stats[userCore.ActivitykeyWishListItems]
		reviewsStat = stats[userCore.ActivityKeyProductReviews]
	}
	sub := u.rdb.Subscribe(ctx, userCore.KeyStatUpdate)
	iterator := func(yield func(stat userCore.UpdatedStat) bool) {
		for {
			select {
			case <-ctx.Done():
				return
			case msg := <-sub.Channel():
				var stat userCore.UpdatedStat
				err = json.Unmarshal([]byte(msg.Payload), &stat)
				if err != nil {
					continue
				}
				if !yield(stat) {
					break
				}
			}
		}
	}
	return map[userCore.ActivityKey]int{
		userCore.ActivitykeyWishListItems:  wishListStat,
		userCore.ActivityKeyProductReviews: reviewsStat,
	}, iterator, nil
}
