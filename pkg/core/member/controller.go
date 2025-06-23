package member

import (
	"fmt"
	"front-office/common/constant"
	"front-office/helper"
	"front-office/pkg/core/log/operation"
	"front-office/pkg/core/role"
	"front-office/utility/mailjet"
	"log"
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

	member, err := ctrl.Svc.GetMemberBy(&FindUserQuery{
		Email:    email,
		Username: username,
		Key:      key,
	})
	if err != nil {
		return err
	}

	resp := helper.ResponseSuccess(
		"succeed to get a user",
		member,
	)

	return c.Status(fiber.StatusOK).JSON(resp)
}

func (ctrl *controller) GetById(c *fiber.Ctx) error {
	id := c.Params("id")

	member, err := ctrl.Svc.GetMemberBy(&FindUserQuery{
		Id: id,
	})
	if err != nil {
		return err
	}

	resp := helper.ResponseSuccess(
		"succeed to get a user",
		member,
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
		roles, err := ctrl.RoleSvc.GetRoles(role.RoleFilter{
			Name: roleName,
		})
		if err != nil {
			return err
		}

		// if len(roles) == 0 {
		// 	statusCode, resp := helper.GetError(constant.DataNotFound)
		// 	return c.Status(statusCode).JSON(resp)
		// }

		roleID = fmt.Sprintf("%v", roles[0].RoleId)
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

		member, err := ctrl.Svc.GetMemberBy(&FindUserQuery{
			Email: *req.Email,
		})

		if err != nil {
			return err
		}

		oldEmail = member.Email
	}

	err := ctrl.Svc.UpdateProfile(userId, oldEmail, req)
	if err != nil {
		return err
	}

	updatedMember, err := ctrl.Svc.GetMemberBy(&FindUserQuery{
		Id: userId,
	})
	if err != nil {
		return err
	}

	addLogRequest := &operation.AddLogRequest{
		MemberId:  updatedMember.MemberId,
		CompanyId: updatedMember.CompanyId,
		Action:    constant.EventUpdateProfile,
	}

	err = ctrl.LogOperationSvc.AddLogOperation(addLogRequest)
	if err != nil {
		log.Println("Failed to log operation for update profile")
	}

	dataResponse := &UserUpdateResponse{
		Id:        updatedMember.MemberId,
		Name:      updatedMember.Name,
		Email:     updatedMember.Email,
		Active:    updatedMember.Active,
		CompanyId: updatedMember.CompanyId,
		RoleId:    updatedMember.RoleId,
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

	member, err := ctrl.Svc.GetMemberBy(&FindUserQuery{
		Id: userId,
	})
	if err != nil {
		return err
	}

	err = ctrl.Svc.UploadProfileImage(userId, &filename)
	if err != nil {
		return err
	}

	addLogRequest := &operation.AddLogRequest{
		MemberId:  member.MemberId,
		CompanyId: member.CompanyId,
		Action:    constant.EventUpdateProfile,
	}

	err = ctrl.LogOperationSvc.AddLogOperation(addLogRequest)
	if err != nil {
		log.Println("Failed to log operation for upload profile photo")
	}

	dataResponse := &UserUpdateResponse{
		Id:        member.MemberId,
		Name:      member.Name,
		Email:     member.Email,
		Active:    member.Active,
		CompanyId: member.CompanyId,
		RoleId:    member.RoleId,
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
	companyId := fmt.Sprintf("%v", c.Locals("companyId"))

	currentUserId, err := helper.InterfaceToUint(c.Locals("userId"))
	if err != nil {
		statusCode, resp := helper.GetError(err.Error())
		return c.Status(statusCode).JSON(resp)
	}

	member, err := ctrl.Svc.GetMemberBy(&FindUserQuery{
		Id:        memberId,
		CompanyId: companyId,
	})
	if err != nil {
		statusCode, resp := helper.GetError(err.Error())
		return c.Status(statusCode).JSON(resp)
	}

	if member.MemberId == 0 {
		statusCode, resp := helper.GetError(constant.DataNotFound)
		return c.Status(statusCode).JSON(resp)
	}

	err = ctrl.Svc.UpdateMemberById(memberId, req)
	if err != nil {
		return err
	}

	currentTime := time.Now()
	formattedTime := helper.FormatWIB(currentTime)

	if req.Email != nil && member.Email != *req.Email {
		err := mailjet.SendConfirmationEmailUserEmailChangeSuccess(member.Name, member.Email, *req.Email, formattedTime)
		if err != nil {
			statusCode, resp := helper.GetError(err.Error())
			return c.Status(statusCode).JSON(resp)
		}
	}

	var logEvents []string

	if req.Name != nil || req.Email != nil || req.RoleId != nil {
		logEvents = append(logEvents, constant.EventUpdateUserData)
	}

	if req.Active != nil {
		if *req.Active {
			logEvents = append(logEvents, constant.EventActivateUser)
		} else {
			logEvents = append(logEvents, constant.EventInactivateUser)
		}
	}

	for _, event := range logEvents {
		logRequest := &operation.AddLogRequest{
			MemberId:  currentUserId,
			CompanyId: member.CompanyId,
			Action:    event,
		}

		err := ctrl.LogOperationSvc.AddLogOperation(logRequest)
		if err != nil {
			log.Println("Failed to log operation for update member data by admin")
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
	companyId := fmt.Sprintf("%v", c.Locals("companyId"))

	_, err := ctrl.Svc.GetMemberBy(&FindUserQuery{
		Id:        id,
		CompanyId: companyId,
	})
	if err != nil {
		return err
	}

	err = ctrl.Svc.DeleteMemberById(id)
	if err != nil {
		return err
	}

	resp := helper.ResponseSuccess(
		"succeed to delete member",
		nil,
	)

	return c.Status(fiber.StatusOK).JSON(resp)
}
