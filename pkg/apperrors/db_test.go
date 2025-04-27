package apperrors

import (
	"errors"
	"net/http"
	"testing"

	"github.com/jackc/pgerrcode"
	"github.com/lib/pq"
	"google.golang.org/grpc/codes"
)

func TestFromError(t *testing.T) {
	type testCase struct {
		err      error
		expected Error
	}
	testCases := []testCase{
		{
			err: errors.New("no rows in result set"),
			expected: Error{
				Status:  http.StatusNotFound,
				Code:    "NOT_FOUND",
				Message: "not found resources",
			},
		},
		{
			err: &pq.Error{
				Code: pgerrcode.CaseNotFound,
			},
			expected: Error{
				Status:  http.StatusNotFound,
				Code:    "NOT_FOUND",
				Message: "not found resources",
			},
		},
		{
			err: &pq.Error{
				Code: pgerrcode.ObjectInUse,
			},
			expected: Error{
				Status:  http.StatusBadRequest,
				Code:    "OBJECT_NOT_IN_PREREQUISITE_STATE",
				Message: "object not in prerequisite state",
			},
		},
		{
			err: errors.New("undefined"),
			expected: Error{
				Status:      500,
				Code:        codes.Unknown.String(),
				Message:     "undefined",
				Description: "undefined",
			},
		},
	}
	for _, tc := range testCases {
		wrappedError := FromError(tc.err)
		if !errors.Is(wrappedError, tc.expected) {
			t.Errorf("expected %v but got %v", tc.expected, wrappedError)
		}
	}
}
