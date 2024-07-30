package company

func NewService(repo Repository) Service {
	return &service{Repo: repo}
}

type service struct {
	Repo Repository
}

type Service interface {
	FindCompanyByIdSvc(id string) (*Company, error)
	UpdateCompanyByIdSvc(req UpdateCompanyRequest, id string) (Company, error)
}

func (svc *service) FindCompanyByIdSvc(id string) (*Company, error) {
	result, err := svc.Repo.FindOneById(id)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (svc *service) UpdateCompanyByIdSvc(req UpdateCompanyRequest, id string) (Company, error) {
	dataReq := Company{
		CompanyName:    req.CompanyName,
		CompanyAddress: req.CompanyAddress,
		CompanyPhone:   req.CompanyPhone,
		PaymentScheme:  req.PaymentScheme,
		IndustryId:     req.IndustryId,
	}

	company, err := svc.Repo.UpdateOneById(dataReq, id)
	if err != nil {
		return company, err
	}

	return company, nil
}
