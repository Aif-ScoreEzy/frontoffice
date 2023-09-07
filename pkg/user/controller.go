package user

import (
	"front-office/helper"
	"front-office/pkg/company"
	"front-office/pkg/role"

	"github.com/gofiber/fiber/v2"
)

func Register(c *fiber.Ctx) error {
	req := c.Locals("request").(*RegisterUserRequest)

	user, _ := FindUserByUsernameSvc(req.Username)
	if user != nil {

		resp := helper.ResponseFailed("Username already exists")

		return c.Status(fiber.StatusBadRequest).JSON(resp)
	}

	user, _ = FindUserByEmailSvc(req.Email)
	if user != nil {
		resp := helper.ResponseFailed("Email already exists")

		return c.Status(fiber.StatusBadRequest).JSON(resp)
	}

	user, err := RegisterUserSvc(*req)
	if err != nil {
		resp := helper.ResponseFailed(err.Error())

		return c.Status(fiber.StatusInternalServerError).JSON(resp)
	}

	dataResponse := UserResponse{
		ID:       user.ID,
		Name:     user.Name,
		Username: user.Username,
		Email:    user.Email,
		Phone:    user.Phone,
		Active:   user.Active,
		Key:      user.Key,
		Company:  user.Company,
		Role:     user.Role,
	}

	resp := helper.ResponseSuccess(
		"Success to register",
		dataResponse,
	)

	return c.Status(fiber.StatusOK).JSON(resp)
}

func Login(c *fiber.Ctx) error {
	req := c.Locals("request").(*UserLoginRequest)

	user, _ := FindUserByEmailSvc(req.Email)
	if user == nil {
		resp := helper.ResponseFailed("Email or password is incorrect")

		return c.Status(fiber.StatusBadRequest).JSON(resp)
	}

	token, err := LoginSvc(req, user)
	if err != nil && err.Error() == "password is incorrect" {
		resp := helper.ResponseFailed("Email or password is incorrect")

		return c.Status(fiber.StatusBadRequest).JSON(resp)
	} else if err != nil {
		resp := helper.ResponseFailed(err.Error())

		return c.Status(fiber.StatusInternalServerError).JSON(resp)
	}

	data := UserLoginResponse{
		ID:        user.ID,
		Name:      user.Name,
		Username:  user.Username,
		Email:     user.Email,
		CompanyID: user.CompanyID,
		RoleID:    user.RoleID,
		Key:       user.Key,
		Token:     token,
	}

	resp := helper.ResponseSuccess(
		"Success to login",
		data,
	)

	return c.Status(fiber.StatusOK).JSON(resp)
}

func ActivateUser(c *fiber.Ctx) error {
	key := c.Params("key")

	user, _ := FindUserByKeySvc(key)
	if user == nil {
		resp := helper.ResponseFailed("User not found")

		return c.Status(fiber.StatusNotFound).JSON(resp)
	}

	user, err := ActivateUserByKeySvc(key)
	if err != nil {
		resp := helper.ResponseFailed(err.Error())

		return c.Status(fiber.StatusInternalServerError).JSON(resp)
	}

	dataResponse := UserUpdateResponse{
		ID:        user.ID,
		Name:      user.Name,
		Username:  user.Username,
		Email:     user.Email,
		Phone:     user.Phone,
		Active:    user.Active,
		Key:       user.Key,
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
		Username:  user.Username,
		Email:     user.Email,
		Phone:     user.Phone,
		Active:    user.Active,
		Key:       user.Key,
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

	user, _ := FindUserByUsernameSvc(req.Username)
	if user != nil {
		resp := helper.ResponseFailed("Username already exists")

		return c.Status(fiber.StatusBadRequest).JSON(resp)
	}

	user, _ = FindUserByEmailSvc(req.Email)
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
		Username:  user.Username,
		Email:     user.Email,
		Phone:     user.Phone,
		Active:    user.Active,
		Key:       user.Key,
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
