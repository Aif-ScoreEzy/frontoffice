package auth

import (
	"front-office/config/database"
	"front-office/pkg/company"
	"front-office/pkg/user"

	"gorm.io/gorm"
)

func CreateAdmin(company *company.Company, user *user.User, activationToken *user.ActivationToken) (*user.User, error) {
	database.DBConn.Preload("Company").Preload("Role").First(&user)
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

	database.DBConn.Preload("Company").Preload("Company.Industry").Preload("Role").Preload("Role.Permissions").First(&user)

	return user, errTx
}

func CreatePasswordResetToken(data *PasswordResetToken) error {
	err := database.DBConn.Debug().Create(&data).Error
	if err != nil {
		return err
	}

	return nil
}

func FindOneByPasswordResetToken(token string) (*PasswordResetToken, error) {
	var result *PasswordResetToken
	err := database.DBConn.Debug().First(&result, "token = ?", token).Error
	if err != nil {
		return nil, err
	}

	return result, nil
}

func FindOneByActivationToken(token string) (*user.ActivationToken, error) {
	var result *user.ActivationToken
	err := database.DBConn.Debug().First(&result, "token = ?", token).Error
	if err != nil {
		return nil, err
	}

	return result, nil
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

func UpdateOne(id string, req map[string]interface{}) (*user.User, error) {
	var user *user.User

	err := database.DBConn.Debug().Model(&user).Where("id = ?", id).Updates(req).Error
	if err != nil {
		return nil, err
	}

	return user, nil
}
