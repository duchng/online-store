package user

import (
	"context"
	"encoding/base64"
	"fmt"
	"iter"

	defined_errors "store-management/internal/core/defined-errors"
	"store-management/internal/core/product"
	"store-management/pkg/atomicity"
	"store-management/pkg/jwttoken"
	"store-management/pkg/password"
)

type UseCase interface {
	SignIn(ctx context.Context, userName, inputPassword string) (Token, error)
	SignUp(ctx context.Context, user User, inputPassword string) (User, error)
	AddToWishList(ctx context.Context, userId int, productId int) error
	RemoveFromWishList(ctx context.Context, userId int, productId int) error
	GetWishList(ctx context.Context, userId int) ([]product.Product, error)
	GetProfile(ctx context.Context, userId int) (User, error)
	ChangePassword(ctx context.Context, userId int, currentPassword, newPassword string) error
	ListUsers(ctx context.Context, filters UserFilters) ([]User, error)
	UpdateUserRole(ctx context.Context, userId int, role UserRole) error
	GetActivityStats(ctx context.Context) (map[ActivityKey]int, iter.Seq[UpdatedStat], error)
}

type usecase struct {
	persistencePort UsePersistencePort
	signParser      jwttoken.SignParser
	atomicExecutor  atomicity.AtomicExecutor
}

func NewUseCase(persistencePort UsePersistencePort, signParser jwttoken.SignParser, atomicExecutor atomicity.AtomicExecutor) UseCase {
	return &usecase{
		persistencePort: persistencePort,
		signParser:      signParser,
		atomicExecutor:  atomicExecutor,
	}
}

func (u *usecase) SignIn(ctx context.Context, userName, inputPassword string) (Token, error) {
	user, err := u.persistencePort.GetByUserName(ctx, userName)
	if err != nil {
		return Token{}, fmt.Errorf("usecase.SignIn %w", err)
	}
	match, err := password.Compare(inputPassword, user.PasswordHash, user.GetSalt())
	if err != nil {
		return Token{}, fmt.Errorf("usecase.SignIn %w", err)
	}
	if !match {
		return Token{}, fmt.Errorf("usecase.SignIn %w", defined_errors.ErrIncorrectPassword)
	}
	claims := user.AccessTokenClaims()
	tokenString, err := u.signParser.SignClaims(claims)
	if err != nil {
		return Token{}, fmt.Errorf("usecase.SignIn %w", err)
	}

	return Token{
		AccessToken: tokenString,
	}, nil
}

func (u *usecase) SignUp(ctx context.Context, user User, inputPassword string) (User, error) {
	user.PasswordHash = base64.RawStdEncoding.EncodeToString(password.HashPassword(inputPassword, user.GetSalt()))
	created, err := u.persistencePort.Create(ctx, user)
	if err != nil {
		return User{}, fmt.Errorf("usecase.SignUp %w", err)
	}
	return created, nil
}

func (u *usecase) AddToWishList(ctx context.Context, userId int, productId int) error {
	if err := u.persistencePort.AddProductToWishList(ctx, userId, productId); err != nil {
		return fmt.Errorf("usecase.AddToWishList %w", err)
	}
	return nil
}

func (u *usecase) RemoveFromWishList(ctx context.Context, userId int, productId int) error {
	if err := u.persistencePort.RemoveProductFromWishList(ctx, userId, productId); err != nil {
		return fmt.Errorf("usecase.RemoveFromWishList %w", err)
	}
	return nil
}

func (u *usecase) GetWishList(ctx context.Context, userId int) ([]product.Product, error) {
	products, err := u.persistencePort.GetWishList(ctx, userId)
	if err != nil {
		return nil, fmt.Errorf("usecase.GetWishList %w", err)
	}
	return products, nil
}

func (u *usecase) GetProfile(ctx context.Context, userId int) (User, error) {
	user, err := u.persistencePort.GetByID(ctx, userId)
	if err != nil {
		return User{}, fmt.Errorf("usecase.GetProfile %w", err)
	}
	return user, nil
}

func (u *usecase) ChangePassword(ctx context.Context, userId int, currentPassword, newPassword string) error {
	atomicFunc := func(ctx context.Context) error {
		user, err := u.persistencePort.GetByIDWithLock(ctx, userId)
		if err != nil {
			return err
		}
		match, err := password.Compare(currentPassword, user.PasswordHash, user.GetSalt())
		if err != nil {
			return err
		}
		if !match {
			return defined_errors.ErrPasswordMismath
		}
		newPasswordHash := base64.RawStdEncoding.EncodeToString(password.HashPassword(newPassword, user.GetSalt()))
		return u.persistencePort.UpdatePassword(ctx, userId, newPasswordHash)
	}
	err := u.atomicExecutor.Execute(ctx, atomicFunc)
	if err != nil {
		return fmt.Errorf("usecase.ChangePassword %w", err)
	}

	return nil
}

func (u *usecase) ListUsers(ctx context.Context, filters UserFilters) ([]User, error) {
	users, err := u.persistencePort.ListUsers(ctx, filters)
	if err != nil {
		return nil, fmt.Errorf("usecase.ListUsers %w", err)
	}
	return users, nil
}

func (u *usecase) UpdateUserRole(ctx context.Context, userId int, role UserRole) error {
	err := u.persistencePort.UpdateRole(ctx, userId, role)
	if err != nil {
		return fmt.Errorf("usecase.UpdateUserRole %w", err)
	}
	return nil
}

func (u *usecase) GetActivityStats(ctx context.Context) (map[ActivityKey]int, iter.Seq[UpdatedStat], error) {
	stats, iterator, err := u.persistencePort.GetActivityStats(ctx)
	if err != nil {
		return nil, nil, fmt.Errorf("usecase.GetActivityStats %w", err)
	}
	return stats, iterator, nil
}

var _ UseCase = (*usecase)(nil)
