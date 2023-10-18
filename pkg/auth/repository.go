package auth

import (
	"front-office/config/database"
	"front-office/pkg/company"
	"front-office/pkg/user"

	"gorm.io/gorm"
)

func CreateAdmin(company *company.Company, user *user.User, activationToken *user.ActivationToken) (*user.User, error) {
	errTx := database.DBConn.Transaction(func(tx *gorm.DB) error {
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

	database.DBConn.Preload("Company").Preload("Role").First(&user)

	return user, errTx
}

func CreatePasswordResetToken(passwordResetToken *PasswordResetToken) (*PasswordResetToken, error) {
	err := database.DBConn.Debug().Create(&passwordResetToken).Error
	if err != nil {
		return nil, err
	}

	return passwordResetToken, nil
}

func FindOnePasswordTokenByToken(token string) (*PasswordResetToken, error) {
	var result *PasswordResetToken
	err := database.DBConn.Debug().First(&result, "token = ?", token).Error
	if err != nil {
		return nil, err
	}

	return result, nil
}

func FindOnePasswordTokenByUserID(userID string) (*PasswordResetToken, error) {
	var passwordResetToken *PasswordResetToken

	err := database.DBConn.Debug().First(&passwordResetToken, "user_id = ?", userID).Error
	if err != nil {
		return nil, err
	}

	return passwordResetToken, nil
}

func ResetPassword(id, token string, req *PasswordResetRequest) error {
	errTX := database.DBConn.Transaction(func(tx *gorm.DB) error {
		err := tx.Debug().Model(&user.User{}).Where("id = ?", id).Update("password", user.SetPassword(req.Password)).Error
		if err != nil {
			return err
		}

		if err := tx.Debug().Model(&PasswordResetToken{}).Where("token = ?", token).Update("activation", true).Error; err != nil {
			return err
		}

		return nil
	})

	if errTX != nil {
		return nil
	}

	return nil
}
