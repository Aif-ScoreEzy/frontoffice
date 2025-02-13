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
	GetRoleById(id string) (*AifResponse, error)
	GetAllRoles() (*AifResponseWithMultipleData, error)
	CreateRoleSvc(req *CreateRoleRequest) (Role, error)
	GetRoleByNameSvc(name string) (*Role, error)
	UpdateRoleByIdSvc(req *UpdateRoleRequest, id string) (*Role, error)
	DeleteRoleByIdSvc(id string) error
}

func (s *service) GetRoleById(id string) (*AifResponse, error) {
	res, err := s.Repo.FindOneById(id)
	if err != nil {
		return nil, err
	}

	return parseSingleResponse(res)
}

func (s *service) GetAllRoles() (*AifResponseWithMultipleData, error) {
	res, err := s.Repo.FindAll()
	if err != nil {
		return nil, err
	}

	return parseMultipleResponse(res)
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

func parseResponse(response *http.Response, result interface{}) error {
	if response == nil {
		return fmt.Errorf("response is nil")
	}
	defer response.Body.Close()

	dataByte, err := io.ReadAll(response.Body)
	if err != nil {
		return fmt.Errorf("failed to read response body: %w", err)
	}

	if err := json.Unmarshal(dataByte, result); err != nil {
		return fmt.Errorf("failed to unmarshal response: %w", err)
	}

	return nil
}

func parseSingleResponse(response *http.Response) (*AifResponse, error) {
	var baseResponse AifResponse
	if err := parseResponse(response, &baseResponse); err != nil {
		return nil, err
	}

	return &baseResponse, nil
}

func parseMultipleResponse(response *http.Response) (*AifResponseWithMultipleData, error) {
	var baseResponse AifResponseWithMultipleData
	if err := parseResponse(response, &baseResponse); err != nil {
		return nil, err
	}

	return &baseResponse, nil
}
