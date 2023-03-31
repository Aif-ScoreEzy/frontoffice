package role

import (
	"github.com/google/uuid"
	"github.com/usepzaka/validator"
)

func CreateRoleSvc(roleRequest RoleRequest) error {
	if errValid := validator.ValidateStruct(roleRequest); errValid != nil {
		return errValid
	}

	roleID := uuid.NewString()
	dataRole := Role{
		ID:   roleID,
		Name: roleRequest.Name,
	}

	_, err := Create(dataRole)
	if err != nil {
		return err
	}

	return nil
}

func GetRoleByIDSvc(id string) (Role, error) {
	result, err := FindOneByID(id)
	if err != nil {
		return result, err
	}

	return result, nil
}
