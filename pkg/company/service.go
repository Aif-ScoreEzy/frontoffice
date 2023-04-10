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
