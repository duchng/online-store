package apperrors

import (
	"errors"
	"net/http"
	"strings"

	"github.com/jackc/pgerrcode"
	"github.com/lib/pq"
	"github.com/redis/go-redis/v9"
)

func IsConstraintViolationError(err error) bool {
	var pqErr *pq.Error
	if errors.As(err, &pqErr) {
		if pgerrcode.IsIntegrityConstraintViolation(string(pqErr.Code)) {
			return true
		}
	}
	return false
}

func IsNotFoundError(err error) bool {
	if err == nil {
		return false
	}
	var pqErr *pq.Error
	if errors.As(err, &pqErr) {
		if pgerrcode.IsCaseNotFound(string(pqErr.Code)) {
			return true
		}
	}
	if strings.Contains(err.Error(), "no rows in result set") {
		return true
	}
	var appErr Error
	if errors.As(err, &appErr) {
		if appErr.Status == http.StatusNotFound {
			return true
		}
	}
	if errors.Is(err, redis.Nil) {
		return true
	}
	return false
}

func IsObjectNotInPrerequisiteStateError(err error) bool {
	var pqErr *pq.Error
	if errors.As(err, &pqErr) {
		if pgerrcode.IsObjectNotInPrerequisiteState(string(pqErr.Code)) {
			return true
		}
	}
	return false
}
