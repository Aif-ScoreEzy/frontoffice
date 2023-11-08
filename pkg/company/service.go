package company

func FindCompanyByIDSvc(id string) (*Company, error) {
	result, err := FindOneByID(id)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func UpdateCompanyByIDSvc(req UpdateCompanyRequest, id string) (Company, error) {
	dataReq := Company{
		CompanyName:    req.CompanyName,
		CompanyAddress: req.CompanyAddress,
		CompanyPhone:   req.CompanyPhone,
		PaymentScheme:  req.PaymentScheme,
		IndustryID:     req.IndustryID,
	}

	company, err := UpdateOneByID(dataReq, id)
	if err != nil {
		return company, err
	}

	return company, nil
}
