package transaction

import (
	"fmt"
	"front-office/common/constant"
	"front-office/helper"
	"front-office/pkg/core/member"

	"github.com/gofiber/fiber/v2"
)

func (ctrl *controller) GetLogScoreezy(c *fiber.Ctx) error {
	resLogTrans, statusCode, errRequest := ctrl.Svc.GetLogScoreezy()
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

func (ctrl *controller) GetLogScoreezyByDate(c *fiber.Ctx) error {
	date := c.Query("date")
	companyId := c.Query("company_id")

	resLogTrans, statusCode, errRequest := ctrl.Svc.GetLogScoreezyByDate(companyId, date)
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

func (ctrl *controller) GetLogScoreezyByRangeDate(c *fiber.Ctx) error {
	startDate := c.Query("start_date")
	endDate := c.Query("end_date")
	companyId := c.Query("company_id")
	page := c.Query("page", "1")

	resLogTrans, statusCode, errRequest := ctrl.Svc.GetLogScoreezyByRangeDate(startDate, endDate, companyId, page)
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

func (ctrl *controller) GetLogScoreezyByMonth(c *fiber.Ctx) error {
	companyId := c.Query("company_id")
	month := c.Query("month")

	resLogTrans, statusCode, errRequest := ctrl.Svc.GetLogScoreezyByMonth(companyId, month)
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
