package user

import (
	"errors"
	"fmt"
	"front-office/constant"
	"front-office/helper"
	"front-office/pkg/role"
	"front-office/utility/mailjet"
	"os"
	"strconv"
	"time"

	"github.com/google/uuid"
)

func RegisterMemberSvc(req *RegisterMemberRequest, companyID string) (*User, string, error) {
	userID := uuid.NewString()

	var tierLevel uint
	if req.RoleID != "" {
		result, err := role.FindRoleByIDSvc(req.RoleID)
		if result == nil {
			return nil, "", errors.New(constant.DataNotFound)
		} else if err != nil {
			return nil, "", err
		} else {
			tierLevel = result.TierLevel
		}
	}

	dataUser := &User{
		ID:        userID,
		Name:      req.Name,
		Email:     req.Email,
		Key:       helper.GenerateAPIKey(),
		Image:     "default-profile-image.jpg",
		RoleID:    req.RoleID,
		CompanyID: companyID,
	}

	secret := os.Getenv("JWT_SECRET_KEY")
	minutesToExpired, _ := strconv.Atoi(os.Getenv("JWT_ACTIVATION_EXPIRES_MINUTES"))

	token, err := helper.GenerateToken(secret, minutesToExpired, userID, dataUser.CompanyID, tierLevel)
	if err != nil {
		return nil, "", err
	}

	tokenID := uuid.NewString()
	dataToken := &ActivationToken{
		ID:     tokenID,
		Token:  token,
		UserID: userID,
	}

	user, err := CreateMember(dataUser, dataToken)
	if err != nil {
		return nil, "", err
	}

	return user, token, nil
}

