package passwordresettoken

import (
	"bytes"
	"encoding/json"
	"fmt"
	"front-office/app/config"
	"front-office/common/constant"
	"net/http"

	"gorm.io/gorm"
)

func NewRepository(db *gorm.DB, cfg *config.Config) Repository {
	return &repository{DB: db, Cfg: cfg}
}

type repository struct {
	DB  *gorm.DB
	Cfg *config.Config
}

type Repository interface {
	FindOnePasswordResetTokenByToken(token string) (*PasswordResetToken, error)
	FindOnePasswordResetTokenByUserId(userId string) (*PasswordResetToken, error)
	CreatePasswordResetTokenAifCore(req *CreatePasswordResetTokenRequest, userId string) (*http.Response, error)
}

// func (repo *repository) CreatePasswordResetToken(passwordResetToken *PasswordResetToken) (*PasswordResetToken, error) {
// 	err := repo.DB.Create(&passwordResetToken).Error
// 	if err != nil {
// 		return nil, err
// 	}

// 	return passwordResetToken, nil
// }

func (repo *repository) FindOnePasswordResetTokenByToken(token string) (*PasswordResetToken, error) {
	var result *PasswordResetToken
	err := repo.DB.First(&result, "token = ?", token).Error
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (repo *repository) FindOnePasswordResetTokenByUserId(userId string) (*PasswordResetToken, error) {
	var passwordResetToken *PasswordResetToken

	err := repo.DB.First(&passwordResetToken, "user_id = ?", userId).Error
	if err != nil {
		return nil, err
	}

	return passwordResetToken, nil
}

func (repo *repository) CreatePasswordResetTokenAifCore(req *CreatePasswordResetTokenRequest, userId string) (*http.Response, error) {
	apiUrl := fmt.Sprintf(`%v/api/core/member/%v/password-reset-token`, repo.Cfg.Env.AifcoreHost, userId)

	jsonBodyValue, _ := json.Marshal(req)
	request, _ := http.NewRequest(http.MethodPost, apiUrl, bytes.NewBuffer(jsonBodyValue))
	request.Header.Set(constant.HeaderContentType, constant.HeaderApplicationJSON)

	client := &http.Client{}
	return client.Do(request)
}
