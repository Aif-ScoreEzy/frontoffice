package user

import (
	"fmt"
	"front-office/constant"
	"front-office/helper"
	"front-office/pkg/company"
	"front-office/pkg/role"

	"github.com/gofiber/fiber/v2"
)

func RegisterMember(c *fiber.Ctx) error {
	req := c.Locals("request").(*RegisterMemberRequest)
	userID := fmt.Sprintf("%v", c.Locals("userID"))

	loggedUser, err := FindUserByIDSvc(userID)
	// Only user with role of admin can register another user
	if err != nil {
		statusCode, resp := helper.GetError(err.Error())
		return c.Status(statusCode).JSON(resp)
	} else if loggedUser.Role.Name != "admin" {
		statusCode, resp := helper.GetError(constant.RequestProhibited)
		return c.Status(statusCode).JSON(resp)
	}

	user, err := FindUserByEmailSvc(req.Email)
	if err != nil {
		statusCode, resp := helper.GetError(err.Error())
		return c.Status(statusCode).JSON(resp)
	} else if user != nil {
		statusCode, resp := helper.GetError(constant.DataAlreadyExist)
		return c.Status(statusCode).JSON(resp)
	}

	_, err = RegisterMemberSvc(req, loggedUser)
	if err != nil {
		statusCode, resp := helper.GetError(err.Error())
		return c.Status(statusCode).JSON(resp)
	}

	resp := helper.ResponseSuccess(
		"We sent an email with the credentials for login to "+req.Email,
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
		Phone:     user.Phone,
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
		resp := helper.ResponseFailed("User not found")

		return c.Status(fiber.StatusNotFound).JSON(resp)
	}

	user, err := DeactivateUserByEmailSvc(email)
	if err != nil {
		resp := helper.ResponseFailed(err.Error())

		return c.Status(fiber.StatusInternalServerError).JSON(resp)
	}

	dataResponse := UserUpdateResponse{
		ID:        user.ID,
		Name:      user.Name,
		Email:     user.Email,
		Phone:     user.Phone,
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
	id := c.Params("id")

	_, err := FindUserByIDSvc(id)
	if err != nil {
		resp := helper.ResponseFailed(err.Error())

		return c.Status(fiber.StatusBadRequest).JSON(resp)
	}

	user, _ := FindUserByEmailSvc(req.Email)
	if user != nil {
		resp := helper.ResponseFailed("Email already exists")

		return c.Status(fiber.StatusBadRequest).JSON(resp)
	}

	_, err = company.IsCompanyIDExistSvc(req.CompanyID)
	if err != nil {
		resp := helper.ResponseFailed(err.Error())

		return c.Status(fiber.StatusBadRequest).JSON(resp)
	}

	_, err = role.IsRoleIDExistSvc(req.RoleID)
	if err != nil {
		resp := helper.ResponseFailed(err.Error())

		return c.Status(fiber.StatusBadRequest).JSON(resp)
	}

	user, err = UpdateUserByIDSvc(req, id)
	if err != nil {
		resp := helper.ResponseFailed(err.Error())

		return c.Status(fiber.StatusInternalServerError).JSON(resp)
	}

	dataResponse := UserUpdateResponse{
		ID:        user.ID,
		Name:      user.Name,
		Email:     user.Email,
		Phone:     user.Phone,
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
	users, err := GetAllUsersSvc()
	if err != nil {
		resp := helper.ResponseFailed(err.Error())

		return c.Status(fiber.StatusInternalServerError).JSON(resp)
	}

	resp := helper.ResponseSuccess(
		"Succeed to get all users",
		users,
	)

	return c.Status(fiber.StatusOK).JSON(resp)
}

func GetUserByID(c *fiber.Ctx) error {
	id := c.Params("id")

	user, _ := FindUserByIDSvc(id)
	if user == nil {
		resp := helper.ResponseFailed("User not found")

		return c.Status(fiber.StatusNotFound).JSON(resp)
	}

	resp := helper.ResponseSuccess(
		"Succeed to get a user by ID",
		user,
	)

	return c.Status(fiber.StatusOK).JSON(resp)
}
