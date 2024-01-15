package store

import (
	"context"
	"database/sql"
	"time"

	"github.com/Planck1858/all-wallet/internal/consts"
	"github.com/jmoiron/sqlx"
)

const (
	userTable           = "users"
	accountTable        = "account"
	spendingTable       = "spending"
	spendingRecordTable = "spending_record"
)

type Store interface {
	BeginTx(ctx context.Context, txOpt *sql.TxOptions) (*sqlx.Tx, error)

	GetUser(ctx context.Context, userId int32) (u *User, err error)
	CreateUser(ctx context.Context, u *User) error
	ExistUser(ctx context.Context, userId int32) (ok bool, err error)
	DeleteUser(ctx context.Context, userId int32) error

	GetLastSpendingRecords(ctx context.Context, userId int32) ([]SpendingRecord, error)
	GetSpendingTotal(ctx context.Context, userId int32) (total float64, currency string, err error)
	AddSpending(ctx context.Context, userId int64, amount float64, currency string, date time.Time) error
	AddSpendings(ctx context.Context, userId int64, req []AddSpendingsReq) error
	SetSpendingsDefaultCurrency(ctx context.Context, userId int32, currency string) (err error)
}

type User struct {
	Id         int32     `db:"id"`
	TelegramId string    `db:"telegram_id"`
	TotalMoney float64   `db:"total_money"`
	Currency   string    `db:"currency"`
	CreatedAt  time.Time `db:"created_at"`
}

type Account struct {
	Id       int32              `db:"id"`
	Name     string             `db:"name"`
	Type     consts.AccountType `db:"type"`
	Balance  float64            `db:"balance"`
	Currency string             `db:"currency"`
	UserID   int32              `db:"user_id"`
}

type Spending struct {
	Id       int32            `db:"id"`
	Total    float64          `db:"total"`
	Currency string           `db:"currency"`
	UserID   int32            `db:"user_id"`
	Records  []SpendingRecord `db:"-"`
}

type AddSpendingsReq struct {
	Amount   float64
	Currency string
	Date     time.Time
}

type SpendingRecord struct {
	Id         int32     `db:"id"`
	Amount     float64   `db:"amount"`
	Currency   string    `db:"currency"`
	SpendingId int32     `db:"spending_id"`
	CreatedAt  time.Time `db:"created_at"`
}
