package permission

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
	CreatePermissionSvc(permissionReq PermissionRequest) (Permission, error)
	IsPermissionExistSvc(id string) (Permission, error)
	GetPermissionByNameSvc(name string) (Permission, error)
	UpdatePermissionByIDSvc(req PermissionRequest, id string) (Permission, error)
	DeletePermissionByIDSvc(id string) error
}

func (svc *service) CreatePermissionSvc(permissionReq PermissionRequest) (Permission, error) {
	permissionID := uuid.NewString()
	dataPermission := Permission{
		ID:   permissionID,
		Name: permissionReq.Name,
	}

	result, err := svc.Repo.Create(dataPermission)
	if err != nil {
		return result, err
	}

	return result, err
}

func (svc *service) IsPermissionExistSvc(id string) (Permission, error) {
	permission := Permission{
		ID: id,
	}

	result, err := svc.Repo.FindOneByID(permission)

	return result, err
}

func (svc *service) GetPermissionByNameSvc(name string) (Permission, error) {
	result, err := svc.Repo.FindOneByName(name)
	if err != nil {
		return result, err
	}

	return result, nil
}

func (svc *service) UpdatePermissionByIDSvc(req PermissionRequest, id string) (Permission, error) {
	result, err := svc.Repo.UpdateByID(req, id)
	if err != nil {
		return result, err
	}

	return result, nil
}

func (svc *service) DeletePermissionByIDSvc(id string) error {
	err := svc.Repo.Delete(id)
	if err != nil {
		return err
	}

	return nil
}
