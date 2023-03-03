package statustocode

import (
	"net/http"

	"google.golang.org/grpc/codes"
)

func GetCodeFromResponse(response *http.Response) codes.Code {
	if response == nil {
		return codes.Unknown
	}
	return codeFromHTTPStatus(response.StatusCode)
}

func GetStatusFromResponse(response *http.Response) int {
	if response == nil {
		return http.StatusNotFound
	}
	return response.StatusCode
}

func codeFromHTTPStatus(code int) codes.Code {
	switch code {
	case http.StatusBadRequest:
		return codes.InvalidArgument
	case http.StatusUnauthorized:
		return codes.Unauthenticated
	case http.StatusForbidden:
		return codes.PermissionDenied
	case http.StatusNotFound:
		return codes.NotFound
	case http.StatusConflict:
		return codes.AlreadyExists
	case http.StatusTooManyRequests:
		return codes.ResourceExhausted
	case http.StatusInternalServerError:
		return codes.Internal
	case http.StatusNotImplemented:
		return codes.Unimplemented
	case http.StatusServiceUnavailable:
		return codes.Unavailable
	case http.StatusGatewayTimeout:
		return codes.DeadlineExceeded
	case http.StatusOK:
		return codes.OK
	}
	return codes.Unknown
}
