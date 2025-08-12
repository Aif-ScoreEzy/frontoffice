package genretail

import (
	"fmt"
	"front-office/common/constant"
	"front-office/common/model"
	"front-office/internal/apperror"
	"front-office/pkg/core/grade"
)

func NewService(repo Repository, gradeRepo grade.Repository) Service {
	return &service{repo, gradeRepo}
}

type service struct {
	repo      Repository
	gradeRepo grade.Repository
}

type Service interface {
	GenRetailV3(memberId, companyId string, payload *genRetailRequest) (*model.ScoreezyAPIResponse[dataGenRetailV3], error)
	// BulkSearchUploadSvc(req []BulkSearchRequest, tempType, apiKey, userId, companyId string) error
	// GetBulkSearchSvc(tierLevel uint, userId, companyId string) ([]BulkSearchResponse, error)
	// GetTotalDataBulk(tierLevel uint, userId, companyId string) (int64, error)
}

func (svc *service) GenRetailV3(memberId, companyId string, payload *genRetailRequest) (*model.ScoreezyAPIResponse[dataGenRetailV3], error) {
	// make sure parameter settings are set
	productSlug := constant.SlugGenRetailV3
	grade, err := svc.gradeRepo.GetGradesAPI(productSlug, fmt.Sprintf("%v", companyId))
	if err != nil {
		return nil, apperror.MapRepoError(err, "failed to get grades")
	}

	if len(grade.Grades) < 1 {
		return nil, apperror.BadRequest(constant.ParamSettingIsNotSet)
	}

	result, err := svc.repo.GenRetailV3API(memberId, payload)
	if err != nil {
		apperror.MapRepoError(err, "failed to process gen retail v3")
	}

	return result, err
}

// func (svc *service) BulkSearchUploadSvc(req []BulkSearchRequest, tempType, apiKey, userId, companyId string) error {
// 	var bulkSearches []*BulkSearch
// 	uploadId := uuid.NewString()

// 	for _, v := range req {
// 		// check api aif-core to get grade data

// 		genRetailRequest := &genRetailRequest{
// 			LoanNo:   v.LoanNo,
// 			Name:     v.Name,
// 			IdCardNo: v.NIK,
// 			PhoneNo:  v.PhoneNumber,
// 		}

// 		genRetailResponse, errRequest := svc.GenRetailV3(genRetailRequest, apiKey)
// 		if errRequest != nil {
// 			return errRequest
// 		}

// 		bulkSearch := &BulkSearch{
// 			UploadId:             uploadId,
// 			TransactionId:        genRetailResponse.Data.TransactionId, // from API
// 			Name:                 v.Name,
// 			IdCardNo:             v.NIK,
// 			PhoneNo:              v.PhoneNumber,
// 			LoanNo:               v.LoanNo,
// 			ProbabilityToDefault: genRetailResponse.Data.ProbabilityToDefault, // from API
// 			Grade:                genRetailResponse.Data.Grade,                // from API
// 			Date:                 genRetailResponse.Data.Date,                 // from API
// 			Type:                 tempType,
// 			UserId:               userId,
// 			CompanyId:            companyId,
// 		}

// 		bulkSearches = append(bulkSearches, bulkSearch)
// 	}

// 	err := svc.Repo.StoreImportData(bulkSearches, userId)
// 	if err != nil {
// 		return err
// 	}

// 	return nil
// }

// func (svc *service) GetBulkSearchSvc(tierLevel uint, userId, companyId string) ([]BulkSearchResponse, error) {

// 	bulkSearches, err := svc.Repo.GetAllBulkSearch(tierLevel, userId, companyId)
// 	if err != nil {
// 		return nil, err
// 	}

// 	var responseBulkSearches []BulkSearchResponse
// 	for _, v := range bulkSearches {
// 		bulkSearch := BulkSearchResponse{
// 			TransactionId:        v.TransactionId,
// 			Name:                 v.Name,
// 			IdCardNo:             v.IdCardNo,
// 			PhoneNo:              v.PhoneNo,
// 			LoanNo:               v.LoanNo,
// 			ProbabilityToDefault: v.ProbabilityToDefault,
// 			Grade:                v.Grade,
// 			Type:                 v.Type,
// 			Date:                 v.Date,
// 		}

// 		if tierLevel != 2 {
// 			// make sure only pick from the member uploads
// 			if userId != v.UserId {
// 				bulkSearch.PIC = v.User.Name
// 			}
// 		}

// 		responseBulkSearches = append(responseBulkSearches, bulkSearch)
// 	}

// 	return responseBulkSearches, nil
// }

// func (svc *service) GetTotalDataBulk(tierLevel uint, userId, companyId string) (int64, error) {
// 	count, err := svc.Repo.CountData(tierLevel, userId, companyId)
// 	return count, err
// }
