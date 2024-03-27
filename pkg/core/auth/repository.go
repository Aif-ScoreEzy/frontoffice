package auth

import (
	"front-office/pkg/core/activationtoken"
	"front-office/pkg/core/company"
	"front-office/pkg/core/passwordresettoken"
	"front-office/pkg/core/user"

	"gorm.io/gorm"
)

func NewRepository(db *gorm.DB) Repository {
	return &repository{DB: db}
}

type repository struct {
	DB *gorm.DB
}

type Repository interface {
	CreateAdmin(company *company.Company, user *user.User, activationToken *activationtoken.ActivationToken) (*user.User, error)
	CreateMember(user *user.User, activationToken *activationtoken.ActivationToken) (*user.User, error)
	ResetPassword(id, token string, req *PasswordResetRequest) error
	VerifyUserTx(req map[string]interface{}, userID, token string) (*user.User, error)
}

func (repo *repository) CreateAdmin(company *company.Company, user *user.User, activationToken *activationtoken.ActivationToken) (*user.User, error) {
	errTx := repo.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(&company).Error; err != nil {
			return err
		}

		user.CompanyID = company.ID
		if err := tx.Create(&user).Error; err != nil {
			return err
		}

		if err := tx.Debug().Create(&activationToken).Error; err != nil {
			return err
		}

		return nil
	})

	if errTx != nil {
		return user, errTx
	}

	repo.DB.Preload("Company").Preload("Company.Industry").Preload("Role").First(&user)

	return user, errTx
}

func (repo *repository) CreateMember(user *user.User, activationToken *activationtoken.ActivationToken) (*user.User, error) {
	errTx := repo.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Debug().Create(&user).Error; err != nil {
			return err
		}

		if err := tx.Debug().Create(&activationToken).Error; err != nil {
			return err
		}

		return nil
	})

	if errTx != nil {
		return nil, errTx
	}

	return user, nil
}

func (repo *repository) ResetPassword(id, token string, req *PasswordResetRequest) error {
	errTX := repo.DB.Transaction(func(tx *gorm.DB) error {
		err := tx.Debug().Model(&user.User{}).Where("id = ?", id).Update("password", user.SetPassword(req.Password)).Error
		if err != nil {
			return err
		}

		if err := tx.Debug().Model(&passwordresettoken.PasswordResetToken{}).Where("token = ?", token).Update("activation", true).Error; err != nil {
			return err
		}

		return nil
	})

	if errTX != nil {
		return errTX
	}

	return nil
}

func (repo *repository) VerifyUserTx(req map[string]interface{}, userID, token string) (*user.User, error) {
	var user *user.User

	errTX := repo.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Debug().Model(&activationtoken.ActivationToken{}).Where("token = ?", token).Update("activation", true).Error; err != nil {
			return err
		}

		if err := tx.Debug().Model(&user).Where("id = ?", userID).Updates(req).Error; err != nil {
			return err
		}

		return nil
	})

	if errTX != nil {
		return nil, errTX
	}

	return user, nil
}
