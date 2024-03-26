package passwordresettoken

import (
	"front-office/helper"
	"front-office/pkg/user"
	"os"
	"strconv"

	"github.com/google/uuid"
)

func NewService(repo Repository) Service {
	return &service{Repo: repo}
}

type service struct {
	Repo Repository
}

type Service interface {
	CreatePasswordResetTokenSvc(user *user.User) (string, *PasswordResetToken, error)
	FindPasswordResetTokenByTokenSvc(token string) (*PasswordResetToken, error)
}

func (svc *service) CreatePasswordResetTokenSvc(user *user.User) (string, *PasswordResetToken, error) {
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

	passwordResetToken, err = svc.Repo.CreatePasswordResetToken(passwordResetToken)
	if err != nil {
		return "", nil, err
	}

	return token, passwordResetToken, nil
}

func (svc *service) FindPasswordResetTokenByTokenSvc(token string) (*PasswordResetToken, error) {
	result, err := svc.Repo.FindOnePasswordResetTokenByToken(token)
	if err != nil {
		return nil, err
	}

	return result, nil
}
