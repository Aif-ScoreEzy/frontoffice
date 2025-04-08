package member

import (
	"fmt"
	"front-office/common/constant"
	"front-office/helper"
	"front-office/pkg/core/log/operation"
	"front-office/pkg/core/role"
	"front-office/utility/mailjet"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
)

func NewController(
	service Service,
	roleService role.Service,
	logOperationService operation.Service) Controller {
	return &controller{
		Svc:             service,
		RoleSvc:         roleService,
		LogOperationSvc: logOperationService,
	}
}

type controller struct {
	Svc             Service
	RoleSvc         role.Service
	LogOperationSvc operation.Service
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
	page := c.Query("page", "1")
	limit := c.Query("limit", "10")
	keyword := c.Query("keyword", "")
	roleName := c.Query("role", "")
	status := c.Query("status", "")
	startDate := c.Query("startDate", "")
	endDate := c.Query("endDate", "")

	var roleID string
	if roleName != "" {
		result, err := ctrl.RoleSvc.GetAllRoles(role.RoleFilter{
			Name: roleName,
		})

		if err != nil || result == nil || !result.Success {
			statusCode, resp := helper.GetError(err.Error())
			return c.Status(statusCode).JSON(resp)
		}

		if len(result.Data) == 0 {
			statusCode, resp := helper.GetError(constant.DataNotFound)
			return c.Status(statusCode).JSON(resp)
		}

		roleID = fmt.Sprintf("%v", result.Data[0].RoleId)
	}

	filter := MemberFilter{
		CompanyID: companyId,
		Page:      page,
		Limit:     limit,
		Keyword:   keyword,
		RoleID:    roleID,
		Status:    status,
		StartDate: startDate,
		EndDate:   endDate,
	}

	result, err := ctrl.Svc.GetMemberList(&filter)
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

	addLogRequest := &operation.AddLogRequest{
		MemberId:  updatedMember.Data.MemberId,
		CompanyId: updatedMember.Data.CompanyId,
		Action:    constant.EventUpdateProfile,
	}

	_, err = ctrl.LogOperationSvc.AddLogOperation(addLogRequest)
	if err != nil {
		fmt.Println("Failed to log operation for update profile:", err)
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

	addLogRequest := &operation.AddLogRequest{
		MemberId:  user.Data.MemberId,
		CompanyId: user.Data.CompanyId,
		Action:    constant.EventUpdateProfile,
	}

	_, err = ctrl.LogOperationSvc.AddLogOperation(addLogRequest)
	if err != nil {
		fmt.Println("Failed to log operation for upload profile photo:", err)
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

	member, err := ctrl.Svc.GetMemberBy(&FindUserQuery{
		Id: memberId,
	})
	if err != nil {
		statusCode, resp := helper.GetError(err.Error())
		return c.Status(statusCode).JSON(resp)
	}

	if member.Data.MemberId == 0 {
		statusCode, resp := helper.GetError(constant.DataNotFound)
		return c.Status(statusCode).JSON(resp)
	}

	_, err = ctrl.Svc.UpdateMemberById(memberId, req)
	if err != nil {
		statusCode, resp := helper.GetError(err.Error())
		return c.Status(statusCode).JSON(resp)
	}

	currentTime := time.Now()
	formattedTime := helper.FormatWIB(currentTime)

	if req.Email != nil && member.Data.Email != *req.Email {
		err := mailjet.SendConfirmationEmailUserEmailChangeSuccess(member.Data.Name, member.Data.Email, *req.Email, formattedTime)
		if err != nil {
			statusCode, resp := helper.GetError(err.Error())
			return c.Status(statusCode).JSON(resp)
		}
	}

	addLogRequest := &operation.AddLogRequest{
		MemberId:  member.Data.MemberId,
		CompanyId: member.Data.CompanyId,
		Action:    constant.EventUpdateProfile,
	}

	_, err = ctrl.LogOperationSvc.AddLogOperation(addLogRequest)
	if err != nil {
		fmt.Println("Failed to log operation for update member data by admin:", err)
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
