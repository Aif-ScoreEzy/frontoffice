package user

import (
	"fmt"
	"front-office/constant"
	"front-office/helper"
	"front-office/pkg/role"

	"github.com/gofiber/fiber/v2"
)

func RegisterMember(c *fiber.Ctx) error {
	req := c.Locals("request").(*RegisterMemberRequest)
	companyID := fmt.Sprintf("%v", c.Locals("companyID"))

	user, _ := FindUserByEmailSvc(req.Email)
	if user != nil {
		statusCode, resp := helper.GetError(constant.DataAlreadyExist)
		return c.Status(statusCode).JSON(resp)
	}

	user, token, err := RegisterMemberSvc(req, companyID)
	if err != nil {
		statusCode, resp := helper.GetError(err.Error())
		return c.Status(statusCode).JSON(resp)
	}

	err = SendEmailActivationSvc(req.Email, token)
	if err != nil {
		fmt.Println("errrrrrrrr", err)
		resend := "resend"
		req := &UpdateUserRequest{
			Status: &resend,
		}

		_, err = UpdateUserByIDSvc(req, user.ID, companyID)
		if err != nil {
			statusCode, resp := helper.GetError(err.Error())
			return c.Status(statusCode).JSON(resp)
		}

		respFailed := helper.ResponseFailed(
			"Send email failed",
		)
		return c.Status(fiber.StatusInternalServerError).JSON(respFailed)
	}

	resp := helper.ResponseSuccess(
		"the activation link has been sent to your email address",
		nil,
	)

	return c.Status(fiber.StatusCreated).JSON(resp)
}

func ActivateUser(c *fiber.Ctx) error {
	key := c.Params("key")

	user, _ := FindUserByKeySvc(key)
	if user == nil {
		statusCode, resp := helper.GetError(constant.DataNotFound)
		return c.Status(statusCode).JSON(resp)
	}

	user, err := ActivateUserByKeySvc(key)
	if err != nil {
		statusCode, resp := helper.GetError(err.Error())
		return c.Status(statusCode).JSON(resp)
	}

	dataResponse := UserUpdateResponse{
		ID:        user.ID,
		Name:      user.Name,
		Email:     user.Email,
		Active:    user.Active,
		CompanyID: user.CompanyID,
		RoleID:    user.RoleID,
	}

	resp := helper.ResponseSuccess(
		"success in activating the user",
		dataResponse,
	)

	return c.Status(fiber.StatusOK).JSON(resp)
}

func DeactiveUser(c *fiber.Ctx) error {
	email := c.Params("email")

	user, _ := FindUserByEmailSvc(email)
	if user == nil {
		statusCode, resp := helper.GetError(constant.DataNotFound)
		return c.Status(statusCode).JSON(resp)
	}

	user, err := DeactivateUserByEmailSvc(email)
	if err != nil {
		statusCode, resp := helper.GetError(err.Error())
		return c.Status(statusCode).JSON(resp)
	}

	dataResponse := UserUpdateResponse{
		ID:        user.ID,
		Name:      user.Name,
		Email:     user.Email,
		Active:    user.Active,
		CompanyID: user.CompanyID,
		RoleID:    user.RoleID,
	}

	resp := helper.ResponseSuccess(
		"success in deactivating the user",
		dataResponse,
	)

	return c.Status(fiber.StatusOK).JSON(resp)
}

func UpdateUserByID(c *fiber.Ctx) error {
	req := c.Locals("request").(*UpdateUserRequest)
	companyID := fmt.Sprintf("%v", c.Locals("companyID"))
	id := c.Params("id")

	_, err := UpdateUserByIDSvc(req, id, companyID)
	if err != nil {
		statusCode, resp := helper.GetError(err.Error())
		return c.Status(statusCode).JSON(resp)
	}

	user, err := FindUserByIDSvc(id, companyID)
	if err != nil {
		statusCode, resp := helper.GetError(err.Error())
		return c.Status(statusCode).JSON(resp)
	}

	dataResponse := UserUpdateResponse{
		ID:        user.ID,
		Name:      user.Name,
		Email:     user.Email,
		Status:    user.Status,
		Active:    user.Active,
		CompanyID: user.CompanyID,
		RoleID:    user.RoleID,
	}

	resp := helper.ResponseSuccess(
		"Success to update user",
		dataResponse,
	)

	return c.Status(fiber.StatusOK).JSON(resp)
}

func GetAllUsers(c *fiber.Ctx) error {
	page := c.Query("page", "1")
	limit := c.Query("limit", "10")
	keyword := c.Query("keyword", "")
	roleName := c.Query("role", "")
	active := c.Query("active", "")
	startDate := c.Query("startDate", "")
	endDate := c.Query("endDate", "")
	userID := fmt.Sprintf("%v", c.Locals("userID"))
	companyID := fmt.Sprintf("%v", c.Locals("companyID"))

	var roleID string
	if roleName != "" {
		role, err := role.GetRoleByNameSvc(roleName)
		if err != nil {
			statusCode, resp := helper.GetError(constant.DataNotFound)
			return c.Status(statusCode).JSON(resp)
		}
		roleID = role.ID
	}

	user, err := FindUserByIDSvc(userID, companyID)
	if err != nil {
		statusCode, resp := helper.GetError(err.Error())
		return c.Status(statusCode).JSON(resp)
	}

	users, err := GetAllUsersSvc(limit, page, keyword, roleID, active, startDate, endDate, user.CompanyID)
	if err != nil {
		statusCode, resp := helper.GetError(err.Error())
		return c.Status(statusCode).JSON(resp)
	}

	totalData, _ := GetTotalDataSvc(keyword, roleID, active, startDate, endDate, user.CompanyID)

	fullResponsePage := map[string]interface{}{
		"total_data": totalData,
		"data":       users,
	}

	resp := helper.ResponseSuccess(
		"Succeed to get all users",
		fullResponsePage,
	)

	return c.Status(fiber.StatusOK).JSON(resp)
}

func GetUserByID(c *fiber.Ctx) error {
	id := c.Params("id")
	companyID := fmt.Sprintf("%v", c.Locals("companyID"))

	user, _ := FindUserByIDSvc(id, companyID)
	if user == nil {
		statusCode, resp := helper.GetError(constant.DataNotFound)
		return c.Status(statusCode).JSON(resp)
	}

	resp := helper.ResponseSuccess(
		"succeed to get a user by ID",
		user,
	)

	return c.Status(fiber.StatusOK).JSON(resp)
}

func DeleteUserByID(c *fiber.Ctx) error {
	id := c.Params("id")
	companyID := fmt.Sprintf("%v", c.Locals("companyID"))

	user, _ := FindUserByIDSvc(id, companyID)
	if user == nil {
		statusCode, resp := helper.GetError(constant.DataNotFound)
		return c.Status(statusCode).JSON(resp)
	}

	err := DeleteUserByIDSvc(id)
	if err != nil {
		statusCode, resp := helper.GetError(err.Error())
		return c.Status(statusCode).JSON(resp)
	}

	resp := helper.ResponseSuccess(
		"user successfully deleted",
		nil,
	)

	return c.Status(fiber.StatusOK).JSON(resp)
}
