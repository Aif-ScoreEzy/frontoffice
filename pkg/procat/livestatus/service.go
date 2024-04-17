package livestatus

func NewService(repo Repository) Service {
	return &service{Repo: repo}
}

type service struct {
	Repo Repository
}

type Service interface {
	CreateJob(data []LiveStatusRequest, totalData int) (uint, error)
	GetJobDetails(jobID uint) ([]*JobDetail, error)
}

func (svc *service) CreateJob(data []LiveStatusRequest, totalData int) (uint, error) {
	dataJob := &Job{
		Total: totalData,
	}

	jobID, err := svc.Repo.CreateJobInTx(dataJob, data)
	if err != nil {
		return 0, err
	}

	return jobID, nil
}

func (svc *service) GetJobDetails(jobID uint) ([]*JobDetail, error) {
	jobDetails, err := svc.Repo.GetJobDetailsByJobID(jobID)
	if err != nil {
		return nil, err
	}

	return jobDetails, nil
}
