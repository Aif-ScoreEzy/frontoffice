package transaction

import (
	"fmt"
	"front-office/common/constant"
	"front-office/helper"
	"front-office/pkg/core/member"

	"github.com/gofiber/fiber/v2"
)

func NewController(service Service, memberService member.Service) Controller {
	return &controller{Svc: service, MemberSvc: memberService}
}

type controller struct {
	Svc       Service
	MemberSvc member.Service
}

type Controller interface {
	GetLogTransactions(c *fiber.Ctx) error
	GetLogTransactionsByDate(c *fiber.Ctx) error
	GetLogTransactionsByRangeDate(c *fiber.Ctx) error
	GetLogTransactionsByMonth(c *fiber.Ctx) error
}

func (ctrl *controller) GetLogTransactions(c *fiber.Ctx) error {
	resLogTrans, statusCode, errRequest := ctrl.Svc.GetLogTransactions()
	if errRequest != nil {
		_, resp := helper.GetError(errRequest.Error())
		return c.Status(statusCode).JSON(resp)
	}

	var transactions []DataLogTrans
	for _, data := range resLogTrans.Data {
		resMember, err := ctrl.MemberSvc.GetMemberBy(&member.FindUserQuery{
			Id: fmt.Sprintf("%d", data.MemberID),
		})

		if err != nil {
			statusCode, resp := helper.GetError(err.Error())
			return c.Status(statusCode).JSON(resp)
		}

		if resMember == nil || !resMember.Success || resMember.Data.MemberId == 0 {
			statusCode, resp := helper.GetError(constant.DataNotFound)
			return c.Status(statusCode).JSON(resp)
		}

		transaction := DataLogTrans{
			Name:      resMember.Data.Name,
			Grade:     data.Grade,
			CreatedAt: data.CreatedAt,
		}

		transactions = append(transactions, transaction)
	}

	responseBody := helper.ResponseSuccess(
		"succeed to get list of log transaction",
		transactions,
	)

	return c.Status(statusCode).JSON(responseBody)
}

func (ctrl *controller) GetLogTransactionsByDate(c *fiber.Ctx) error {
	date := c.Query("date")
	companyId := c.Query("company_id")

	resLogTrans, statusCode, errRequest := ctrl.Svc.GetLogTransactionsByDate(companyId, date)
	if errRequest != nil {
		_, resp := helper.GetError(errRequest.Error())
		return c.Status(statusCode).JSON(resp)
	}

	var transactions []DataLogTrans
	for _, data := range resLogTrans.Data {
		resMember, err := ctrl.MemberSvc.GetMemberBy(&member.FindUserQuery{
			Id: fmt.Sprintf("%d", data.MemberID),
		})

		if err != nil {
			statusCode, resp := helper.GetError(err.Error())
			return c.Status(statusCode).JSON(resp)
		}

		if resMember == nil || !resMember.Success || resMember.Data.MemberId == 0 {
			statusCode, resp := helper.GetError(constant.DataNotFound)
			return c.Status(statusCode).JSON(resp)
		}

		transaction := DataLogTrans{
			Name:      resMember.Data.Name,
			Grade:     data.Grade,
			CreatedAt: data.CreatedAt,
		}

		transactions = append(transactions, transaction)
	}

	responseBody := helper.ResponseSuccess(
		"succeed to get list of log transaction",
		transactions,
	)

	return c.Status(statusCode).JSON(responseBody)
}

func (ctrl *controller) GetLogTransactionsByRangeDate(c *fiber.Ctx) error {
	startDate := c.Query("start_date")
	endDate := c.Query("end_date")
	companyId := c.Query("company_id")
	page := c.Query("page", "1")

	resLogTrans, statusCode, errRequest := ctrl.Svc.GetLogTransactionsByRangeDate(startDate, endDate, companyId, page)
	if errRequest != nil {
		_, resp := helper.GetError(errRequest.Error())
		return c.Status(statusCode).JSON(resp)
	}

	var transactions []DataLogTrans
	for _, data := range resLogTrans.Data {
		resMember, err := ctrl.MemberSvc.GetMemberBy(&member.FindUserQuery{
			Id: fmt.Sprintf("%d", data.MemberID),
		})

		if err != nil {
			statusCode, resp := helper.GetError(err.Error())
			return c.Status(statusCode).JSON(resp)
		}

		if resMember == nil || !resMember.Success || resMember.Data.MemberId == 0 {
			statusCode, resp := helper.GetError(constant.DataNotFound)
			return c.Status(statusCode).JSON(resp)
		}

		transaction := DataLogTrans{
			Name:      resMember.Data.Name,
			Grade:     data.Grade,
			CreatedAt: data.CreatedAt,
		}

		transactions = append(transactions, transaction)
	}

	responseBody := helper.ResponseSuccess(
		"succeed to get list of log transaction",
		transactions,
	)

	return c.Status(statusCode).JSON(responseBody)
}

func (ctrl *controller) GetLogTransactionsByMonth(c *fiber.Ctx) error {
	companyId := c.Query("company_id")
	month := c.Query("month")

	resLogTrans, statusCode, errRequest := ctrl.Svc.GetLogTransactionsByMonth(companyId, month)
	if errRequest != nil {
		_, resp := helper.GetError(errRequest.Error())
		return c.Status(statusCode).JSON(resp)
	}

	var transactions []DataLogTrans
	for _, data := range resLogTrans.Data {
		resMember, err := ctrl.MemberSvc.GetMemberBy(&member.FindUserQuery{
			Id: fmt.Sprintf("%d", data.MemberID),
		})

		if err != nil {
			statusCode, resp := helper.GetError(err.Error())
			return c.Status(statusCode).JSON(resp)
		}

		if resMember == nil || !resMember.Success || resMember.Data.MemberId == 0 {
			statusCode, resp := helper.GetError(constant.DataNotFound)
			return c.Status(statusCode).JSON(resp)
		}

		transaction := DataLogTrans{
			Name:      resMember.Data.Name,
			Grade:     data.Grade,
			CreatedAt: data.CreatedAt,
		}

		transactions = append(transactions, transaction)
	}

	responseBody := helper.ResponseSuccess(
		"succeed to get list of log transaction",
		transactions,
	)

	return c.Status(statusCode).JSON(responseBody)
}
