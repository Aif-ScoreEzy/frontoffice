package role

import (
	"fmt"
	"front-office/app/config"
	"front-office/common/constant"
	"net/http"

	"gorm.io/gorm"
)

func NewRepository(db *gorm.DB, cfg *config.Config) Repository {
	return &repository{
		Db:  db,
		Cfg: cfg,
	}
}

type repository struct {
	Db  *gorm.DB
	Cfg *config.Config
}

type Repository interface {
	Create(role Role) (Role, error)
	FindAll() (*http.Response, error)
	FindOneById(id string) (*http.Response, error)
	FindOneByName(name string) (*Role, error)
	UpdateById(req *Role, id string) (*Role, error)
	Delete(id string) error
}

func (repo *repository) FindOneById(id string) (*http.Response, error) {
	apiUrl := fmt.Sprintf(`%v/api/core/role/%v`, repo.Cfg.Env.AifcoreHost, id)

	request, err := http.NewRequest(http.MethodGet, apiUrl, nil)
	if err != nil {
		return nil, err
	}

	request.Header.Set(constant.HeaderContentType, constant.HeaderApplicationJSON)

	client := &http.Client{}

	return client.Do(request)
}

func (repo *repository) Create(role Role) (Role, error) {
	result := repo.Db.Create(&role)

	repo.Db.Preload("Permissions").First(&role, "id = ?", role.Id)

	return role, result.Error
}

func (repo *repository) FindAll() (*http.Response, error) {
	apiUrl := fmt.Sprintf(`%v/api/core/role`, repo.Cfg.Env.AifcoreHost)

	request, err := http.NewRequest(http.MethodGet, apiUrl, nil)
	if err != nil {
		return nil, err
	}

	request.Header.Set(constant.HeaderContentType, constant.HeaderApplicationJSON)

	client := &http.Client{}

	return client.Do(request)
}

func (repo *repository) FindOneByName(name string) (*Role, error) {
	var role *Role
	result := repo.Db.First(&role, "name = ?", name)
	if result.Error != nil {
		return nil, result.Error
	}

	return role, nil
}

func (repo *repository) UpdateById(req *Role, id string) (*Role, error) {
	var role *Role

	result := repo.Db.Model(&role).
		Where("id = ?", id).Updates(req)
	if result.Error != nil {
		return nil, result.Error
	}

	return role, nil
}

func (repo *repository) Delete(id string) error {
	var role Role
	err := repo.Db.Where("id = ?", id).Delete(&role).Error

	return err
}
