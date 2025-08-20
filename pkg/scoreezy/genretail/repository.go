package genretail

import (
	"bytes"
	"encoding/json"
	"fmt"
	"front-office/app/config"
	"front-office/common/constant"
	"front-office/common/model"
	"front-office/helper"
	"front-office/internal/httpclient"
	"front-office/internal/jsonutil"
	"net/http"
)

func NewRepository(cfg *config.Config, client httpclient.HTTPClient, marshalFn jsonutil.Marshaller) Repository {
	if marshalFn == nil {
		marshalFn = json.Marshal // default behavior
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
	GenRetailV3API(memberId string, payload *genRetailRequest) (*model.ScoreezyAPIResponse[dataGenRetailV3], error)
	GetLogsScoreezyAPI(companyId string) (*model.AifcoreAPIResponse[[]*logTransScoreezy], error)
	GetLogsByRangeDateAPI(filter *filterLogs) (*model.AifcoreAPIResponse[[]*logTransScoreezy], error)
	// StoreImportData(newData []*BulkSearch, userId string) error
	// GetAllBulkSearch(tierLevel uint, userId, companyId string) ([]*BulkSearch, error)
	// CountData(tierLevel uint, userId, companyId string) (int64, error)
}

func (repo *repository) GenRetailV3API(memberId string, payload *genRetailRequest) (*model.ScoreezyAPIResponse[dataGenRetailV3], error) {
	url := fmt.Sprintf("%s/api/score/genretail/v3", repo.cfg.Env.ScoreezyHost)

	bodyBytes, err := repo.marshalFn(payload)
	if err != nil {
		return nil, fmt.Errorf(constant.ErrMsgMarshalReqBody, err)
	}

	req, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(bodyBytes))
	if err != nil {
		return nil, fmt.Errorf(constant.ErrMsgHTTPReqFailed, err)
	}

	req.Header.Set(constant.HeaderContentType, constant.HeaderApplicationJSON)
	req.Header.Set(constant.XUIDKey, memberId)

	res, err := repo.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf(constant.ErrMsgHTTPReqFailed, err)
	}
	defer res.Body.Close()

	return helper.ParseScoreezyAPIResponse[dataGenRetailV3](res)
}

func (repo *repository) fetchLogsAPI(path string, query map[string]string) (*model.AifcoreAPIResponse[[]*logTransScoreezy], error) {
	url := fmt.Sprintf("%s%s", repo.cfg.Env.AifcoreHost, path)

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf(constant.ErrMsgHTTPReqFailed, err)
	}
	req.Header.Set(constant.HeaderContentType, constant.HeaderApplicationJSON)

	q := req.URL.Query()
	for k, v := range query {
		q.Add(k, v)
	}
	req.URL.RawQuery = q.Encode()

	resp, err := repo.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf(constant.ErrMsgHTTPReqFailed, err)
	}
	defer resp.Body.Close()

	apiResp, err := helper.ParseAifcoreAPIResponse[[]*logTransScoreezy](resp)
	if err != nil {
		return nil, err
	}

	return apiResp, nil
}

func (repo *repository) GetLogsScoreezyAPI(companyId string) (*model.AifcoreAPIResponse[[]*logTransScoreezy], error) {
	return repo.fetchLogsAPI("/api/core/logging/transaction/scoreezy/list", map[string]string{"company_id": companyId})
}

func (repo *repository) GetLogsByRangeDateAPI(filter *filterLogs) (*model.AifcoreAPIResponse[[]*logTransScoreezy], error) {
	return repo.fetchLogsAPI("/api/core/logging/transaction/scoreezy/range", map[string]string{
		"company_id": filter.CompanyId,
		"date_start": filter.StartDate,
		"date_end":   filter.EndDate,
		"grade":      filter.Grade,
	})
}

// func (repo *repository) StoreImportData(newData []*BulkSearch, userId string) error {
// 	errTx := repo.DB.Transaction(func(tx *gorm.DB) error {
// 		// remove data existing in table
// 		if err := repo.DB.Delete(&BulkSearch{}, "user_id = ?", userId).Error; err != nil {
// 			return err
// 		}

// 		// replace existing with new data
// 		if err := tx.Create(&newData).Error; err != nil {
// 			return err
// 		}

// 		return nil
// 	})

// 	if errTx != nil {
// 		return errTx
// 	}

// 	return nil
// }

// func (repo *repository) GetAllBulkSearch(tierLevel uint, userId, companyId string) ([]*BulkSearch, error) {
// 	var bulkSearches []*BulkSearch

// 	query := repo.DB.Preload("User")

// 	if tierLevel == 1 {
// 		// admin
// 		query = query.Where("company_id = ?", companyId)
// 	} else {
// 		// user
// 		query = query.Where("user_id = ?", userId)
// 	}

// 	err := query.Find(&bulkSearches)

// 	if err.Error != nil {
// 		return nil, err.Error
// 	}

// 	return bulkSearches, nil
// }

// func (repo *repository) CountData(tierLevel uint, userId, companyId string) (int64, error) {
// 	var bulkSearches []*BulkSearch
// 	var count int64

// 	query := repo.DB.Debug()

// 	if tierLevel == 1 {
// 		// admin
// 		query = query.Where("company_id = ?", companyId)
// 	} else {
// 		// user
// 		query = query.Where("user_id = ?", userId)
// 	}

// 	err := query.Find(&bulkSearches).Count(&count).Error

// 	return count, err
// }
