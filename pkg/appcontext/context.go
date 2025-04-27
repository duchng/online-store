package appcontext

import (
	"errors"

	"github.com/labstack/echo/v4"

	"store-management/pkg/jwttoken"
)

type Key string

const (
	AuthenticatedUser Key = "AuthenticatedUser"
)

func ContextGetUserData(ctx echo.Context) *jwttoken.Claims {
	if user, ok := ctx.Get(string(AuthenticatedUser)).(*jwttoken.Claims); ok {
		return user
	}
	return nil
}

func ContextGetUserId(ctx echo.Context) (int, error) {
	if user := ContextGetUserData(ctx); user != nil {
		return user.UserId, nil
	}
	return 0, errors.New("unauthenticated user")
}
