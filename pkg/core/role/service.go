package role

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/google/uuid"
)

func NewService(repo Repository) Service {
	return &service{Repo: repo}
}

type service struct {
	Repo Repository
}

type Service interface {
	FindRoleById(id string) (*AifResponse, error)
	CreateRoleSvc(req *CreateRoleRequest) (Role, error)
	GetAllRolesSvc() ([]Role, error)
	GetRoleByNameSvc(name string) (*Role, error)
	UpdateRoleByIdSvc(req *UpdateRoleRequest, id string) (*Role, error)
	DeleteRoleByIdSvc(id string) error
}

func (s *service) FindRoleById(id string) (*AifResponse, error) {
	res, err := s.Repo.FindOneById(id)
	if err != nil {
		return nil, err
	}

	return parseResponse(res)
}

func (svc *service) CreateRoleSvc(req *CreateRoleRequest) (Role, error) {
	roleId := uuid.NewString()
	dataReq := Role{
		Id:          roleId,
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

func (svc *service) GetRoleByNameSvc(name string) (*Role, error) {
	result, err := svc.Repo.FindOneByName(name)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (svc *service) UpdateRoleByIdSvc(req *UpdateRoleRequest, id string) (*Role, error) {
	dataReq := &Role{}

	if req.Name != "" {
		dataReq.Name = req.Name
	}

	if req.Permissions != nil {
		dataReq.Permissions = req.Permissions
	}

	role, err := svc.Repo.UpdateById(dataReq, id)
	if err != nil {
		return role, err
	}

	return role, nil
}

func (svc *service) DeleteRoleByIdSvc(id string) error {
	err := svc.Repo.Delete(id)
	if err != nil {
		return err
	}

	return nil
}

func parseResponse(response *http.Response) (*AifResponse, error) {
	var baseResponse *AifResponse

	if response == nil {
		return nil, fmt.Errorf("response is nil")
	}
	defer response.Body.Close()

	dataByte, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	if err := json.Unmarshal(dataByte, &baseResponse); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	return baseResponse, nil
}
