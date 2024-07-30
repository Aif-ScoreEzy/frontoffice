package genretail

import (
	"bytes"
	"encoding/json"
	"front-office/app/config"
	"front-office/common/constant"
	"io"
	"net/http"

	"github.com/google/uuid"
)

func NewService(repo Repository, cfg *config.Config) Service {
	return &service{Repo: repo, Cfg: cfg}
}

type service struct {
	Repo Repository
	Cfg  *config.Config
}

type Service interface {
	GenRetailV3(requestData *GenRetailRequest, apiKey string) (*GenRetailV3ModelResponse, error)
	BulkSearchUploadSvc(req []BulkSearchRequest, tempType, apiKey, userId, companyId string) error
	GetBulkSearchSvc(tierLevel uint, userId, companyId string) ([]BulkSearchResponse, error)
	GetTotalDataBulk(tierLevel uint, userId, companyId string) (int64, error)
}

func (svc *service) GenRetailV3(requestData *GenRetailRequest, apiKey string) (*GenRetailV3ModelResponse, error) {
	var dataResponse *GenRetailV3ModelResponse

	url := svc.Cfg.Env.AifcoreHost + "/api/score/genretail/v3"

	requestByte, _ := json.Marshal(requestData)
	request, _ := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(requestByte))
	request.Header.Set(constant.HeaderContentType, constant.HeaderApplicationJSON)
	request.Header.Set(constant.XAPIKey, apiKey)

	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		return nil, err
	}

	responseBodyBytes, _ := io.ReadAll(response.Body)
	defer response.Body.Close()

	json.Unmarshal(responseBodyBytes, &dataResponse)
	dataResponse.StatusCode = response.StatusCode

	return dataResponse, nil
}

func (svc *service) BulkSearchUploadSvc(req []BulkSearchRequest, tempType, apiKey, userId, companyId string) error {
	var bulkSearches []*BulkSearch
	uploadId := uuid.NewString()

	for _, v := range req {
		// check api aif-core to get grade data

		genRetailRequest := &GenRetailRequest{
			LoanNo:   v.LoanNo,
			Name:     v.Name,
			IdCardNo: v.NIK,
			PhoneNo:  v.PhoneNumber,
		}

		genRetailResponse, errRequest := svc.GenRetailV3(genRetailRequest, apiKey)
		if errRequest != nil {
			return errRequest
		}

		bulkSearch := &BulkSearch{
			UploadId:             uploadId,
			TransactionId:        genRetailResponse.Data.TransactionId, // from API
			Name:                 v.Name,
			IdCardNo:             v.NIK,
			PhoneNo:              v.PhoneNumber,
			LoanNo:               v.LoanNo,
			ProbabilityToDefault: genRetailResponse.Data.ProbabilityToDefault, // from API
			Grade:                genRetailResponse.Data.Grade,                // from API
			Date:                 genRetailResponse.Data.Date,                 // from API
			Type:                 tempType,
			UserId:               userId,
			CompanyId:            companyId,
		}

		bulkSearches = append(bulkSearches, bulkSearch)
	}

	err := svc.Repo.StoreImportData(bulkSearches, userId)
	if err != nil {
		return err
	}

	return nil
}

func (svc *service) GetBulkSearchSvc(tierLevel uint, userId, companyId string) ([]BulkSearchResponse, error) {

	bulkSearches, err := svc.Repo.GetAllBulkSearch(tierLevel, userId, companyId)
	if err != nil {
		return nil, err
	}

	var responseBulkSearches []BulkSearchResponse
	for _, v := range bulkSearches {
		bulkSearch := BulkSearchResponse{
			TransactionId:        v.TransactionId,
			Name:                 v.Name,
			IdCardNo:             v.IdCardNo,
			PhoneNo:              v.PhoneNo,
			LoanNo:               v.LoanNo,
			ProbabilityToDefault: v.ProbabilityToDefault,
			Grade:                v.Grade,
			Type:                 v.Type,
			Date:                 v.Date,
		}

		if tierLevel != 2 {
			// make sure only pick from the member uploads
			if userId != v.UserId {
				bulkSearch.PIC = v.User.Name
			}
		}

		responseBulkSearches = append(responseBulkSearches, bulkSearch)
	}

	return responseBulkSearches, nil
}

func (svc *service) GetTotalDataBulk(tierLevel uint, userId, companyId string) (int64, error) {
	count, err := svc.Repo.CountData(tierLevel, userId, companyId)
	return count, err
}
