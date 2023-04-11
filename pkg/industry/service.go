package industry

func IsIndustryIDExistSvc(id string) (Industry, error) {
	industry := Industry{
		ID: id,
	}

	result, err := FindOneByID(industry)
	if err != nil {
		return result, err
	}

	return result, nil
}
