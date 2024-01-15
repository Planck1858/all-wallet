package telegram

import (
	"context"
	"fmt"

	"github.com/pkg/errors"
	"gopkg.in/telebot.v3"
)

func (t *Telegram) ChangeSpendingsCurrencyHandler(c telebot.Context) (err error) {
	const method = "telegram.ChangeSpendingsCurrencyHandler()"
	t.log.With("method", method, "userId", c.Sender().ID).Info("started...")
	defer func() {
		t.logDefer(method, err, c)
	}()

	selectKeys := make([][]telebot.InlineButton, 0)
	selectKeys = append(selectKeys, spendingsBtn.Inline()[1:len(spendingsBtn.Inline())-1])
	selectKeys = append(selectKeys, []telebot.InlineButton{spendingsBtn.SelectAnotherCurrency})

	return c.Send("Choose currency", &telebot.ReplyMarkup{InlineKeyboard: selectKeys})
}

func (t *Telegram) SelectAnotherCurrencyRequestHandler(c telebot.Context) (err error) {
	const method = "telegram.SelectAnotherCurrencyRequestHandler()"
	t.log.With("method", method, "userId", c.Sender().ID).Info("started...")
	defer func() {
		t.logDefer(method, err, c)
	}()

	newMsg, err := t.bot.Send(c.Sender(), "Write currency index as a reply to this message", &telebot.SendOptions{
		ReplyTo: c.Message(),
	})
	if err != nil {
		return errors.Wrap(err, "send message")
	}

	fmt.Printf("after msg id: %v \n", newMsg.ID)
	k := fmt.Sprintf("%v-%v", replyCacheKey, c.Sender().ID)
	t.cache.add(k, newMsg.ID)

	return nil
}

func (t *Telegram) SelectAnotherCurrencyResponseHandler(c telebot.Context) (err error) {
	const method = "telegram.SelectAnotherCurrencyResponseHandler()"
	t.log.With("method", method, "userId", c.Sender().ID).Info("started...")
	defer func() {
		t.logDefer(method, err, c)
	}()

	cur := c.Message().Text
	if cur == "" {
		return c.Send("Currency can't be empty", emptyOpt)
	}

	foundCur := t.currencyClient.CheckCurrency(cur)
	if !foundCur {
		return c.Send("Currency not found", emptyOpt)
	}

	err = t.store.SetSpendingsDefaultCurrency(context.Background(), int32(c.Sender().ID), cur)
	if err != nil {
		return errors.Wrap(err, "set spendings currency")
	}

	err = c.Send("Currency selected", emptyOpt)
	if err != nil {
		return errors.Wrap(err, "send message")
	}

	return t.SpendingsHandler(c)
}

func (t *Telegram) SelectUsdHandler(c telebot.Context) (err error) {
	ctx := context.Background()
	const method = "telegram.SelectUsdHandler()"
	t.log.With("method", method, "userId", c.Sender().ID).Info("started...")
	defer func() {
		t.logDefer(method, err, c)
	}()

	err = t.store.SetSpendingsDefaultCurrency(ctx, int32(c.Sender().ID), usd)
	if err != nil {
		return errors.Wrap(err, "set spendings currency")
	}

	err = c.Send("USD selected as spendings currency", emptyOpt)
	if err != nil {
		return errors.Wrap(err, "send message")
	}

	return t.SpendingsHandler(c)
}

func (t *Telegram) SelectEurHandler(c telebot.Context) (err error) {
	ctx := context.Background()
	const method = "telegram.SelectEurHandler()"
	t.log.With("method", method, "userId", c.Sender().ID).Info("started...")
	defer func() {
		t.logDefer(method, err, c)
	}()

	err = t.store.SetSpendingsDefaultCurrency(ctx, int32(c.Sender().ID), eur)
	if err != nil {
		return errors.Wrap(err, "set spendings currency")
	}

	err = c.Send("EUR selected as spendings currency", emptyOpt)
	if err != nil {
		return errors.Wrap(err, "send message")
	}

	return t.SpendingsHandler(c)
}

func (t *Telegram) SelectGbpHandler(c telebot.Context) (err error) {
	ctx := context.Background()
	const method = "telegram.SelectGbpHandler()"
	t.log.With("method", method, "userId", c.Sender().ID).Info("started...")
	defer func() {
		t.logDefer(method, err, c)
	}()

	err = t.store.SetSpendingsDefaultCurrency(ctx, int32(c.Sender().ID), gbp)
	if err != nil {
		return errors.Wrap(err, "set spendings currency")
	}

	err = c.Send("GBP selected as spendings currency", emptyOpt)
	if err != nil {
		return errors.Wrap(err, "send message")
	}

	return t.SpendingsHandler(c)
}

func (t *Telegram) SelectRubHandler(c telebot.Context) (err error) {
	ctx := context.Background()
	const method = "telegram.SelectRubHandler()"
	t.log.With("method", method, "userId", c.Sender().ID).Info("started...")
	defer func() {
		t.logDefer(method, err, c)
	}()

	err = t.store.SetSpendingsDefaultCurrency(ctx, int32(c.Sender().ID), rub)
	if err != nil {
		return errors.Wrap(err, "set spendings currency")
	}

	err = c.Send("RUB selected as spendings currency", emptyOpt)
	if err != nil {
		return errors.Wrap(err, "send message")
	}

	return t.SpendingsHandler(c)
}

func (t *Telegram) SelectKztHandler(c telebot.Context) (err error) {
	ctx := context.Background()
	const method = "telegram.SelectKztHandler()"
	t.log.With("method", method, "userId", c.Sender().ID).Info("started...")
	defer func() {
		t.logDefer(method, err, c)
	}()

	err = t.store.SetSpendingsDefaultCurrency(ctx, int32(c.Sender().ID), kzt)
	if err != nil {
		return errors.Wrap(err, "set spendings currency")
	}

	err = c.Send("KZT selected as spendings currency", emptyOpt)
	if err != nil {
		return errors.Wrap(err, "send message")
	}

	return t.SpendingsHandler(c)
}

func (t *Telegram) SelectBtcHandler(c telebot.Context) (err error) {
	ctx := context.Background()
	const method = "telegram.SelectBtcHandler()"
	t.log.With("method", method, "userId", c.Sender().ID).Info("started...")
	defer func() {
		t.logDefer(method, err, c)
	}()

	err = t.store.SetSpendingsDefaultCurrency(ctx, int32(c.Sender().ID), btc)
	if err != nil {
		return errors.Wrap(err, "set spendings currency")
	}

	err = c.Send("BTC selected as spendings currency", emptyOpt)
	if err != nil {
		return errors.Wrap(err, "send message")
	}

	return t.SpendingsHandler(c)
}
