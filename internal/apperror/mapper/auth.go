package mapper

import (
	"front-office/common/constant"
	"front-office/internal/apperror"
	"strings"
)

func MapAuthError(err *apperror.ExternalAPIError) error {
	if err.StatusCode == 401 {
		if strings.Contains(err.Message, "not Active") {
			return apperror.Unauthorized("your account is not active")
		}

		return apperror.Unauthorized(constant.InvalidEmailOrPassword)
	}

	return MapExternalAPIError(err)
}
