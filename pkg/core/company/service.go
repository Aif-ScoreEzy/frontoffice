package company

func NewService(repo Repository) Service {
	return &service{Repo: repo}
}

type service struct {
	Repo Repository
}

type Service interface {
	FindCompanyByIDSvc(id string) (*Company, error)
	UpdateCompanyByIDSvc(req UpdateCompanyRequest, id string) (Company, error)
}

func (svc *service) FindCompanyByIDSvc(id string) (*Company, error) {
	result, err := svc.Repo.FindOneByID(id)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (svc *service) UpdateCompanyByIDSvc(req UpdateCompanyRequest, id string) (Company, error) {
	dataReq := Company{
		CompanyName:    req.CompanyName,
		CompanyAddress: req.CompanyAddress,
		CompanyPhone:   req.CompanyPhone,
		PaymentScheme:  req.PaymentScheme,
		IndustryID:     req.IndustryID,
	}

	company, err := svc.Repo.UpdateOneByID(dataReq, id)
	if err != nil {
		return company, err
	}

	return company, nil
}
