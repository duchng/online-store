package apperrors

import (
	"net/http"

	"google.golang.org/grpc/codes"
)

// GRPCStatusToHTTPStatus converts a gRPC status code to an HTTP status code.
func GRPCStatusToHTTPStatus(code codes.Code) int {
	switch code {
	case codes.OK:
		return http.StatusOK

	case codes.Canceled:
		return http.StatusRequestTimeout

	case codes.Unknown:
		return http.StatusInternalServerError

	case codes.InvalidArgument:
		return http.StatusBadRequest

	case codes.DeadlineExceeded:
		return http.StatusGatewayTimeout

	case codes.NotFound:
		return http.StatusNotFound

	case codes.AlreadyExists:
		return http.StatusConflict

	case codes.PermissionDenied:
		return http.StatusForbidden

	case codes.ResourceExhausted:
		return http.StatusTooManyRequests

	case codes.FailedPrecondition:
		return http.StatusPreconditionFailed

	case codes.Aborted:
		return http.StatusLocked

	case codes.OutOfRange:
		return http.StatusRequestedRangeNotSatisfiable

	case codes.Unimplemented:
		return http.StatusNotImplemented

	case codes.Internal:
		return http.StatusInternalServerError

	case codes.Unavailable:
		return http.StatusServiceUnavailable

	case codes.DataLoss:
		return http.StatusInternalServerError

	case codes.Unauthenticated:
		return http.StatusUnauthorized

	default:
		return http.StatusInternalServerError
	}
}

// HTTPStatusToGRPCStatus converts an HTTP status code to a gRPC status code.
func HTTPStatusToGRPCStatus(code int) codes.Code {
	switch code {
	case http.StatusOK: // 200
		return codes.OK

	case http.StatusBadRequest: // 400
		return codes.InvalidArgument

	case http.StatusUnauthorized: // 401
		return codes.Unauthenticated

	case http.StatusForbidden: // 403
		return codes.PermissionDenied

	case http.StatusNotFound: // 404
		return codes.NotFound

	case http.StatusRequestTimeout: // 408
		return codes.Canceled

	case http.StatusConflict: // 409
		return codes.AlreadyExists

	case http.StatusPreconditionFailed: // 412
		return codes.FailedPrecondition

	case http.StatusRequestedRangeNotSatisfiable: // 416
		return codes.OutOfRange

	case http.StatusExpectationFailed: // 417
		return codes.FailedPrecondition

	case http.StatusLocked: // 423
		return codes.Aborted

	case http.StatusFailedDependency:
		return codes.Aborted

	case http.StatusTooManyRequests: // 429
		return codes.ResourceExhausted

	case http.StatusInternalServerError: // 500
		return codes.Unknown

	case http.StatusNotImplemented: // 501
		return codes.Unimplemented

	case http.StatusServiceUnavailable: // 503
		return codes.Unavailable

	case http.StatusGatewayTimeout: // 504
		return codes.DeadlineExceeded

	default:
		return codes.Unknown
	}
}
