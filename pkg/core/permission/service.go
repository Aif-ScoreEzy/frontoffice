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
	UpdatePermissionByIdSvc(req PermissionRequest, id string) (Permission, error)
	DeletePermissionByIdSvc(id string) error
}

func (svc *service) CreatePermissionSvc(permissionReq PermissionRequest) (Permission, error) {
	permissionId := uuid.NewString()
	dataPermission := Permission{
		Id:   permissionId,
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
		Id: id,
	}

	result, err := svc.Repo.FindOneById(permission)

	return result, err
}

func (svc *service) GetPermissionByNameSvc(name string) (Permission, error) {
	result, err := svc.Repo.FindOneByName(name)
	if err != nil {
		return result, err
	}

	return result, nil
}

func (svc *service) UpdatePermissionByIdSvc(req PermissionRequest, id string) (Permission, error) {
	result, err := svc.Repo.UpdateById(req, id)
	if err != nil {
		return result, err
	}

	return result, nil
}

func (svc *service) DeletePermissionByIdSvc(id string) error {
	err := svc.Repo.Delete(id)
	if err != nil {
		return err
	}

	return nil
}
