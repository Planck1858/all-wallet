package telegram

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/Planck1858/all-wallet/internal/consts"
	"github.com/Planck1858/all-wallet/internal/store"
	"github.com/Planck1858/all-wallet/internal/utils"
	"github.com/pkg/errors"
	"gopkg.in/telebot.v3"
)

type expanse struct {
	amount   float64
	currency string
	date     time.Time
}

func (t *Telegram) SpendingsHandler(c telebot.Context) (err error) {
	userId := c.Sender().ID
	ctx := context.Background()

	const method = "telegram.SpendingsHandler()"
	t.log.With("method", method, "userId", userId).Info("started...")
	defer func() {
		t.logDefer(method, err, c)
	}()

	repSpendings, err := t.store.GetLastSpendingRecords(ctx, int32(userId))
	if err != nil {
		return errors.Wrap(err, "get last spending records")
	}

	total, currency, err := t.store.GetSpendingTotal(ctx, int32(userId))
	if err != nil {
		return errors.Wrap(err, "get spending total")
	}

	spendingsStr := "Last 10 spendings:\n"
	for _, v := range repSpendings {
		spendingsStr += fmt.Sprintf("∙ <b>%v %v</b>, %v UTC\n", utils.FormatFloatToStr(v.Amount, consts.TotalPrecision), v.Currency, v.CreatedAt.Format(time.DateOnly))
	}

	spendingsStr += fmt.Sprintf("\n<b>Total: %v %v</b>\n", utils.FormatFloatToStr(total, consts.TotalPrecision), currency)

	var selectKeys [][]telebot.InlineButton
	selectKeys = append(selectKeys, []telebot.InlineButton{spendingsBtn.ChangeSpendingsCurrency})

	return c.Send(spendingsStr, &telebot.ReplyMarkup{InlineKeyboard: selectKeys}, telebot.ModeHTML)
}

func (t *Telegram) AddSpendingHandler(c telebot.Context) (err error) {
	ctx := context.Background()

	const method = "telegram.AddSpendingHandler()"
	t.log.With("method", method, "userId", c.Sender().ID).Info("started...")
	defer func() {
		t.logDefer(method, err, c)
	}()

	spendings, err := getExpanseArrFromStr(c.Message().Text)
	if err != nil {
		return errors.Wrap(err, "invalid spendings")
	}

	reqAddSpendReq := make([]store.AddSpendingsReq, 0, len(spendings))
	for _, spend := range spendings {
		reqAddSpendReq = append(reqAddSpendReq, store.AddSpendingsReq{
			Amount:   spend.amount,
			Currency: spend.currency,
			Date:     spend.date,
		})
	}

	err = t.store.AddSpendings(ctx, c.Sender().ID, reqAddSpendReq)
	if err != nil {
		return errors.Wrap(err, "add spending")
	}

	return c.Send("Spending added", emptyOpt)
}

func getExpanseArrFromStr(str string) ([]expanse, error) {
	if str == "" {
		return nil, errors.New("empty string")
	}

	spendingsArr := strings.Split(str, "\n")

	var spendings []expanse
	for _, v := range spendingsArr {
		spend, err := getExpanseFromStr(v)
		if err != nil {
			return nil, errors.Wrap(err, "invalid expanse string")
		}

		spendings = append(spendings, spend)
	}

	return spendings, nil
}

// getExpanseFromStr - str format should be "operator 123.4 currency date";
// operator can be +/-/++/--; currency is uds/eur/rub etc.; date is optional,
// can be dd.mm.yyyy or dd.mm with current year
func getExpanseFromStr(str string) (exp expanse, err error) {
	expArr := strings.Split(str, " ")
	if !validateExpanseStr(expArr) {
		return exp, fmt.Errorf("invalid expanse format: %s", str)
	}

	float, err := strconv.ParseFloat(strings.ReplaceAll(expArr[1], ",", "."), 64)
	if err != nil {
		return exp, errors.Wrap(err, "parse float")
	}

	if expArr[0] == "-" || expArr[0] == "--" {
		exp.amount = utils.RoundFloat(-float, consts.FloatPrecision)
	} else {
		exp.amount = utils.RoundFloat(float, consts.FloatPrecision)
	}

	switch len(expArr) {
	case 3:
		exp.currency = expArr[2]
		exp.date = time.Now().UTC().Round(time.Minute)
		return exp, nil

	case 4:
		exp.currency = expArr[2]

		dataArr := strings.Split(expArr[3], ".")
		if len(dataArr) == 2 {
			exp.date, err = time.Parse("02.01.2006", expArr[3]+"."+fmt.Sprintf("%02d", time.Now().Year()))
			if err != nil {
				return exp, errors.Wrap(err, "parse date")
			}
		}

		if len(dataArr) == 3 {
			exp.date, err = time.Parse("02.01.2006", expArr[3])
			if err != nil {
				return exp, errors.Wrap(err, "parse date")
			}
		}
	}

	return exp, nil
}

func validateExpanseStr(arr []string) bool {
	if (len(arr) > 4 && len(arr) < 3) ||
		((arr[0] != "+" && arr[0] != "-") &&
			(arr[0] != "++" && arr[0] != "--")) {
		return false
	}

	return true
}
