package main

import (
	"embed"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"

	"github.com/Planck1858/all-wallet/internal/app"
	"github.com/Planck1858/all-wallet/internal/config"
)

//go:embed migrations/*.sql
var embedMigrations embed.FS
var cfgFlag string

func init() {
	var ok bool
	cfgFlag, ok = os.LookupEnv("ALL_WALLET_CONFIG")
	if !ok {
		cfgFlag = "config.yaml"
	}
}

func main() {
	stopCh := make(chan os.Signal, 1)
	signal.Notify(stopCh, os.Interrupt, syscall.SIGTERM, syscall.SIGINT)

	conf := config.Config{}
	err := config.NewConfig(cfgFlag, &conf)
	if err != nil {
		panic(fmt.Errorf("config init: %v", err))
	}

	zapLog, err := loggerLevel(conf.DebugMode)
	if err != nil {
		panic(fmt.Errorf("logger init: %v", err))
	}
	defer zapLog.Sync()

	log := zapLog.Sugar()
	log.Infof("config print: %s", conf)

	log.Info("application initialization")
	a, err := app.New(log, conf, embedMigrations)
	if err != nil {
		panic(fmt.Errorf("app init: %v", err))
	}

	log.Info("starting application")
	a.Run()

	stop := <-stopCh
	log.Infof("catch stop signal: %v", stop)

	log.Info("finishing application")
	err = a.Stop()
	if err != nil {
		panic(fmt.Errorf("app stop: %v", err))
	}
}

func loggerLevel(debug bool) (*zap.Logger, error) {
	var loggerConfig zap.Config
	if debug {
		loggerConfig = zap.NewDevelopmentConfig()
	}

	loggerConfig = zap.NewProductionConfig()
	loggerConfig.EncoderConfig.TimeKey = "timestamp"
	loggerConfig.EncoderConfig.EncodeTime = zapcore.RFC3339TimeEncoder

	return loggerConfig.Build()
}
