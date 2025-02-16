package app

import (
	"fmt"
	"github.com/Blxssy/AvitoTest/config"
	"github.com/Blxssy/AvitoTest/internal/repo/pg"
	"github.com/Blxssy/AvitoTest/internal/services"
	"github.com/Blxssy/AvitoTest/internal/transport/http"
	"github.com/Blxssy/AvitoTest/pkg/logger"
	"github.com/Blxssy/AvitoTest/pkg/postgres"
	"github.com/Blxssy/AvitoTest/pkg/token"
	"go.uber.org/zap"
	"os"
	"os/signal"
	"syscall"
)

func Run(cfg *config.Config) {
	log := logger.New(cfg.Logger.Level)

	pgRepo, err := postgres.New(cfg.PG)
	if err != nil {
		log.Fatal(fmt.Sprintf("error conntecting to PostgreSQL: %v", err))
	} else {
		log.Info("successfully connected to PostgreSQL")
	}
	version, err := postgres.RunMigrations(pgRepo.DB, cfg.PG)
	if err != nil {
		log.Fatal(fmt.Sprintf("error while running migraions: %v", err))
	}
	log.Info("Migrations version", zap.Uint("v", version))

	t := token.NewTokenGen(token.TokenConfig{
		TokenKey: cfg.Token.TokenKey,
		TokenTTL: cfg.Token.TokenTTL,
	})
	_ = t

	coinRepo := pg.NewCoinRepo(pgRepo)
	coinService := services.NewCoinService(coinRepo, t)

	httpServer := http.NewServer(http.ServerConfig{
		Addr:        cfg.Server.Addr,
		CoinService: coinService,
		Logger:      log,
	})

	go func() {
		if err = httpServer.Run(); err != nil {
			log.Fatal(fmt.Sprintf("error occurred while running HTTP server: %v", err))
		}
	}()
	log.Info(fmt.Sprintf("HTTP server successfully started on %s", cfg.Server.Addr))

	// Graceful Shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT, os.Interrupt)

	<-quit

	// Shutdown HTTP server
	log.Info("shutdown HTTP server...")
	if err = httpServer.Shutdown(); err != nil {
		log.Error(fmt.Sprintf("failed to shutdown HTTP server: %v", err))
	} else {
		log.Info("HTTP server successfully shutdown")
	}

	log.Info("closing PostgreSQL connections...")
	if err = pgRepo.Close(); err != nil {
		log.Error(fmt.Sprintf("failed to closing PostgreSQL connections: %v", err))
	} else {
		log.Info("PostgreSQL connections successfully closed")
	}
}
