package main

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/cordalace/wireguard-for-homies/internal/db/badgerdb"
	"github.com/cordalace/wireguard-for-homies/internal/ip"
	"github.com/cordalace/wireguard-for-homies/internal/manager"
	"github.com/cordalace/wireguard-for-homies/internal/telegram"
	badger "github.com/dgraph-io/badger/v2"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func loggerSync(logger *zap.Logger) {
	_ = logger.Sync()
}

type Config struct {
	LogLevel      zapcore.Level
	TelegramToken string
}

func main() {
	preLogger, err := zap.NewProduction()
	if err != nil {
		panic(err)
	}
	defer loggerSync(preLogger)

	// cfg, err := config.Parse()
	// if err != nil {
	// 	preLogger.Error("error loading config", zap.Error(err))
	// 	return
	// }
	cfg := Config{
		LogLevel:      zapcore.InfoLevel,
		TelegramToken: os.Getenv("TOKEN"),
	}

	logConfig := zap.NewProductionConfig()
	logConfig.Level = zap.NewAtomicLevelAt(cfg.LogLevel)
	logger, err := logConfig.Build()
	if err != nil {
		preLogger.Error("error instanstiating logger", zap.Error(err))
	}
	defer loggerSync(logger)

	i := ip.NewIP()
	if err := i.Init(); err != nil {
		logger.Fatal("error creating ip", zap.Error(err))
	}
	defer i.Close()

	badgerOptions := badger.DefaultOptions("/tmp/badger")
	badgerDB := badgerdb.NewBadgerDB(badgerOptions)

	if err := badgerDB.Init(); err != nil {
		logger.Fatal("error opening badger db", zap.Error(err))
	}
	defer badgerDB.Close()

	wgManager := manager.NewManager(badgerDB.AsManagerDB(), i, logger)
	if err = wgManager.Init(); err != nil {
		logger.Fatal("error initializing wireguard", zap.Error(err))
	}
	defer wgManager.Close()

	tg := telegram.NewTelegram(cfg.TelegramToken, logger)
	if err = tg.Init(); err != nil {
		logger.Fatal("error initializing telegram", zap.Error(err))
	}

	go func() {
		if err := tg.Run(); err != nil {
			logger.Fatal("error running telegram", zap.Error(err))
		}
	}()

	// thanks to Sergey, pretty graceful shutdown here
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	logger.Info("start graceful shutdown")

	tg.Close()

	logger.Info("exiting")
}
