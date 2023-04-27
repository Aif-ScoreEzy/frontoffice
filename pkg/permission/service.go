package permission

import (
	"github.com/google/uuid"
)

func CreatePermissionSvc(permissionReq PermissionRequest) (Permission, error) {
	permissionID := uuid.NewString()
	dataPermission := Permission{
		ID:   permissionID,
		Name: permissionReq.Name,
	}

	result, err := Create(dataPermission)
	if err != nil {
		return result, err
	}

	return result, err
}

func IsPermissionExistSvc(id string) (Permission, error) {
	permission := Permission{
		ID: id,
	}

	result, err := FindOneByID(permission)

	return result, err
}

func GetPermissionByNameSvc(name string) (Permission, error) {
	result, err := FindOneByName(name)
	if err != nil {
		return result, err
	}

	return result, nil
}

func UpdatePermissionByIDSvc(req PermissionRequest, id string) (Permission, error) {
	result, err := UpdateByID(req, id)
	if err != nil {
		return result, err
	}

	return result, nil
}

func DeletePermissionByIDSvc(id string) error {
	err := Delete(id)
	if err != nil {
		return err
	}

	return nil
}
