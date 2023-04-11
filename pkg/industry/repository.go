package industry

import (
	"fmt"
	"front-office/config/database"

	"gorm.io/gorm"
)

func FindOneByID(industry Industry) (Industry, error) {
	err := database.DBConn.First(&industry).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return industry, fmt.Errorf("Industry with ID %s not found", industry.ID)
		}

		return industry, fmt.Errorf("Failed to find industry with ID %s: %v", industry.ID, err)
	}

	return industry, nil
}
