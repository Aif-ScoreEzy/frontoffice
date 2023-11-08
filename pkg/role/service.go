package role

import (
	"github.com/google/uuid"
)

func CreateRoleSvc(req *CreateRoleRequest) (Role, error) {
	roleID := uuid.NewString()
	dataReq := Role{
		ID:          roleID,
		Name:        req.Name,
		Permissions: req.Permissions,
		TierLevel:   req.TierLevel,
	}

	role, err := Create(dataReq)
	if err != nil {
		return role, err
	}

	return role, nil
}

func GetAllRolesSvc() ([]Role, error) {
	roles, err := FindAll()
	if err != nil {
		return roles, err
	}

	return roles, nil
}

func FindRoleByIDSvc(id string) (*Role, error) {
	result, err := FindOneByID(id)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func GetRoleByNameSvc(name string) (*Role, error) {
	result, err := FindOneByName(name)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func UpdateRoleByIDSvc(req *UpdateRoleRequest, id string) (*Role, error) {
	dataReq := &Role{}

	if req.Name != "" {
		dataReq.Name = req.Name
	}

	if req.Permissions != nil {
		dataReq.Permissions = req.Permissions
	}

	role, err := UpdateByID(dataReq, id)
	if err != nil {
		return role, err
	}

	return role, nil
}

func DeleteRoleByIDSvc(id string) error {
	err := Delete(id)
	if err != nil {
		return err
	}

	return nil
}
