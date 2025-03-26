package industry

func NewService(repo Repository) Service {
	return &service{Repo: repo}
}

type service struct {
	Repo Repository
}

type Service interface {
	IsIndustryIdExistSvc(id string) (Industry, error)
}

func (svc *service) IsIndustryIdExistSvc(id string) (Industry, error) {
	industry := Industry{
		Id: id,
	}

	result, err := svc.Repo.FindOneById(industry)
	if err != nil {
		return result, err
	}

	return result, nil
}
