package store

import (
	pgxw "github.com/Planck1858/pgxwrapper"
	"go.uber.org/zap"

	"github.com/Planck1858/all-wallet/internal/currency_api"
)

// Store as lazymonostore (im lazy)
type store struct {
	log *zap.SugaredLogger
	db  pgxw.PgDatabase

	currencyApi *currency_api.Client
}

func New(log *zap.SugaredLogger, db pgxw.PgDatabase, currencyApi *currency_api.Client) *store {
	return &store{
		log:         log,
		db:          db,
		currencyApi: currencyApi,
	}
}
