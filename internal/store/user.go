package store

import (
	"context"
	"database/sql"
	"time"

	sq "github.com/Masterminds/squirrel"
	"github.com/Planck1858/all-wallet/internal/consts"
	"github.com/pkg/errors"
)

func (s *store) GetUser(ctx context.Context, userId int32) (u *User, err error) {
	q := sq.Select("*").
		From(userTable).
		Where(sq.Eq{"telegram_id": userId}).
		PlaceholderFormat(sq.Dollar)

	err = s.db.GetSq(ctx, &u, q)
	if err != nil {
		return nil, err
	}

	return u, nil
}

func (s *store) CreateUser(ctx context.Context, u *User) (err error) {
	if u.Currency == "" {
		u.Currency = consts.DefaultCurrency
	}

	tx, err := s.db.Tx(ctx, &sql.TxOptions{})
	if err != nil {
		return err
	}
	defer func() {
		if err != nil {
			_ = tx.Rollback()
		}
	}()

	sqUser := sq.Insert(userTable).
		Columns("telegram_id", "total_money", "currency", "created_at").
		Values(u.TelegramId, u.TotalMoney, u.Currency, time.Now().UTC()).
		Suffix("returning id").
		PlaceholderFormat(sq.Dollar)

	var userId int32
	err = s.db.GetSqTx(tx, ctx, &userId, sqUser)
	if err != nil {
		return err
	}

	sqAccount := sq.Insert(accountTable).
		Columns("name", "type", "balance", "currency", "user_id").
		Values("account", consts.AccountTypeCard, 0, u.Currency, userId).
		PlaceholderFormat(sq.Dollar)

	err = s.db.InsertSqTx(tx, ctx, sqAccount)
	if err != nil {
		return err
	}

	sqSpending := sq.Insert(spendingTable).
		Columns("user_id", "total", "currency").
		Values(userId, 0, u.Currency).
		PlaceholderFormat(sq.Dollar)

	err = s.db.InsertSqTx(tx, ctx, sqSpending)
	if err != nil {
		return err
	}

	return tx.Commit()
}

func (s *store) ExistUser(ctx context.Context, userId int32) (ok bool, err error) {
	q := sq.Select("1").
		From(userTable).
		Where(sq.Eq{"telegram_id": userId}).
		PlaceholderFormat(sq.Dollar)

	err = s.db.GetSq(ctx, &ok, q)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return false, err
	}

	return ok, nil
}

func (s *store) DeleteUser(ctx context.Context, userTelegramId int32) error {
	tx, err := s.db.Tx(ctx, &sql.TxOptions{})
	if err != nil {
		return err
	}

	defer func() {
		if err != nil {
			_ = tx.Rollback()
		}
	}()

	// select user to delete
	qSelectUser := sq.Select("*").
		From(userTable).
		Where(sq.Eq{"telegram_id": userTelegramId}).
		PlaceholderFormat(sq.Dollar)

	var user User
	err = s.db.GetSqTx(tx, ctx, &user, qSelectUser)
	if err != nil {
		return err
	}

	// delete user spendigs
	qDelSpendRecords := sq.Delete(spendingTable).
		Where(sq.Eq{"user_id": userTelegramId}).
		PlaceholderFormat(sq.Dollar)

	err = s.db.DeleteSqTx(tx, ctx, qDelSpendRecords)
	if err != nil {
		return err
	}

	// delete user accounts
	qDelAccounts := sq.Delete(accountTable).
		Where(sq.Eq{"user_id": userTelegramId}).
		PlaceholderFormat(sq.Dollar)

	err = s.db.DeleteSqTx(tx, ctx, qDelAccounts)
	if err != nil {
		return err
	}

	// delete user
	qDelUser := sq.Delete(userTable).
		Where(sq.Eq{"telegram_id": userTelegramId}).
		PlaceholderFormat(sq.Dollar)

	err = s.db.DeleteSqTx(tx, ctx, qDelUser)
	if err != nil {
		return err
	}

	return tx.Commit()
}