func FindUserByEmailSvc(email string) (*User, error) {
	user, err := FindOneByEmail(email)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func FindUserByKeySvc(key string) (*User, error) {
	user, err := FindOneByKey(key)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func FindUserByIDSvc(id, companyID string) (*User, error) {
	user, err := FindOneByID(id, companyID)
	if err != nil {
		return nil, err
	}

	return user, err
}

func CreateActivationTokenSvc(user *User) (string, *ActivationToken, error) {
	secret := os.Getenv("JWT_SECRET_KEY")
	minutesToExpired, _ := strconv.Atoi(os.Getenv("JWT_ACTIVATION_EXPIRES_MINUTES"))

	token, err := helper.GenerateToken(secret, minutesToExpired, user.ID, user.CompanyID, user.Role.TierLevel)
	if err != nil {
		return "", nil, err
	}

	tokenID := uuid.NewString()
	activationToken := &ActivationToken{
		ID:     tokenID,
		Token:  token,
		UserID: user.ID,
	}

	activationToken, err = CreateActivationToken(activationToken)
	if err != nil {
		return "", nil, err
	}

	return token, activationToken, nil
}

func SendEmailActivationSvc(email, token string) error {
	baseURL := os.Getenv("FRONTEND_BASE_URL")

	variables := map[string]interface{}{
		"link": fmt.Sprintf("%s/activation?key=%s", baseURL, token),
	}

	err := mailjet.CreateMailjet(email, 5188578, variables)
	if err != nil {
		return err
	}

	return nil
}

func FindActivationTokenByTokenSvc(token string) (*ActivationToken, error) {
	result, err := FindOneActivationTokenBytoken(token)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func FindActivationTokenByUserIDSvc(userID string) (*ActivationToken, error) {
	result, err := FindOneActivationTokenByUserID(userID)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func UpdateProfileSvc(req *UpdateProfileRequest, user *User) (*User, error) {
	updateUser := map[string]interface{}{}

	if req.Name != nil {
		updateUser["name"] = *req.Name
	}

	if req.Email != nil {
		result, _ := FindOneByID(user.ID, user.CompanyID)
		if result.Role.TierLevel == 2 {
			return nil, errors.New(constant.RequestProhibited)
		}

		result, _ = FindUserByEmailSvc(*req.Email)
		if result != nil {
			return nil, errors.New(constant.EmailAlreadyExists)
		}

		updateUser["email"] = *req.Email
	}

	updateUser["updated_at"] = time.Now()

	user, err := UpdateOneByID(updateUser, user)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func UploadProfileImageSvc(user *User, filename *string) (*User, error) {
	updateUser := map[string]interface{}{}

	if filename != nil {
		updateUser["image"] = *filename
	}

	updateUser["updated_at"] = time.Now()

	user, err := UpdateOneByID(updateUser, user)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func UpdateUserByIDSvc(req *UpdateUserRequest, user *User) (*User, error) {
	updateUser := map[string]interface{}{}

	if req.Name != nil {
		updateUser["name"] = *req.Name
	}

	if req.Email != nil {
		user, err := FindUserByEmailSvc(*req.Email)
		if err != nil {
			return nil, err
		} else if user != nil {
			return nil, errors.New(constant.EmailAlreadyExists)
		}

		updateUser["email"] = *req.Email
	}

	if req.RoleID != nil {
		role, err := role.FindOneByID(*req.RoleID)
		if role == nil {
			return nil, errors.New(constant.DataNotFound)
		} else if err != nil {
			return nil, err
		}

		updateUser["role_id"] = *req.RoleID
	}

	if req.Active != nil {
		if *req.Active {
			updateUser["status"] = "active"
			updateUser["active"] = true
		} else {
			updateUser["status"] = "inactive"
			updateUser["active"] = false
		}
	}

	if req.Status != nil {
		updateUser["status"] = *req.Status
	}

	updateUser["updated_at"] = time.Now()

	user, err := UpdateOneByID(updateUser, user)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func GetAllUsersSvc(limit, page, keyword, roleID, status, startDate, endDate, companyID string) ([]GetUsersResponse, error) {
	intPage, _ := strconv.Atoi(page)
	intLimit, _ := strconv.Atoi(limit)
	offset := (intPage - 1) * intLimit

	if status != "" && (status != "active" && status != "inactive" && status != "resend" && status != "pending") {
		return nil, errors.New(constant.InvalidStatusValue)
	}

	var startTime, endTime string
	layoutPostgreSQLDate := "2006-01-02"
	if startDate != "" {
		err := helper.ParseDate(layoutPostgreSQLDate, startDate)
		if err != nil {
			return nil, errors.New(constant.InvalidDateFormat)
		}

		startTime = helper.FormatStartTimeForSQL(startDate)

		if endDate == "" {
			endTime = helper.FormatEndTimeForSQL(startDate)
		}
	}

	if endDate != "" {
		err := helper.ParseDate(layoutPostgreSQLDate, endDate)
		if err != nil {
			return nil, errors.New(constant.InvalidDateFormat)
		}

		endTime = helper.FormatEndTimeForSQL(endDate)
	}

	users, err := FindAll(intLimit, offset, keyword, roleID, status, startTime, endTime, companyID)
	if err != nil {
		return nil, err
	}

	var responseUsers []GetUsersResponse
	for _, user := range users {
		responseUser := GetUsersResponse{
			ID:         user.ID,
			Name:       user.Name,
			Email:      user.Email,
			Status:     user.Status,
			Active:     user.Active,
			IsVerified: user.IsVerified,
			CompanyID:  user.CompanyID,
			Role:       user.Role,
			CreatedAt:  user.CreatedAt,
		}
		responseUsers = append(responseUsers, responseUser)
	}

	return responseUsers, nil
}

func GetTotalDataSvc(keyword, roleID, active, startDate, endDate, companyID string) (int64, error) {
	var startTime, endTime string
	layoutPostgreSQLDate := "2006-01-02"
	if startDate != "" {
		err := helper.ParseDate(layoutPostgreSQLDate, startDate)
		if err != nil {
			return 0, errors.New(constant.InvalidDateFormat)
		}

		startTime = helper.FormatStartTimeForSQL(startDate)

		if endDate == "" {
			endTime = helper.FormatEndTimeForSQL(startDate)
		}
	}

	if endDate != "" {
		err := helper.ParseDate(layoutPostgreSQLDate, endDate)
		if err != nil {
			return 0, errors.New(constant.InvalidDateFormat)
		}

		endTime = helper.FormatEndTimeForSQL(endDate)
	}

	count, err := GetTotalData(keyword, roleID, active, startTime, endTime, companyID)
	return count, err
}

func DeleteUserByIDSvc(id string) error {
	err := DeleteByID(id)
	if err != nil {
		return err
	}

	return nil
}
