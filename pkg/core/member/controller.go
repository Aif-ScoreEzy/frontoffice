package member

import (
	"fmt"
	"front-office/common/constant"
	"front-office/helper"
	"front-office/utility/mailjet"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
)

func NewController(service Service) Controller {
	return &controller{
		Svc: service,
	}
}

type controller struct {
	Svc Service
}

type Controller interface {
	GetBy(c *fiber.Ctx) error
	GetById(c *fiber.Ctx) error
	GetList(c *fiber.Ctx) error
	UpdateProfile(c *fiber.Ctx) error
	UploadProfileImage(c *fiber.Ctx) error
	UpdateMemberById(c *fiber.Ctx) error
	DeleteById(c *fiber.Ctx) error
}

func (ctrl *controller) GetBy(c *fiber.Ctx) error {
	email := c.Query("email")
	username := c.Query("username")
	key := c.Query("key")

	result, err := ctrl.Svc.GetMemberBy(&FindUserQuery{
		Email:    email,
		Username: username,
		Key:      key,
	})
	if err != nil || result == nil || !result.Success {
		statusCode, resp := helper.GetError(err.Error())
		return c.Status(statusCode).JSON(resp)
	}

	resp := helper.ResponseSuccess(
		"succeed to get a user",
		result.Data,
	)

	return c.Status(fiber.StatusOK).JSON(resp)
}

func (ctrl *controller) GetById(c *fiber.Ctx) error {
	id := c.Params("id")

	result, err := ctrl.Svc.GetMemberBy(&FindUserQuery{
		Id: id,
	})
	if err != nil {
		statusCode, resp := helper.GetError(err.Error())
		return c.Status(statusCode).JSON(resp)
	}

	if result == nil || !result.Success || result.Data.MemberId == 0 {
		statusCode, resp := helper.GetError(constant.DataNotFound)
		return c.Status(statusCode).JSON(resp)
	}

	resp := helper.ResponseSuccess(
		"succeed to get a user",
		result.Data,
	)

	return c.Status(fiber.StatusOK).JSON(resp)
}

func (ctrl *controller) GetList(c *fiber.Ctx) error {
	companyId := fmt.Sprintf("%v", c.Locals("companyId"))

	result, err := ctrl.Svc.GetMemberList(companyId)
	if err != nil || result == nil || !result.Success {
		statusCode, resp := helper.GetError(err.Error())
		return c.Status(statusCode).JSON(resp)
	}

	fullResponsePage := map[string]interface{}{
		"data":       result.Data,
		"total_data": result.Meta.Total,
	}

	resp := helper.ResponseSuccess(
		"succeed to get member list",
		fullResponsePage,
	)

	return c.Status(fiber.StatusOK).JSON(resp)
}

func (ctrl *controller) UpdateProfile(c *fiber.Ctx) error {
	req := c.Locals("request").(*UpdateProfileRequest)
	userId := fmt.Sprintf("%v", c.Locals("userId"))
	roleId := fmt.Sprintf("%v", c.Locals("roleId"))

	var oldEmail string
	if req.Email != nil {
		roleIdInt, err := strconv.Atoi(roleId)
		if err != nil {
			statusCode, resp := helper.GetError(err.Error())
			return c.Status(statusCode).JSON(resp)
		}

		roleIdMember := 2
		if roleIdInt == roleIdMember {
			statusCode, resp := helper.GetError(constant.RequestProhibited)
			return c.Status(statusCode).JSON(resp)
		}

		result, err := ctrl.Svc.GetMemberBy(&FindUserQuery{
			Email: *req.Email,
		})

		if err != nil {
			statusCode, resp := helper.GetError(err.Error())
			return c.Status(statusCode).JSON(resp)
		}

		if result != nil && result.Data.MemberId != 0 {
			statusCode, resp := helper.GetError(constant.EmailAlreadyExists)
			return c.Status(statusCode).JSON(resp)
		}

		oldEmail = result.Data.Email
	}

	result, err := ctrl.Svc.UpdateProfile(userId, oldEmail, req)
	if err != nil || result == nil || !result.Success {
		statusCode, resp := helper.GetError(err.Error())
		return c.Status(statusCode).JSON(resp)
	}

	updatedMember, err := ctrl.Svc.GetMemberBy(&FindUserQuery{
		Id: userId,
	})
	if err != nil {
		statusCode, resp := helper.GetError(err.Error())
		return c.Status(statusCode).JSON(resp)
	}

	dataResponse := &UserUpdateResponse{
		Id:        updatedMember.Data.MemberId,
		Name:      updatedMember.Data.Name,
		Email:     updatedMember.Data.Email,
		Active:    updatedMember.Data.Active,
		CompanyId: updatedMember.Data.CompanyId,
		RoleId:    updatedMember.Data.RoleId,
	}

	resp := helper.ResponseSuccess(
		"succeed to update profile",
		dataResponse,
	)

	return c.Status(fiber.StatusOK).JSON(resp)
}

