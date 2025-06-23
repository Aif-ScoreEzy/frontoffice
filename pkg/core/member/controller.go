package member

import (
	"fmt"
	"front-office/common/constant"
	"front-office/helper"
	"front-office/internal/apperror"
	"front-office/pkg/core/log/operation"
	"front-office/pkg/core/role"
	"front-office/utility/mailjet"
	"log"
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

	filter := &MemberFilter{
		CompanyID: companyId,
		Page:      c.Query("page", "1"),
		Limit:     c.Query("limit", "10"),
		Keyword:   c.Query("keyword", ""),
		RoleName:  c.Query("role", ""),
		Status:    c.Query("status", ""),
		StartDate: c.Query("startDate", ""),
		EndDate:   c.Query("endDate", ""),
	}

	users, meta, err := ctrl.Svc.GetMemberList(filter)
	if err != nil {
		return err
	}

	return c.Status(fiber.StatusOK).JSON(helper.ResponseSuccess(
		"succeed to get member list",
		map[string]interface{}{
			"data":       users,
			"total_data": meta.Total,
		},
	))
}

func (ctrl *controller) UpdateProfile(c *fiber.Ctx) error {
	req := c.Locals("request").(*UpdateProfileRequest)

	userId := fmt.Sprintf("%v", c.Locals("userId"))
	roleId, err := helper.InterfaceToUint(c.Locals("roleId"))
	if err != nil {
		return apperror.Unauthorized("invalid role id session")
	}

	updateResp, err := ctrl.Svc.UpdateProfile(userId, roleId, req)
	if err != nil {
		return err
	}

	return c.Status(fiber.StatusOK).JSON(helper.ResponseSuccess(
		"succeed to update profile",
		updateResp,
	))
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

	dataResponse := &userUpdateResponse{
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
