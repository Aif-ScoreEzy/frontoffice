package mapper

import (
	"errors"
	"front-office/internal/apperror"
)

func MapExternalAPIError(err *apperror.ExternalAPIError) error {
	switch err.StatusCode {
	case 400:
		return apperror.BadRequest(err.Message)
	case 401:
		return apperror.Unauthorized("unauthorized")
	case 403:
		return apperror.Forbidden("access forbidden")
	case 404:
		return apperror.NotFound("resource not found")
	case 409:
		return apperror.Conflict(err.Message)
	case 422:
		return apperror.UnprocessableEntity(err.Message)
	case 429:
		return apperror.TooManyRequests("too many requests")
	case 500:
		return apperror.Internal("internal server error", nil)
	case 502:
		return apperror.BadGateway("bad gateway from external service")
	case 503:
		return apperror.ServiceUnavailable("external service unavailable")
	case 504:
		return apperror.GatewayTimeout("external service timeout")
	default:
		return apperror.BadGateway("unexpected external service error")
	}
}

func MapRepoError(err error, context string) error {
	var apiErr *apperror.ExternalAPIError
	if errors.As(err, &apiErr) {
		return MapExternalAPIError(apiErr)
	}

	return apperror.Internal(context, err)
}
