package passwordresettoken

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
	CreatePasswordResetToken(passwordResetToken *PasswordResetToken) (*PasswordResetToken, error)
	FindOnePasswordResetTokenByToken(token string) (*PasswordResetToken, error)
	FindOnePasswordResetTokenByUserID(userID string) (*PasswordResetToken, error)
}

func (repo *repository) CreatePasswordResetToken(passwordResetToken *PasswordResetToken) (*PasswordResetToken, error) {
	err := repo.DB.Debug().Create(&passwordResetToken).Error
	if err != nil {
		return nil, err
	}

	return passwordResetToken, nil
}

func (repo *repository) FindOnePasswordResetTokenByToken(token string) (*PasswordResetToken, error) {
	var result *PasswordResetToken
	err := repo.DB.Debug().First(&result, "token = ?", token).Error
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (repo *repository) FindOnePasswordResetTokenByUserID(userID string) (*PasswordResetToken, error) {
	var passwordResetToken *PasswordResetToken

	err := repo.DB.Debug().First(&passwordResetToken, "user_id = ?", userID).Error
	if err != nil {
		return nil, err
	}

	return passwordResetToken, nil
}
