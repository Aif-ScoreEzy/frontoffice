package grade

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"front-office/app/config"
	"front-office/common/constant"
	"front-office/helper"
	"front-office/internal/httpclient"
	"front-office/internal/jsonutil"
	"net/http"
	"time"
)

func NewRepository(cfg *config.Config, client httpclient.HTTPClient, marshalFn jsonutil.Marshaller) Repository {
	if marshalFn == nil {
		marshalFn = json.Marshal
	}

	return &repository{
		cfg:       cfg,
		client:    client,
		marshalFn: marshalFn,
	}
}

type repository struct {
	cfg       *config.Config
	client    httpclient.HTTPClient
	marshalFn jsonutil.Marshaller
}

type Repository interface {
	CreateGradesAPI(payload *createGradePayload) error
}

func (repo *repository) CreateGradesAPI(payload *createGradePayload) error {
	url := fmt.Sprintf("%s/api/core/product/%s/grades", repo.cfg.Env.AifcoreHost, payload.ProductSlug)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	bodyBytes, err := repo.marshalFn(payload.Request)
	if err != nil {
		return fmt.Errorf(constant.ErrMsgMarshalReqBody, err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewBuffer(bodyBytes))
	if err != nil {
		return fmt.Errorf(constant.ErrMsgHTTPReqFailed, err)
	}

	req.Header.Set(constant.HeaderContentType, constant.HeaderApplicationJSON)
	req.Header.Set(constant.XCompanyId, payload.CompanyId)

	resp, err := repo.client.Do(req)
	if err != nil {
		return fmt.Errorf(constant.ErrMsgHTTPReqFailed, err)
	}
	defer resp.Body.Close()

	_, err = helper.ParseAifcoreAPIResponse[any](resp)
	if err != nil {
		return err
	}

	return nil
}
