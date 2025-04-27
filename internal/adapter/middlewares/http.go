package middlewares

import (
	"errors"
	"log/slog"
	"net/http"
	"strings"

	"github.com/labstack/echo/v4"

	"store-management/pkg/appcontext"
	"store-management/pkg/apperrors"
	"store-management/pkg/jwttoken"
)

// AuthenticationMiddleware provides JWT-based authentication for Echo requests by validating tokens from request headers.
// It extracts and parses the "Authorization" header, verifies token validity, and sets user claims in the request context.
// Requires a SignParser implementation for token parsing and validation. Returns an Echo middleware function.
func AuthenticationMiddleware(parser jwttoken.SignParser) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			// Get the public key from the request header
			authHeader := c.Request().Header.Get("Authorization")
			if authHeader != "" {
				headerParts := strings.Split(authHeader, " ")
				if len(headerParts) == 2 && headerParts[0] == "Bearer" {
					token, err := parser.ParseClaims(headerParts[1])
					if err != nil {
						return echo.NewHTTPError(http.StatusUnauthorized, "Invalid authentication provided")
					}
					if claims, ok := token.Claims.(*jwttoken.Claims); ok && token.Valid {
						c.Set(string(appcontext.AuthenticatedUser), claims)
					} else {
						return echo.NewHTTPError(http.StatusUnauthorized, "Invalid authentication token")
					}
				}
			}
			return next(c)
		}
	}
}

// RequireOneOfRoles checks if the user has at least one of the specified roles
func RequireOneOfRoles(roles ...string) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			authenticatedUser, ok := c.Get(string(appcontext.AuthenticatedUser)).(*jwttoken.Claims)
			if !ok || authenticatedUser == nil {
				return echo.NewHTTPError(http.StatusUnauthorized, "Unauthorized")
			}
			if !matchRole(authenticatedUser.Role, roles) {
				return echo.NewHTTPError(http.StatusUnauthorized, "You are not allowed to perform this action")
			}
			return next(c)
		}
	}
}

func matchRole(userRole string, allowedRoles []string) bool {
	for _, allowedRole := range allowedRoles {
		if strings.EqualFold(userRole, allowedRole) {
			return true
		}
	}
	return false
}

// CustomHTTPErrorHandler customize Echo's HTTP error handler to use custom Error message
// It sends a JSON response with status code.
func CustomHTTPErrorHandler(err error, ec echo.Context) {
	slog.ErrorContext(ec.Request().Context(), err.Error())
	var echoError *echo.HTTPError
	if errors.As(err, &echoError) {
		_ = ec.JSON(echoError.Code, echoError.Message)
		return
	}
	var appError apperrors.Error
	if errors.As(err, &appError) {
		scopedErr := ec.JSON(appError.Status, appError)
		if scopedErr != nil {
			ec.Echo().Logger.Error(err)
		}
	} else {
		// unknown error
		defaultErr := &apperrors.Error{
			Status:  http.StatusInternalServerError,
			Message: err.Error(),
		}
		scopedErr := ec.JSON(http.StatusInternalServerError, defaultErr)
		if scopedErr != nil {
			ec.Echo().Logger.Error(err)
		}
		return
	}
}