func (ctrl *controller) UploadProfileImage(c *fiber.Ctx) error {
	userId := fmt.Sprintf("%v", c.Locals("userId"))
	filename := fmt.Sprintf("%v", c.Locals("filename"))

	_, err := ctrl.Svc.GetMemberBy(&FindUserQuery{
		Id: userId,
	})
	if err != nil {
		statusCode, resp := helper.GetError(err.Error())
		return c.Status(statusCode).JSON(resp)
	}

	user, err := ctrl.Svc.UploadProfileImage(userId, &filename)
	if err != nil {
		statusCode, resp := helper.GetError(err.Error())
		return c.Status(statusCode).JSON(resp)
	}

	dataResponse := &UserUpdateResponse{
		Id:        user.Data.MemberId,
		Name:      user.Data.Name,
		Email:     user.Data.Email,
		Active:    user.Data.Active,
		CompanyId: user.Data.CompanyId,
		RoleId:    user.Data.RoleId,
	}

	resp := helper.ResponseSuccess(
		"success to upload profile image",
		dataResponse,
	)

	return c.Status(fiber.StatusOK).JSON(resp)
}

func (ctrl *controller) UpdateMemberById(c *fiber.Ctx) error {
	req := c.Locals("request").(*UpdateUserRequest)
	memberId := c.Params("id")

	result, err := ctrl.Svc.GetMemberBy(&FindUserQuery{
		Id: memberId,
	})
	if err != nil {
		statusCode, resp := helper.GetError(err.Error())
		return c.Status(statusCode).JSON(resp)
	}

	if result.Data.MemberId == 0 {
		statusCode, resp := helper.GetError(constant.DataNotFound)
		return c.Status(statusCode).JSON(resp)
	}

	_, err = ctrl.Svc.UpdateMemberByIdSvc(memberId, req)
	if err != nil {
		statusCode, resp := helper.GetError(err.Error())
		return c.Status(statusCode).JSON(resp)
	}

	currentTime := time.Now()
	formattedTime := helper.FormatWIB(currentTime)

	if req.Email != nil && result.Data.Email != *req.Email {
		err := mailjet.SendConfirmationEmailUserEmailChangeSuccess(result.Data.Name, result.Data.Email, *req.Email, formattedTime)
		if err != nil {
			statusCode, resp := helper.GetError(err.Error())
			return c.Status(statusCode).JSON(resp)
		}
	}

	resp := helper.ResponseSuccess(
		"Success to update user",
		nil,
	)

	return c.Status(fiber.StatusOK).JSON(resp)
}

func (ctrl *controller) DeleteById(c *fiber.Ctx) error {
	id := c.Params("id")
	companyId, err := strconv.Atoi(fmt.Sprintf("%v", c.Locals("companyId")))
	if err != nil {
		return err
	}

	result, err := ctrl.Svc.GetMemberBy(&FindUserQuery{
		Id: id,
	})

	if err != nil {
		statusCode, resp := helper.GetError(err.Error())
		return c.Status(statusCode).JSON(resp)
	}

	if result == nil || !result.Success || result.Data.MemberId == 0 || result.Data.CompanyId != uint(companyId) {
		statusCode, resp := helper.GetError(constant.DataNotFound)
		return c.Status(statusCode).JSON(resp)
	}

	result, err = ctrl.Svc.DeleteMemberById(id)
	if err != nil || result == nil || !result.Success {
		statusCode, resp := helper.GetError(err.Error())
		return c.Status(statusCode).JSON(resp)
	}

	resp := helper.ResponseSuccess(
		"succeed to delete member",
		nil,
	)

	return c.Status(fiber.StatusOK).JSON(resp)
}
