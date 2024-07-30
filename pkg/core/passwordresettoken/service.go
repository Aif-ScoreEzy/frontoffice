package passwordresettoken

import (
	"front-office/app/config"
)

func NewService(repo Repository, cfg *config.Config) Service {
	return &service{Repo: repo, Cfg: cfg}
}

type service struct {
	Repo Repository
	Cfg  *config.Config
}

type Service interface {
	// CreatePasswordResetTokenSvc(user *user.User) (string, *PasswordResetToken, error)
	FindPasswordResetTokenByTokenSvc(token string) (*PasswordResetToken, error)
}

// func (svc *service) CreatePasswordResetTokenSvc(user *user.User) (string, *PasswordResetToken, error) {
// 	secret := svc.Cfg.Env.JwtSecretKey
// 	minutesToExpired, _ := strconv.Atoi(svc.Cfg.Env.JwtResetPasswordExpiresMinutes)

// 	token, err := helper.GenerateToken(secret, minutesToExpired, user.Id, user.CompanyId, user.Role.TierLevel)
// 	if err != nil {
// 		return "", nil, err
// 	}

// 	tokenId := uuid.NewString()
// 	passwordResetToken := &PasswordResetToken{
// 		Id:     tokenId,
// 		Token:  token,
// 		UserId: user.Id,
// 	}

// 	passwordResetToken, err = svc.Repo.CreatePasswordResetToken(passwordResetToken)
// 	if err != nil {
// 		return "", nil, err
// 	}

// 	return token, passwordResetToken, nil
// }

func (svc *service) FindPasswordResetTokenByTokenSvc(token string) (*PasswordResetToken, error) {
	result, err := svc.Repo.FindOnePasswordResetTokenByToken(token)
	if err != nil {
		return nil, err
	}

	return result, nil
}
