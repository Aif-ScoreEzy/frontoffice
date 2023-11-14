package auth

import (
	"front-office/config/database"
	activation_token "front-office/pkg/activation-token"
	"front-office/pkg/company"
	"front-office/pkg/password_reset_token"
	"front-office/pkg/user"

	"gorm.io/gorm"
)

func CreateAdmin(company *company.Company, user *user.User, activationToken *activation_token.ActivationToken) (*user.User, error) {
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

	database.DBConn.Preload("Company").Preload("Company.Industry").Preload("Role").First(&user)

	return user, errTx
}

func CreateMember(user *user.User, activationToken *activation_token.ActivationToken) (*user.User, error) {
	errTx := database.DBConn.Transaction(func(tx *gorm.DB) error {
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

func ResetPassword(id, token string, req *PasswordResetRequest) error {
	errTX := database.DBConn.Transaction(func(tx *gorm.DB) error {
		err := tx.Debug().Model(&user.User{}).Where("id = ?", id).Update("password", user.SetPassword(req.Password)).Error
		if err != nil {
			return err
		}

		if err := tx.Debug().Model(&password_reset_token.PasswordResetToken{}).Where("token = ?", token).Update("activation", true).Error; err != nil {
			return err
		}

		return nil
	})

	if errTX != nil {
		return errTX
	}

	return nil
}

func VerifyUserTx(req map[string]interface{}, userID, token string) (*user.User, error) {
	var user *user.User

	errTX := database.DBConn.Transaction(func(tx *gorm.DB) error {
		if err := tx.Debug().Model(&activation_token.ActivationToken{}).Where("token = ?", token).Update("activation", true).Error; err != nil {
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
