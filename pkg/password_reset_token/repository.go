package password_reset_token

import "front-office/config/database"

func CreatePasswordResetToken(passwordResetToken *PasswordResetToken) (*PasswordResetToken, error) {
	err := database.DBConn.Debug().Create(&passwordResetToken).Error
	if err != nil {
		return nil, err
	}

	return passwordResetToken, nil
}

func FindOnePasswordResetTokenByToken(token string) (*PasswordResetToken, error) {
	var result *PasswordResetToken
	err := database.DBConn.Debug().First(&result, "token = ?", token).Error
	if err != nil {
		return nil, err
	}

	return result, nil
}

func FindOnePasswordResetTokenByUserID(userID string) (*PasswordResetToken, error) {
	var passwordResetToken *PasswordResetToken

	err := database.DBConn.Debug().First(&passwordResetToken, "user_id = ?", userID).Error
	if err != nil {
		return nil, err
	}

	return passwordResetToken, nil
}
