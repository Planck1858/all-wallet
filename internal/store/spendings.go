package store

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	sq "github.com/Masterminds/squirrel"
	"github.com/Planck1858/all-wallet/internal/consts"
	"github.com/Planck1858/all-wallet/internal/utils"
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
)

const selectUserId = "select id from " + userTable + " where telegram_id = ?"

func (s *store) BeginTx(ctx context.Context, txOpt *sql.TxOptions) (*sqlx.Tx, error) {
	return s.db.Tx(ctx, txOpt)
}

func (s *store) GetLastSpendingRecords(ctx context.Context, teleUserId int32) ([]SpendingRecord, error) {
	sqSelectSpendingId := sq.Select("*").From(spendingRecordTable).
		Where("spending_id = (select id from spending where user_id = ("+selectUserId+"))", teleUserId).
		OrderBy("id DESC").Limit(10).
		PlaceholderFormat(sq.Dollar)

	var spendings []SpendingRecord
	err := s.db.SelectSq(ctx, &spendings, sqSelectSpendingId)
	if err != nil {
		return nil, err
	}

	return spendings, nil
}

func (s *store) GetSpendingTotal(ctx context.Context, teleUserId int32) (float64, string, error) {
	sqTotal := sq.Select("total", "currency").From(spendingTable).
		Where("user_id = ("+selectUserId+")", teleUserId).PlaceholderFormat(sq.Dollar)

	res := struct {
		Total    float64 `db:"total"`
		Currency string  `db:"currency"`
	}{}

	err := s.db.GetSq(ctx, &res, sqTotal)
	if err != nil {
		return 0, "", err
	}

	return res.Total, res.Currency, nil
}

func (s *store) AddSpending(ctx context.Context, teleUserId int64, amount float64, currency string, date time.Time) (err error) {
	if date.IsZero() {
		date = time.Now().UTC()
	}

	tx, err := s.db.Tx(ctx, &sql.TxOptions{})
	if err != nil {
		return errors.Wrap(err, "start transaction")
	}
	defer func() {
		if err != nil {
			_ = tx.Rollback()
		}
	}()

	// select spending
	sqSelectSpending := sq.Select("*").From(spendingTable).
		Where("user_id = ("+selectUserId+")", teleUserId).PlaceholderFormat(sq.Dollar)

	var spending Spending
	err = s.db.GetSqTx(tx, ctx, &spending, sqSelectSpending)
	if err != nil {
		return errors.Wrap(err, "get spending")
	}

	// add spending
	sqAddSpending := sq.Insert(spendingRecordTable).
		Columns("amount", "currency", "created_at", "spending_id").
		Values(amount, currency, date, spending.Id).PlaceholderFormat(sq.Dollar)

	fmt.Println("amount: ", amount)
	err = s.db.InsertSqTx(tx, ctx, sqAddSpending)
	if err != nil {
		return errors.Wrap(err, "add spending")
	}

	// get spending currency
	sqSelectCurrency := sq.Select("currency").From(spendingTable).
		Where("user_id = ("+selectUserId+")", teleUserId).PlaceholderFormat(sq.Dollar)

	var spendingCurrency string
	err = s.db.GetSqTx(tx, ctx, &spendingCurrency, sqSelectCurrency)
	if err != nil {
		return errors.Wrap(err, "get spending currency")
	}

	// get currency rate
	if currency != spendingCurrency {
		rate, err := s.currencyApi.GetCurrency(currency, spendingCurrency, date)
		if err != nil {
			return errors.Wrap(err, "get currency rate")
		}

		amount *= rate
	}

	// update total
	sqUpdateTotal := sq.Update(spendingTable).Set("total", sq.Expr("total + ?", amount)).
		Where("user_id = ("+selectUserId+")", teleUserId).PlaceholderFormat(sq.Dollar)

	err = s.db.UpdateSqTx(tx, ctx, sqUpdateTotal)
	if err != nil {
		return errors.Wrap(err, "update total")
	}

	return tx.Commit()
}

