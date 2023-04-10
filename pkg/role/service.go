package role

import (
	"github.com/google/uuid"
)

func CreateRoleSvc(req RoleRequest) (Role, error) {
	roleID := uuid.NewString()
	dataReq := Role{
		ID:          roleID,
		Name:        req.Name,
		Permissions: req.Permissions,
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

func IsRoleIDExistSvc(id string) (Role, error) {
	role := Role{
		ID: id,
	}

	result, err := FindOneByID(role)
	if err != nil {
		return result, err
	}

	return role, nil
}

func GetRoleByNameSvc(name string) (Role, error) {
	result, err := FindOneByName(name)
	if err != nil {
		return result, err
	}

	return result, nil
}

func UpdateRoleByIDSvc(req RoleRequest, id string) (Role, error) {
	dataReq := Role{
		Name:        req.Name,
		Permissions: req.Permissions,
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
