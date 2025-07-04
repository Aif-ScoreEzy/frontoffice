package role

import (
	"front-office/internal/apperror"
)

func NewService(repo Repository) Service {
	return &service{Repo: repo}
}

type service struct {
	Repo Repository
}

type Service interface {
	GetRoles(filter RoleFilter) ([]*MstRole, error)
	GetRoleById(id string) (*MstRole, error)
}

func (s *service) GetRoles(filter RoleFilter) ([]*MstRole, error) {
	roles, err := s.Repo.CallGetRolesAPI(filter)
	if err != nil {
		return nil, apperror.MapRepoError(err, "failed to fetch roles")
	}

	if len(roles) == 0 {
		return nil, apperror.NotFound("role not found")

	}

	return roles, nil
}

func (s *service) GetRoleById(id string) (*MstRole, error) {
	role, err := s.Repo.CallGetRoleAPI(id)
	if err != nil {
		return nil, apperror.MapRepoError(err, "failed to fetch role")
	}

	if role.RoleId == 0 {
		return nil, apperror.NotFound("role not found")
	}

	return role, nil
}
