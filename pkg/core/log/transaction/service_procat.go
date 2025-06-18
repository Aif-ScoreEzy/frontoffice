package transaction

import (
	"front-office/common/model"
	"front-office/helper"
)

func (svc *service) GetLogTransSuccessCount(jobId string) (*model.AifcoreAPIResponse[*getSuccessCountDataResponse], error) {
	response, err := svc.repo.CallLogTransSuccessCountAPI(jobId)
	if err != nil {
		return nil, err
	}

	result, err := helper.ParseAifcoreAPIResponse[*getSuccessCountDataResponse](response)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (svc *service) UpdateLogProCat(transId string, req *UpdateTransRequest) (*model.AifcoreAPIResponse[any], error) {
	data := map[string]interface{}{}

	if req.Success != nil {
		data["success"] = *req.Success
	}

	response, err := svc.repo.CallUpdateLogTransAPI(transId, data)
	if err != nil {
		return nil, err
	}

	result, err := helper.ParseAifcoreAPIResponse[any](response)
	if err != nil {
		return nil, err
	}

	return result, nil
}
