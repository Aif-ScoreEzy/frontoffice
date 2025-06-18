package transaction

import (
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
	// scoreezy
	GetLogScoreezy(c *fiber.Ctx) error
	GetLogScoreezyByDate(c *fiber.Ctx) error
	GetLogScoreezyByRangeDate(c *fiber.Ctx) error
	GetLogScoreezyByMonth(c *fiber.Ctx) error

	// product catalog
}
