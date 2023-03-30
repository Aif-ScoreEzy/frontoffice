package role

import (
	"github.com/google/uuid"
	"github.com/usepzaka/validator"
)

func CreateRoleService(roleRequest RoleRequest) error {
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