func (s *store) AddSpendings(ctx context.Context, teleUserId int64, req []AddSpendingsReq) (err error) {
	var tx *sqlx.Tx
	tx, err = s.db.Tx(ctx, &sql.TxOptions{})
	if err != nil {
		return errors.Wrap(err, "start transaction")
	}
	defer func() {
		if err != nil {
			_ = tx.Rollback()
		}
	}()

	// select spending
	sqSelectSpending := sq.Select("*").From(spendingTable).
		Where("user_id = ("+selectUserId+")", teleUserId).PlaceholderFormat(sq.Dollar)

	var spending Spending
	err = s.db.GetSqTx(tx, ctx, &spending, sqSelectSpending)
	if err != nil {
		return errors.Wrap(err, "get spending")
	}

	// get spending currency
	sqSelectCurrency := sq.Select("currency").From(spendingTable).
		Where("user_id = ("+selectUserId+")", teleUserId).PlaceholderFormat(sq.Dollar)

	var spendingCurrency string
	err = s.db.GetSqTx(tx, ctx, &spendingCurrency, sqSelectCurrency)
	if err != nil {
		return errors.Wrap(err, "get spending currency")
	}

	for _, r := range req {
		// add spending
		sqAddSpending := sq.Insert(spendingRecordTable).
			Columns("amount", "currency", "created_at", "spending_id").
			Values(r.Amount, r.Currency, r.Date, spending.Id).PlaceholderFormat(sq.Dollar)

		err = s.db.InsertSqTx(tx, ctx, sqAddSpending)
		if err != nil {
			return errors.Wrap(err, "add spending")
		}

		// get currency rate
		if r.Currency != spendingCurrency {
			rate, err := s.currencyApi.GetCurrency(r.Currency, spendingCurrency, r.Date)
			if err != nil {
				return errors.Wrap(err, "get currency rate")
			}

			r.Amount *= rate
		}

		// update total
		sqUpdateTotal := sq.Update(spendingTable).Set("total", sq.Expr("total + ?", r.Amount)).
			Where("user_id = ("+selectUserId+")", teleUserId).PlaceholderFormat(sq.Dollar)

		err = s.db.UpdateSqTx(tx, ctx, sqUpdateTotal)
		if err != nil {
			return errors.Wrap(err, "update total")
		}
	}

	err = tx.Commit()
	if err != nil {
		return errors.Wrap(err, "commit transaction")
	}

	return nil
}

func (s *store) SetSpendingsDefaultCurrency(ctx context.Context, teleUserId int32, newCurrency string) (err error) {
	var tx *sqlx.Tx
	tx, err = s.db.Tx(ctx, &sql.TxOptions{})
	if err != nil {
		return errors.Wrap(err, "start transaction")
	}
	defer func() {
		if err != nil {
			_ = tx.Rollback()
		}
	}()

	// get old spending currency and total amount
	type currencyAndTotalDTO struct {
		Currency string  `db:"currency"`
		Total    float64 `db:"total"`
		UserId   int32   `db:"user_id"`
	}

	qGetOldCurrency := sq.Select("currency", "total", "user_id").From(spendingTable).
		Where("user_id = ("+selectUserId+")", teleUserId).
		PlaceholderFormat(sq.Dollar)

	var curTotal currencyAndTotalDTO
	err = s.db.GetSqTx(tx, ctx, &curTotal, qGetOldCurrency)
	if err != nil {
		return errors.Wrap(err, "get old currency")
	}

	// get exchange rate
	rate, err := s.currencyApi.GetCurrency(curTotal.Currency, newCurrency, time.Now().UTC())
	if err != nil {
		return errors.Wrap(err, "get currency rate")
	}

	newTotal := utils.RoundFloat(curTotal.Total*rate, consts.FloatPrecision)

	// update spending currency
	qUpdateCurrency := sq.Update(spendingTable).
		Set("currency", newCurrency).
		Set("total", newTotal).
		Where(sq.Eq{"user_id": curTotal.UserId}).
		PlaceholderFormat(sq.Dollar)

	err = s.db.UpdateSqTx(tx, ctx, qUpdateCurrency)
	if err != nil {
		return errors.Wrap(err, "update currency")
	}

	err = tx.Commit()
	if err != nil {
		return errors.Wrap(err, "commit transaction")
	}

	return nil
}
