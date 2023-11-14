package activation_token

import (
	"front-office/config/database"
)

func FindOneActivationTokenBytoken(token string) (*ActivationToken, error) {
	var activationToken *ActivationToken

	err := database.DBConn.Debug().First(&activationToken, "token = ?", token).Error
	if err != nil {
		return nil, err
	}

	return activationToken, nil
}

func FindOneActivationTokenByUserID(userID string) (*ActivationToken, error) {
	var activationToken *ActivationToken

	err := database.DBConn.Debug().First(&activationToken, "user_id = ?", userID).Error
	if err != nil {
		return nil, err
	}

	return activationToken, nil
}

func CreateActivationToken(activationToken *ActivationToken) (*ActivationToken, error) {
	err := database.DBConn.Debug().Create(&activationToken).Error
	if err != nil {
		return nil, err
	}

	return activationToken, nil
}
