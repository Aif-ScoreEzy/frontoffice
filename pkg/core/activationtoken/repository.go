package activationtoken

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
	FindOneActivationTokenBytoken(token string) (*MstActivationToken, error)
	FindOneActivationTokenByUserId(userId string) (*MstActivationToken, error)
	CreateActivationToken(activationToken *MstActivationToken) (*MstActivationToken, error)
	CreateActivationTokenAifCore(req *CreateActivationTokenRequest, userId string) (*http.Response, error)
}

func (repo *repository) FindOneActivationTokenBytoken(token string) (*MstActivationToken, error) {
	var activationToken *MstActivationToken

	err := repo.DB.First(&activationToken, "token = ?", token).Error
	if err != nil {
		return nil, err
	}

	return activationToken, nil
}

func (repo *repository) FindOneActivationTokenByUserId(userId string) (*MstActivationToken, error) {
	var activationToken *MstActivationToken

	err := repo.DB.First(&activationToken, "user_id = ?", userId).Error
	if err != nil {
		return nil, err
	}

	return activationToken, nil
}

func (repo *repository) CreateActivationToken(activationToken *MstActivationToken) (*MstActivationToken, error) {
	err := repo.DB.Create(&activationToken).Error
	if err != nil {
		return nil, err
	}

	return activationToken, nil
}

func (repo *repository) CreateActivationTokenAifCore(req *CreateActivationTokenRequest, userId string) (*http.Response, error) {
	apiUrl := fmt.Sprintf(`%v/api/core/member/%v/activation-token`, repo.Cfg.Env.AifcoreHost, userId)

	jsonBodyValue, _ := json.Marshal(req)
	request, _ := http.NewRequest(http.MethodPost, apiUrl, bytes.NewBuffer(jsonBodyValue))
	request.Header.Set(constant.HeaderContentType, constant.HeaderApplicationJSON)

	client := &http.Client{}
	return client.Do(request)
}
