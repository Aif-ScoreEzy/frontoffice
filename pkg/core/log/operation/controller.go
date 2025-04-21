package operation

import (
	"fmt"
	"front-office/common/constant"
	"front-office/helper"
	"strings"

	"github.com/gofiber/fiber/v2"
)

func NewController(service Service) Controller {
	return &controller{Svc: service}
}

type controller struct {
	Svc Service
}

type Controller interface {
	GetList(c *fiber.Ctx) error
	GetListByRange(c *fiber.Ctx) error
}

func (ctrl *controller) GetList(c *fiber.Ctx) error {
	companyId := fmt.Sprintf("%v", c.Locals("companyId"))
	page := c.Query("page", "1")
	size := c.Query("size", "10")
	role := c.Query("role")
	event := c.Query("event")
	name := c.Query("name", "")

	// validation for query input
	var eventMap = map[string]string{
		"sign_in":                     constant.EventSignIn,
		"sign_out":                    constant.EventSignOut,
		"change_password":             constant.EventChangePassword,
		"add_new_member":              constant.EventRegisterMember,
		"request_password_reset":      constant.EventRequestPasswordReset,
		"password_reset":              constant.EventPasswordReset,
		"update_profile_account":      constant.EventUpdateProfile,
		"updates_user_data":           constant.EventUpdateUserData,
		"activate_user":               constant.EventActivateUser,
		"inactivate_user":             constant.EventInactivateUser,
		"calculate_score":             constant.EventCalculateScore,
		"download_history_hit":        constant.EventDownloadScoreHistory,
		"change_billing_information":  constant.EventChangeBillingInformation,
		"topup_balance":               constant.EventTopupBalance,
		"submit_payment_confirmation": constant.EventSubmitPaymentConfirmation,
	}

	normalizedEventQuery := strings.ToLower(strings.ReplaceAll(event, " ", "_"))
	mappedEvent, ok := eventMap[normalizedEventQuery]
	if event != "" && !ok {
		statusCode, res := helper.GetError("invalid event type")

		return c.Status(statusCode).JSON(res)
	}

	event = mappedEvent

	filter := &LogOperationFilter{
		CompanyId: companyId,
		Page:      page,
		Size:      size,
		Role:      strings.ToLower(role),
		Event:     event,
		Name:      strings.ToLower(name),
	}

	result, err := ctrl.Svc.GetLogOperations(filter)
	if err != nil {
		statusCode, res := helper.GetError(err.Error())

		return c.Status(statusCode).JSON(res)
	}

	responseBody := helper.ResponseSuccess(
		"succeed to get list of log operation",
		result,
	)

	return c.Status(responseBody.StatusCode).JSON(responseBody)
}

func (ctrl *controller) GetListByRange(c *fiber.Ctx) error {
	companyId := fmt.Sprintf("%v", c.Locals("companyId"))
	page := c.Query("page", "1")
	size := c.Query("size", "10")
	startDate := c.Query("start_date")
	endDate := c.Query(("end_date"))

	filter := &LogRangeFilter{
		Page:      page,
		Size:      size,
		CompanyId: companyId,
		StartDate: startDate,
		EndDate:   endDate,
	}

	result, err := ctrl.Svc.GetByRange(filter)
	if err != nil {
		statusCode, res := helper.GetError(err.Error())

		return c.Status(statusCode).JSON(res)
	}

	responseBody := helper.ResponseSuccess(
		"succeed to get list of log operation",
		result,
	)

	return c.Status(responseBody.StatusCode).JSON(responseBody)
}
