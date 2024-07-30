package user

import (
	"fmt"
	"front-office/helper"
	"front-office/pkg/core/role"

	"github.com/gofiber/fiber/v2"
)

func NewController(service Service, svcRole role.Service) Controller {
	return &controller{Svc: service, SvcRole: svcRole}
}

type controller struct {
	Svc     Service
	SvcRole role.Service
}

type Controller interface {
	GetAllUsers(c *fiber.Ctx) error
	GetUserById(c *fiber.Ctx) error
	UpdateUserById(c *fiber.Ctx) error
	UpdateProfile(c *fiber.Ctx) error
	UploadProfileImage(c *fiber.Ctx) error
	DeleteUserById(c *fiber.Ctx) error
}

func (ctrl *controller) GetAllUsers(c *fiber.Ctx) error {
	page := c.Query("page", "1")
	limit := c.Query("limit", "10")
	keyword := c.Query("keyword", "")
	roleName := c.Query("role", "")
	status := c.Query("status", "")
	startDate := c.Query("startDate", "")
	endDate := c.Query("endDate", "")
	companyId := fmt.Sprintf("%v", c.Locals("companyId"))

	var roleId string
	if roleName != "" {
		role, err := ctrl.SvcRole.GetRoleByNameSvc(roleName)
		if err != nil {
			statusCode, resp := helper.GetError(err.Error())
			return c.Status(statusCode).JSON(resp)
		}

		roleId = role.Id
	}

	users, err := ctrl.Svc.GetAllUsersSvc(limit, page, keyword, roleId, status, startDate, endDate, companyId)
	if err != nil {
		statusCode, resp := helper.GetError(err.Error())
		return c.Status(statusCode).JSON(resp)
	}

	totalData, _ := ctrl.Svc.GetTotalDataSvc(keyword, roleId, status, startDate, endDate, companyId)

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

func (ctrl *controller) GetUserById(c *fiber.Ctx) error {
	userId := c.Params("id")
	companyId := fmt.Sprintf("%v", c.Locals("companyId"))

	user, err := ctrl.Svc.FindUserByIdAndCompanyIdSvc(userId, companyId)
	if err != nil {
		statusCode, resp := helper.GetError(err.Error())
		return c.Status(statusCode).JSON(resp)
	}

	resp := helper.ResponseSuccess(
		"succeed to get a user by Id",
		user,
	)

	return c.Status(fiber.StatusOK).JSON(resp)
}

func (ctrl *controller) UpdateUserById(c *fiber.Ctx) error {
	req := c.Locals("request").(*UpdateUserRequest)
	companyId := fmt.Sprintf("%v", c.Locals("companyId"))
	userId := c.Params("id")

	user, err := ctrl.Svc.FindUserByIdAndCompanyIdSvc(userId, companyId)
	if err != nil {
		statusCode, resp := helper.GetError(err.Error())
		return c.Status(statusCode).JSON(resp)
	}

	user, err = ctrl.Svc.UpdateUserByIdSvc(req, user)
	if err != nil {
		statusCode, resp := helper.GetError(err.Error())
		return c.Status(statusCode).JSON(resp)
	}

	dataResponse := UserUpdateResponse{
		Id:        user.Id,
		Name:      user.Name,
		Email:     user.Email,
		Status:    user.Status,
		Active:    user.Active,
		CompanyId: user.CompanyId,
		RoleId:    user.RoleId,
	}

	resp := helper.ResponseSuccess(
		"Success to update user",
		dataResponse,
	)

	return c.Status(fiber.StatusOK).JSON(resp)
}

func (ctrl *controller) UpdateProfile(c *fiber.Ctx) error {
	req := c.Locals("request").(*UpdateProfileRequest)
	userId := fmt.Sprintf("%v", c.Locals("userId"))
	companyId := fmt.Sprintf("%v", c.Locals("companyId"))

	user, err := ctrl.Svc.FindUserByIdAndCompanyIdSvc(userId, companyId)
	if err != nil {
		statusCode, resp := helper.GetError(err.Error())
		return c.Status(statusCode).JSON(resp)
	}

	user, err = ctrl.Svc.UpdateProfileSvc(req, user)
	if err != nil {
		statusCode, resp := helper.GetError(err.Error())
		return c.Status(statusCode).JSON(resp)
	}

	dataResponse := &UserUpdateResponse{
		Id:        user.Id,
		Name:      user.Name,
		Email:     user.Email,
		Status:    user.Status,
		Active:    user.Active,
		CompanyId: user.CompanyId,
		RoleId:    user.RoleId,
	}

	resp := helper.ResponseSuccess(
		"success to update user",
		dataResponse,
	)

	return c.Status(fiber.StatusOK).JSON(resp)
}

func (ctrl *controller) UploadProfileImage(c *fiber.Ctx) error {
	userId := fmt.Sprintf("%v", c.Locals("userId"))
	companyId := fmt.Sprintf("%v", c.Locals("companyId"))
	filename := fmt.Sprintf("%v", c.Locals("filename"))

	user, err := ctrl.Svc.FindUserByIdAndCompanyIdSvc(userId, companyId)
	if err != nil {
		statusCode, resp := helper.GetError(err.Error())
		return c.Status(statusCode).JSON(resp)
	}

	user, err = ctrl.Svc.UploadProfileImageSvc(user, &filename)
	if err != nil {
		statusCode, resp := helper.GetError(err.Error())
		return c.Status(statusCode).JSON(resp)
	}

	dataResponse := &UserUpdateResponse{
		Id:        user.Id,
		Name:      user.Name,
		Email:     user.Email,
		Status:    user.Status,
		Active:    user.Active,
		CompanyId: user.CompanyId,
		RoleId:    user.RoleId,
	}

	resp := helper.ResponseSuccess(
		"success to upload profile image",
		dataResponse,
	)

	return c.Status(fiber.StatusOK).JSON(resp)
}

func (ctrl *controller) DeleteUserById(c *fiber.Ctx) error {
	userId := c.Params("id")
	companyId := fmt.Sprintf("%v", c.Locals("companyId"))

	_, err := ctrl.Svc.FindUserByIdAndCompanyIdSvc(userId, companyId)
	if err != nil {
		statusCode, resp := helper.GetError(err.Error())
		return c.Status(statusCode).JSON(resp)
	}

	err = ctrl.Svc.DeleteUserByIdSvc(userId)
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
