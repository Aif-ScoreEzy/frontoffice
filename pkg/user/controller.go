package user

import (
	"fmt"
	"front-office/constant"
	"front-office/helper"
	"front-office/pkg/role"
	"os"
	"path/filepath"

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
		resend := "resend"
		req := &UpdateUserRequest{
			Status: &resend,
		}

		_, err = UpdateUserByIDSvc(req, user)
		if err != nil {
			statusCode, resp := helper.GetError(err.Error())
			return c.Status(statusCode).JSON(resp)
		}

		statusCode, resp := helper.GetError(constant.SendEmailFailed)
		return c.Status(statusCode).JSON(resp)
	}

	resp := helper.ResponseSuccess(
		fmt.Sprintf("we've sent an email to %s with a link to activate the account", req.Email),
		nil,
	)

	return c.Status(fiber.StatusCreated).JSON(resp)
}

func UpdateUserByID(c *fiber.Ctx) error {
	req := c.Locals("request").(*UpdateUserRequest)
	companyID := fmt.Sprintf("%v", c.Locals("companyID"))
	userID := c.Params("id")

	user, err := FindUserByIDSvc(userID, companyID)
	if err != nil {
		statusCode, resp := helper.GetError(err.Error())
		return c.Status(statusCode).JSON(resp)
	}

	user, err = UpdateUserByIDSvc(req, user)
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

func UpdateProfile(c *fiber.Ctx) error {
	req := c.Locals("request").(*UpdateProfileRequest)
	userID := fmt.Sprintf("%v", c.Locals("userID"))
	companyID := fmt.Sprintf("%v", c.Locals("companyID"))

	user, err := FindUserByIDSvc(userID, companyID)
	if err != nil {
		statusCode, resp := helper.GetError(err.Error())
		return c.Status(statusCode).JSON(resp)
	}

	user, err = UpdateProfileSvc(req, user)
	if err != nil {
		statusCode, resp := helper.GetError(err.Error())
		return c.Status(statusCode).JSON(resp)
	}

	dataResponse := &UserUpdateResponse{
		ID:        user.ID,
		Name:      user.Name,
		Email:     user.Email,
		Status:    user.Status,
		Active:    user.Active,
		CompanyID: user.CompanyID,
		RoleID:    user.RoleID,
	}

	resp := helper.ResponseSuccess(
		"success to update user",
		dataResponse,
	)

	return c.Status(fiber.StatusOK).JSON(resp)
}

func UploadProfileImage(c *fiber.Ctx) error {
	userID := fmt.Sprintf("%v", c.Locals("userID"))
	companyID := fmt.Sprintf("%v", c.Locals("companyID"))

	user, err := FindUserByIDSvc(userID, companyID)
	if err != nil {
		statusCode, resp := helper.GetError(err.Error())
		return c.Status(statusCode).JSON(resp)
	}

	file, err := c.FormFile("image")
	if err != nil {
		statusCode, resp := helper.GetError(err.Error())
		return c.Status(statusCode).JSON(resp)
	}

	if file.Size > 200*1024 {
		statusCode, resp := helper.GetError(constant.FileSizeIsTooLarge)
		return c.Status(statusCode).JSON(resp)
	}

	ext := filepath.Ext(file.Filename)
	filename := fmt.Sprintf("%s%s", userID, ext)
	filePath := fmt.Sprintf("./public/%s", filename)

	if _, err := os.Stat(filePath); err == nil {
		if err := os.Remove(filePath); err != nil {
			statusCode, resp := helper.GetError(err.Error())
			return c.Status(statusCode).JSON(resp)
		}
	}

	if err := c.SaveFile(file, filePath); err != nil {
		statusCode, resp := helper.GetError(err.Error())
		return c.Status(statusCode).JSON(resp)
	}

	user, err = UploadProfileImageSvc(user, &filename)
	if err != nil {
		statusCode, resp := helper.GetError(err.Error())
		return c.Status(statusCode).JSON(resp)
	}

	dataResponse := &UserUpdateResponse{
		ID:        user.ID,
		Name:      user.Name,
		Email:     user.Email,
		Status:    user.Status,
		Active:    user.Active,
		CompanyID: user.CompanyID,
		RoleID:    user.RoleID,
	}

	resp := helper.ResponseSuccess(
		"success to upload profile image",
		dataResponse,
	)

	return c.Status(fiber.StatusOK).JSON(resp)
}

func GetAllUsers(c *fiber.Ctx) error {
	page := c.Query("page", "1")
	limit := c.Query("limit", "10")
	keyword := c.Query("keyword", "")
	roleName := c.Query("role", "")
	status := c.Query("status", "")
	startDate := c.Query("startDate", "")
	endDate := c.Query("endDate", "")
	companyID := fmt.Sprintf("%v", c.Locals("companyID"))

	var roleID string
	if roleName != "" {
		role, err := role.GetRoleByNameSvc(roleName)
		if err != nil {
			statusCode, resp := helper.GetError(err.Error())
			return c.Status(statusCode).JSON(resp)
		}

		roleID = role.ID
	}

	users, err := GetAllUsersSvc(limit, page, keyword, roleID, status, startDate, endDate, companyID)
	if err != nil {
		statusCode, resp := helper.GetError(err.Error())
		return c.Status(statusCode).JSON(resp)
	}

	totalData, _ := GetTotalDataSvc(keyword, roleID, status, startDate, endDate, companyID)

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
	userID := c.Params("id")
	companyID := fmt.Sprintf("%v", c.Locals("companyID"))

	user, err := FindUserByIDSvc(userID, companyID)
	if err != nil {
		statusCode, resp := helper.GetError(err.Error())
		return c.Status(statusCode).JSON(resp)
	}

	resp := helper.ResponseSuccess(
		"succeed to get a user by ID",
		user,
	)

	return c.Status(fiber.StatusOK).JSON(resp)
}

func DeleteUserByID(c *fiber.Ctx) error {
	userID := c.Params("id")
	companyID := fmt.Sprintf("%v", c.Locals("companyID"))

	_, err := FindUserByIDSvc(userID, companyID)
	if err != nil {
		statusCode, resp := helper.GetError(err.Error())
		return c.Status(statusCode).JSON(resp)
	}

	err = DeleteUserByIDSvc(userID)
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
