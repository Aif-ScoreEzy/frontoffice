package company

func IsCompanyIDExistSvc(id string) (Company, error) {
	company := Company{
		ID: id,
	}

	result, err := FindOneByID(company)
	if err != nil {
		return result, err
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
