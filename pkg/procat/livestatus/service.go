package livestatus

func NewService(repo Repository) Service {
	return &service{Repo: repo}
}

type service struct {
	Repo Repository
}

type Service interface {
	CreateJob(data *FIFRequests, totalData int) error
}

func (svc *service) CreateJob(data *FIFRequests, totalData int) error {
	dataJob := &Job{
		Total: totalData,
	}

	if err := svc.Repo.CreateJobInTx(dataJob, data); err != nil {
		return err
	}

	return nil
}
