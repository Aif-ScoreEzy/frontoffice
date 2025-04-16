package server

import (
	"front-office/app/config"
	"front-office/pkg/core"
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"gorm.io/gorm"
)

type fiberServer struct {
	App *fiber.App
	Db  *gorm.DB
	Cfg *config.Config
}

func NewServer(cfg *config.Config, db *gorm.DB) Server {
	return &fiberServer{
		App: fiber.New(),
		Db:  db,
		Cfg: cfg,
	}
}

func (s *fiberServer) Start() {
	s.App.Use(recover.New())
	s.App.Static("/", "./public")
	s.App.Use(cors.New(cors.Config{
		AllowHeaders:     "Origin,Content-Type,Accept,Content-Length,Accept-Language,Accept-Encoding,Connection,Access-Control-Allow-Origin,Access-Control-Allow-Headers,Authorization",
		AllowOrigins:     s.Cfg.Env.FrontendBaseUrl,
		AllowCredentials: true,
		AllowMethods:     "GET,POST,HEAD,PUT,DELETE,PATCH,OPTIONS",
		ExposeHeaders:    "Set-Cookie",
	}))

	api := s.App.Group("/api/fo")
	core.SetupInit(api, s.Cfg, s.Db)

	log.Fatal(s.App.Listen(":" + s.Cfg.Env.Port))
}
