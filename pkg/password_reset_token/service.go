package password_reset_token

import (
	"front-office/helper"
	"front-office/pkg/user"
	"os"
	"strconv"

	"github.com/google/uuid"
)

func CreatePasswordResetTokenSvc(user *user.User) (string, *PasswordResetToken, error) {
	secret := os.Getenv("JWT_SECRET_KEY")
	minutesToExpired, _ := strconv.Atoi(os.Getenv("JWT_RESET_PASSWORD_EXPIRES_MINUTES"))

	token, err := helper.GenerateToken(secret, minutesToExpired, user.ID, user.CompanyID, user.Role.TierLevel)
	if err != nil {
		return "", nil, err
	}

	tokenID := uuid.NewString()
	passwordResetToken := &PasswordResetToken{
		ID:     tokenID,
		Token:  token,
		UserID: user.ID,
	}

	passwordResetToken, err = CreatePasswordResetToken(passwordResetToken)
	if err != nil {
		return "", nil, err
	}

	return token, passwordResetToken, nil
}

func FindPasswordResetTokenByTokenSvc(token string) (*PasswordResetToken, error) {
	result, err := FindOnePasswordResetTokenByToken(token)
	if err != nil {
		return nil, err
	}

	return result, nil
}
