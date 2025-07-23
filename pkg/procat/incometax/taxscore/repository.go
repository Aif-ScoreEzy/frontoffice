package taxscore

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"front-office/app/config"
	"front-office/common/constant"
	"front-office/common/model"
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
	TaxScoreAPI(apiKey, jobId string, request *taxScoreRequest) (*model.ProCatAPIResponse[taxScoreRespData], error)
}

func (repo *repository) TaxScoreAPI(apiKey, jobId string, reqBody *taxScoreRequest) (*model.ProCatAPIResponse[taxScoreRespData], error) {
	url := fmt.Sprintf("%s/product/incometax/tax-score", repo.cfg.Env.ProductCatalogHost)

	bodyBytes, err := repo.marshalFn(reqBody)
	if err != nil {
		return nil, fmt.Errorf(constant.ErrMsgMarshalReqBody, err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewBuffer(bodyBytes))
	if err != nil {
		return nil, err
	}

	req.Header.Set(constant.HeaderContentType, constant.HeaderApplicationJSON)
	req.Header.Set(constant.XAPIKey, apiKey)

	q := req.URL.Query()
	q.Add("job_id", jobId)
	req.URL.RawQuery = q.Encode()

	resp, err := repo.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf(constant.ErrMsgHTTPReqFailed, err)
	}
	defer resp.Body.Close()

	apiResp, err := helper.ParseProCatAPIResponse[taxScoreRespData](resp)
	if err != nil {
		return nil, err
	}

	return apiResp, err
}
