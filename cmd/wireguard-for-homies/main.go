package main

import (
	"github.com/cordalace/wireguard-for-homies/internal/db/badgerdb"
	"github.com/cordalace/wireguard-for-homies/internal/manager"
	badger "github.com/dgraph-io/badger/v2"
	"github.com/vishvananda/netlink"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func loggerSync(logger *zap.Logger) {
	_ = logger.Sync()
}

type Config struct {
	LogLevel zapcore.Level
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
		LogLevel: zapcore.InfoLevel,
	}

	logConfig := zap.NewProductionConfig()
	logConfig.Level = zap.NewAtomicLevelAt(cfg.LogLevel)
	logger, err := logConfig.Build()
	if err != nil {
		preLogger.Error("error instanstiating logger", zap.Error(err))
	}
	defer loggerSync(logger)

	netlinkHandle, err := netlink.NewHandle()
	if err != nil {
		logger.Fatal("error creating netlink handle", zap.Error(err))
	}
	defer netlinkHandle.Delete()

	badgerOptions := badger.DefaultOptions("/tmp/badger")
	badgerDB := badgerdb.NewBadgerDB(badgerOptions)
	err = badgerDB.Init()
	if err != nil {
		logger.Fatal("error opening badger db", zap.Error(err))
	}
	defer badgerDB.Close()

	wgManager := manager.NewManager(badgerDB.AsManagerDB(), netlinkHandle, logger)
	err = wgManager.Init()
	if err != nil {
		logger.Fatal("error initializing wireguard", zap.Error(err))
	}
	defer wgManager.Close()
}
