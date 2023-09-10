package store

import (
	"context"
	"database/sql"
	"time"

	sq "github.com/Masterminds/squirrel"
	"github.com/pkg/errors"
)

func (s *store) GetLastSpendingRecords(ctx context.Context, userId int32) ([]SpendingRecord, error) {
	sqSelectSpendingId := sq.Select("*").From(spendingRecordTable).
		Where("spending_id = (select id from spending where user_id = (select id from users where telegram_id = ?))", userId).
		OrderBy("id DESC").Limit(10).
		PlaceholderFormat(sq.Dollar)

	var spendings []SpendingRecord
	err := s.db.SelectSq(ctx, &spendings, sqSelectSpendingId)
	if err != nil {
		return nil, err
	}

	return spendings, nil
}

func (s *store) GetSpendingTotal(ctx context.Context, userId int32) (float64, string, error) {
	sqTotal := sq.Select("total", "currency").From(spendingTable).
		Where("user_id = (select id from users where telegram_id = ?)", userId).PlaceholderFormat(sq.Dollar)

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

func (s *store) AddSpending(ctx context.Context, userId int64, amount float64, currency string, date time.Time) (err error) {
	if date.IsZero() {
		date = time.Now().UTC()
	}

	tx, err := s.db.Tx(ctx, &sql.TxOptions{})
	if err != nil {
		return errors.Wrap(err, "can't start transaction")
	}
	defer func() {
		if err != nil {
			_ = tx.Rollback()
		}
	}()

	sqUserId, sqUserIdArgs, err := sq.Select("id").From(userTable).Where("telegram_id = ?", userId).PlaceholderFormat(sq.Dollar).ToSql()
	if err != nil {
		return errors.Wrap(err, "can't build get user id query")
	}

	// select spending
	sqSelectSpending := sq.Select("*").From(spendingTable).
		Where("user_id = ("+sqUserId+")", sqUserIdArgs...).PlaceholderFormat(sq.Dollar)

	var spending Spending
	err = s.db.GetSqTx(tx, ctx, &spending, sqSelectSpending)
	if err != nil {
		return errors.Wrap(err, "can't get spending")
	}

	// add spending
	sqAddSpending := sq.Insert(spendingRecordTable).
		Columns("amount", "currency", "created_at", "spending_id").
		Values(amount, currency, date, spending.Id).PlaceholderFormat(sq.Dollar)

	err = s.db.InsertSqTx(tx, ctx, sqAddSpending)
	if err != nil {
		return errors.Wrap(err, "can't add spending")
	}

	// get spending currency
	sqSelectCurrency := sq.Select("currency").From(spendingTable).
		Where("user_id = ("+sqUserId+")", sqUserIdArgs...).PlaceholderFormat(sq.Dollar)

	var spendingCurrency string
	err = s.db.GetSqTx(tx, ctx, &spendingCurrency, sqSelectCurrency)
	if err != nil {
		return errors.Wrap(err, "can't get spending currency")
	}

	// get currency rate
	if currency != spendingCurrency {
		rate, err := s.currencyApi.GetCurrency(currency, spendingCurrency, date)
		if err != nil {
			return errors.Wrap(err, "can't get currency rate")
		}

		amount *= rate
	}

	// update total
	sqUpdateTotal := sq.Update(spendingTable).Set("total", sq.Expr("total + ?", amount)).
		Where("user_id = (select id from users where telegram_id = ?)", userId).PlaceholderFormat(sq.Dollar)

	err = s.db.UpdateSqTx(tx, ctx, sqUpdateTotal)
	if err != nil {
		return errors.Wrap(err, "can't update total")
	}

	return tx.Commit()
}
