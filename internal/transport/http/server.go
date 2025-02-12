package http

import (
	"fmt"
	"github.com/Blxssy/AvitoTest/internal/services"
	v1 "github.com/Blxssy/AvitoTest/internal/transport/http/v1"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/requestid"
	"go.uber.org/zap"
)

type Server struct {
	addr string

	coinService services.CoinService

	logger *zap.Logger
	app    *fiber.App
}

type ServerConfig struct {
	Addr string

	CoinService services.CoinService

	Logger *zap.Logger
}

func NewServer(cfg ServerConfig) *Server {
	server := &Server{
		addr:        cfg.Addr,
		logger:      cfg.Logger,
		coinService: cfg.CoinService,
		app:         nil,
	}

	server.app = fiber.New(fiber.Config{})

	server.init()

	return server
}

func (s *Server) Run() error {
	if err := s.app.Listen(s.addr); err != nil {
		return fmt.Errorf("listening HTTP server: %w", err)
	}
	return nil
}

func (s *Server) Shutdown() error {
	if err := s.app.Shutdown(); err != nil {
		return fmt.Errorf("shutdown HTTP server: %w", err)
	}
	return nil
}

func (s *Server) init() {
	s.app.Use(cors.New())
	s.app.Use(requestid.New())
	s.app.Use(func(ctx *fiber.Ctx) error {
		ctx.Response().Header.Set(fiber.HeaderContentType, fiber.MIMEApplicationJSON)
		return ctx.Next()
	})

	s.setHandlers()
}

func (s *Server) setHandlers() {
	handlerV1 := v1.NewHandler(v1.HandlerConfig{
		CoinService: s.coinService,
		Logger:      s.logger,
	})
	{
		handlerV1.Init(s.app)
	}
}
