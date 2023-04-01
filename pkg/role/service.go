package role

import (
	"github.com/google/uuid"
)

func CreateRoleSvc(roleRequest RoleRequest) (Role, error) {
	roleID := uuid.NewString()
	dataRole := Role{
		ID:   roleID,
		Name: roleRequest.Name,
	}

	role, err := Create(dataRole)
	if err != nil {
		return role, err
	}

	return role, nil
}

func GetRoleByIDSvc(id string) (Role, error) {
	result, err := FindOneByID(id)
	if err != nil {
		return result, err
	}

	return result, nil
}
