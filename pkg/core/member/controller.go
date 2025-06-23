package member

import (
	"fmt"
	"front-office/helper"
	"front-office/internal/apperror"
	"front-office/pkg/core/log/operation"
	"front-office/pkg/core/role"

	"github.com/gofiber/fiber/v2"
)

func NewController(
	service Service,
	roleService role.Service,
	logOperationService operation.Service) Controller {
	return &controller{
		svc:             service,
		roleSvc:         roleService,
		logOperationSvc: logOperationService,
	}
}

type controller struct {
	svc             Service
	roleSvc         role.Service
	logOperationSvc operation.Service
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
	member, err := ctrl.svc.GetMemberBy(&FindUserQuery{
		Email:    c.Query("email"),
		Username: c.Query("username"),
		Key:      c.Query("key"),
	})
	if err != nil {
		return err
	}

	return c.Status(fiber.StatusOK).JSON(helper.ResponseSuccess(
		"succeed to get a user",
		member,
	))
}

func (ctrl *controller) GetById(c *fiber.Ctx) error {
	id := c.Params("id")

	member, err := ctrl.svc.GetMemberBy(&FindUserQuery{
		Id: id,
	})
	if err != nil {
		return err
	}

	return c.Status(fiber.StatusOK).JSON(helper.ResponseSuccess(
		"succeed to get a user",
		member,
	))
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

	users, meta, err := ctrl.svc.GetMemberList(filter)
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

	updateResp, err := ctrl.svc.UpdateProfile(userId, roleId, req)
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

	resp, err := ctrl.svc.UploadProfileImage(userId, &filename)
	if err != nil {
		return err
	}

	return c.Status(fiber.StatusOK).JSON(helper.ResponseSuccess(
		"success to upload profile image",
		resp,
	))
}

func (ctrl *controller) UpdateMemberById(c *fiber.Ctx) error {
	req := c.Locals("request").(*UpdateUserRequest)

	memberId := c.Params("id")
	companyId := fmt.Sprintf("%v", c.Locals("companyId"))

	currentUserId, err := helper.InterfaceToUint(c.Locals("userId"))
	if err != nil {
		return apperror.Unauthorized("invalid user id session")
	}

	roleId, err := helper.InterfaceToUint(c.Locals("roleId"))
	if err != nil {
		return apperror.Unauthorized("invalid role id session")
	}

	err = ctrl.svc.UpdateMemberById(currentUserId, roleId, companyId, memberId, req)
	if err != nil {
		return err
	}

	return c.Status(fiber.StatusOK).JSON(helper.ResponseSuccess(
		"success to update user",
		nil,
	))
}

func (ctrl *controller) DeleteById(c *fiber.Ctx) error {
	id := c.Params("id")
	companyId := fmt.Sprintf("%v", c.Locals("companyId"))

	_, err := ctrl.svc.GetMemberBy(&FindUserQuery{
		Id:        id,
		CompanyId: companyId,
	})
	if err != nil {
		return err
	}

	err = ctrl.svc.DeleteMemberById(id)
	if err != nil {
		return err
	}

	resp := helper.ResponseSuccess(
		"succeed to delete member",
		nil,
	)

	return c.Status(fiber.StatusOK).JSON(resp)
}
