package role

import (
	"github.com/google/uuid"
)

func NewService(repo Repository) Service {
	return &service{Repo: repo}
}

type service struct {
	Repo Repository
}

type Service interface {
	CreateRoleSvc(req *CreateRoleRequest) (Role, error)
	GetAllRolesSvc() ([]Role, error)
	FindRoleByIDSvc(id string) (*Role, error)
	GetRoleByNameSvc(name string) (*Role, error)
	UpdateRoleByIDSvc(req *UpdateRoleRequest, id string) (*Role, error)
	DeleteRoleByIDSvc(id string) error
}

func (svc *service) CreateRoleSvc(req *CreateRoleRequest) (Role, error) {
	roleID := uuid.NewString()
	dataReq := Role{
		ID:          roleID,
		Name:        req.Name,
		Permissions: req.Permissions,
		TierLevel:   req.TierLevel,
	}

	role, err := svc.Repo.Create(dataReq)
	if err != nil {
		return role, err
	}

	return role, nil
}

func (svc *service) GetAllRolesSvc() ([]Role, error) {
	roles, err := svc.Repo.FindAll()
	if err != nil {
		return roles, err
	}

	return roles, nil
}

func (svc *service) FindRoleByIDSvc(id string) (*Role, error) {
	result, err := svc.Repo.FindOneByID(id)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (svc *service) GetRoleByNameSvc(name string) (*Role, error) {
	result, err := svc.Repo.FindOneByName(name)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (svc *service) UpdateRoleByIDSvc(req *UpdateRoleRequest, id string) (*Role, error) {
	dataReq := &Role{}

	if req.Name != "" {
		dataReq.Name = req.Name
	}

	if req.Permissions != nil {
		dataReq.Permissions = req.Permissions
	}

	role, err := svc.Repo.UpdateByID(dataReq, id)
	if err != nil {
		return role, err
	}

	return role, nil
}

func (svc *service) DeleteRoleByIDSvc(id string) error {
	err := svc.Repo.Delete(id)
	if err != nil {
		return err
	}

	return nil
}
