package telegram

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/pkg/errors"
	"gopkg.in/telebot.v3"
)

type spending struct {
	amount   float64
	currency string
	date     time.Time
}

func (t *Telegram) SpendingsHandler(c telebot.Context) (err error) {
	userId := c.Sender().ID
	ctx := context.Background()

	method := "telegram.SpendingsHandler()"
	t.log.With("method", method, "userId", userId).Info("started...")
	defer func() {
		t.logDefer(method, err, c)
	}()

	repSpendings, err := t.store.GetLastSpendingRecords(ctx, int32(userId))
	if err != nil {
		return errors.Wrap(err, "can't get last spending records")
	}

	total, currency, err := t.store.GetSpendingTotal(ctx, int32(userId))
	if err != nil {
		return errors.Wrap(err, "can't get spending total")
	}

	spendingsStr := "Last 10 spendings:\n"
	for _, v := range repSpendings {
		spendingsStr += fmt.Sprintf("- <b>%.2f %v</b>, %v UTC\n", v.Amount, v.Currency, v.CreatedAt.Format(time.DateTime))
	}

	spendingsStr += fmt.Sprintf("\n<b>Total: %.2f %v</b>\n", total, currency)

	return c.Send(spendingsStr, emptyOpt, telebot.ModeHTML)
}

func (t *Telegram) AddSpendingHandler(c telebot.Context) (err error) {
	ctx := context.Background()

	method := "telegram.AddSpendingHandler()"
	t.log.With("method", method, "userId", c.Sender().ID).Info("started...")
	defer func() {
		t.logDefer(method, err, c)
	}()

	spend, err := getSpendingFromStr(c.Message().Text)
	if err != nil {
		return errors.Wrap(err, "invalid spending format")
	}

	err = t.store.AddSpending(ctx, c.Sender().ID, spend.amount, spend.currency, spend.date)
	if err != nil {
		return errors.Wrap(err, "can't add spending")
	}

	return nil
}

// func (t *Telegram) SubSpendingHandler(c telebot.Context) (err error) {
// 	ctx := context.Background()
//
// 	spend, err := getSpendingFromStr(c.Message().Text)
// 	if err != nil {
// 		return errors.Wrap(err, "invalid spending format")
// 	}
//
// 	if spend.currency == "" {
// 		spend.currency = consts.DefaultCurrency
// 	}
//
// 	err = t.store.AddSpending(ctx, c.Sender().ID, -spend.amount, spend.currency, spend.date)
// 	if err != nil {
// 		return errors.Wrap(err, "can't sub spending")
// 	}
//
// 	return c.Send("", emptyOpt)
// }

// getSpendingFromStr - str format should be "+ 123.4 currency date";
// operator can be +/-; currency is uds/eur/rub etc.; date is dd.mm.yyyy or dd.mm with current year
func getSpendingFromStr(str string) (spending spending, err error) {
	spendArr := strings.Split(str, " ")
	if !validateSpending(spendArr) {
		return spending, fmt.Errorf("invalid spending format: %s", str)
	}

	float, err := strconv.ParseFloat(strings.ReplaceAll(spendArr[1], ",", "."), 32)
	if err != nil {
		return spending, err
	}
	spending.amount = float

	if len(spendArr) == 3 {
		spending.currency = spendArr[2]
		spending.date = time.Now().UTC()
		return spending, nil
	}

	if len(spendArr) == 4 {
		spending.currency = spendArr[2]

		dataArr := strings.Split(spendArr[3], ".")
		if len(dataArr) == 2 {
			spending.date, err = time.Parse("02.01.2006", spendArr[3]+"."+fmt.Sprintf("%02d", time.Now().Year()))
			if err != nil {
				return spending, err
			}
		}

		if len(dataArr) == 3 {
			spending.date, err = time.Parse("02.01.2006", spendArr[3])
			if err != nil {
				return spending, err
			}
		}
	}

	return spending, nil
}

func validateSpending(spendArr []string) bool {
	if (len(spendArr) > 4 && len(spendArr) < 3) ||
		(spendArr[0] != "+" && spendArr[0] != "-") {
		return false
	}
	return true
}
