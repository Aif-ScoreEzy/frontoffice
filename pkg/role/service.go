package role

import (
	"github.com/google/uuid"
)

func CreateRoleSvc(req RoleRequest) (RoleResponse, error) {
	roleID := uuid.NewString()
	dataReq := Role{
		ID:          roleID,
		Name:        req.Name,
		Permissions: req.Permissions,
	}

	var dataRole RoleResponse
	role, err := Create(dataReq)
	if err != nil {
		return dataRole, err
	}

	dataRole = RoleResponse{
		ID:          role.ID,
		Name:        role.Name,
		Permissions: role.Permissions,
	}

	return dataRole, nil
}

func GetRoleByIDSvc(id string) (Role, error) {
	result, err := FindOneByID(id)
	if err != nil {
		return result, err
	}

	return result, nil
}

func GetRoleByNameSvc(name string) (Role, error) {
	result, err := FindOneByName(name)
	if err != nil {
		return result, err
	}

	return result, nil
}

func UpdateRoleByIDSvc(req RoleRequest, id string) (RoleResponse, error) {
	dataReq := Role{
		Name:        req.Name,
		Permissions: req.Permissions,
	}

	var dataRole RoleResponse
	role, err := UpdateByID(dataReq, id)
	if err != nil {
		return dataRole, err
	}

	dataRole = RoleResponse{
		ID:          role.ID,
		Name:        role.Name,
		Permissions: role.Permissions,
	}

	return dataRole, nil
}

func DeleteRoleByIDSvc(id string) error {
	err := Delete(id)
	if err != nil {
		return err
	}

	return nil
}
