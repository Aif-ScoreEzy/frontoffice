package industry

func NewService(repo Repository) Service {
	return &service{Repo: repo}
}

type service struct {
	Repo Repository
}

type Service interface {
	IsIndustryIDExistSvc(id string) (Industry, error)
}

func (svc *service) IsIndustryIDExistSvc(id string) (Industry, error) {
	industry := Industry{
		ID: id,
	}

	result, err := svc.Repo.FindOneByID(industry)
	if err != nil {
		return result, err
	}

	return result, nil
}
