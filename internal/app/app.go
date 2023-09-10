package app

import (
	"embed"
	"fmt"
	"time"

	pgxw "github.com/Planck1858/pgxwrapper"
	"github.com/avast/retry-go"
	_ "github.com/lib/pq"
	"github.com/pkg/errors"
	"github.com/pressly/goose/v3"
	"go.uber.org/zap"
	tele "gopkg.in/telebot.v3"

	"github.com/Planck1858/all-wallet/internal/config"
	"github.com/Planck1858/all-wallet/internal/currency_api"
	"github.com/Planck1858/all-wallet/internal/store"
	"github.com/Planck1858/all-wallet/internal/telegram"
)

type App struct {
	log      *zap.SugaredLogger
	db       pgxw.PgDatabase
	telegram *telegram.Telegram
}

func New(log *zap.SugaredLogger, conf config.Config, embedMigrations embed.FS) (*App, error) {
	dsn := fmt.Sprintf("host=%s user=%s dbname=%s password=%s sslmode=%s",
		conf.Database.Address, conf.Database.User, conf.Database.Database, conf.Database.Password, conf.Database.SSLMode)

	// DB
	var db *pgxw.DB
	if err := retry.Do(func() error {
		var dbErr error

		db, dbErr = pgxw.Open(
			pgxw.OptionDSN(dsn),
			pgxw.OptionTicker(time.Second*3),
			pgxw.OptionAttempts(4),
			pgxw.OptionEnableLogs(true))
		if dbErr != nil {
			return dbErr
		}

		return nil
	}, retry.Attempts(5), retry.Delay(time.Second*3)); err != nil {
		return nil, errors.Wrap(err, "db connection")
	}

	// Bot client
	botClient, err := tele.NewBot(tele.Settings{
		Token:   conf.BotToken,
		Poller:  &tele.LongPoller{Timeout: 10 * time.Second},
		Verbose: conf.DebugMode,
	})
	if err != nil {
		return nil, errors.Wrap(err, "telegram bot initialization")
	}

	// Migrations
	if conf.Database.InitMigrations {
		err = initMigrations(db, embedMigrations)
		if err != nil {
			return nil, errors.Wrap(err, "migrations initialization")
		}
	}

	// Currency client
	currencyClient := currency_api.New(log)

	// Telegram bot
	tgBot := telegram.New(log, store.New(log, db, currencyClient), botClient, currencyClient)

	return &App{
		log:      log,
		db:       db,
		telegram: tgBot,
	}, nil
}

func (a *App) Run() {
	a.telegram.RegisterAllHandlers()
	a.telegram.Run()
}

func (a *App) Stop() error {
	a.telegram.Stop()

	err := a.db.GetDB().Close()
	if err != nil {
		return err
	}

	return nil
}

func initMigrations(db pgxw.PgDatabase, embedMigrations embed.FS) error {
	goose.SetBaseFS(embedMigrations)

	if err := goose.Up(db.GetDB(), "migrations"); err != nil {
		return errors.Wrap(err, "migration goose up")
	}

	return nil
}
