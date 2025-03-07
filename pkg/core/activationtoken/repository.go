package activationtoken

import (
	"bytes"
	"encoding/json"
	"fmt"
	"front-office/app/config"
	"front-office/common/constant"
	"net/http"
)

func NewRepository(cfg *config.Config) Repository {
	return &repository{Cfg: cfg}
}

type repository struct {
	Cfg *config.Config
}

type Repository interface {
	FindOneActivationTokenBytoken(token string) (*http.Response, error)
	CreateActivationTokenAifCore(req *CreateActivationTokenRequest, memberId string) (*http.Response, error)
}

func (repo *repository) FindOneActivationTokenBytoken(token string) (*http.Response, error) {
	apiUrl := fmt.Sprintf(`%v/api/core/member/activation-tokens/%v`, repo.Cfg.Env.AifcoreHost, token)

	request, _ := http.NewRequest(http.MethodGet, apiUrl, nil)
	request.Header.Set(constant.HeaderContentType, constant.HeaderApplicationJSON)

	client := &http.Client{}
	return client.Do(request)
}

func (repo *repository) CreateActivationTokenAifCore(req *CreateActivationTokenRequest, memberId string) (*http.Response, error) {
	apiUrl := fmt.Sprintf(`%v/api/core/member/%v/activation-tokens`, repo.Cfg.Env.AifcoreHost, memberId)

	jsonBodyValue, _ := json.Marshal(req)
	request, _ := http.NewRequest(http.MethodPost, apiUrl, bytes.NewBuffer(jsonBodyValue))
	request.Header.Set(constant.HeaderContentType, constant.HeaderApplicationJSON)

	client := &http.Client{}
	return client.Do(request)
}
