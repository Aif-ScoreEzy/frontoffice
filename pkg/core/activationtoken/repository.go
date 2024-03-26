package activationtoken

import (
	"gorm.io/gorm"
)

func NewRepository(db *gorm.DB) Repository {
	return &repository{DB: db}
}

type repository struct {
	DB *gorm.DB
}

type Repository interface {
	FindOneActivationTokenBytoken(token string) (*ActivationToken, error)
	FindOneActivationTokenByUserID(userID string) (*ActivationToken, error)
	CreateActivationToken(activationToken *ActivationToken) (*ActivationToken, error)
}

func (repo *repository) FindOneActivationTokenBytoken(token string) (*ActivationToken, error) {
	var activationToken *ActivationToken

	err := repo.DB.Debug().First(&activationToken, "token = ?", token).Error
	if err != nil {
		return nil, err
	}

	return activationToken, nil
}

func (repo *repository) FindOneActivationTokenByUserID(userID string) (*ActivationToken, error) {
	var activationToken *ActivationToken

	err := repo.DB.Debug().First(&activationToken, "user_id = ?", userID).Error
	if err != nil {
		return nil, err
	}

	return activationToken, nil
}

func (repo *repository) CreateActivationToken(activationToken *ActivationToken) (*ActivationToken, error) {
	err := repo.DB.Debug().Create(&activationToken).Error
	if err != nil {
		return nil, err
	}

	return activationToken, nil
}
