package telegram

import (
	"gopkg.in/telebot.v3"
)

func (t *Telegram) AccountsHandler(c telebot.Context) (err error) {
	userId := c.Sender().ID

	const method = "telegram.AccountsHandler()"
	t.log.With("method", method, "userId", userId).Info("started...")
	defer func() {
		t.logDefer(method, err, c)
	}()

	// TODO

	return c.Send("Accounts", emptyOpt)
}

func (t *Telegram) AddAccountingHandler(c telebot.Context) (err error) {
	// ctx := context.Background()

	const method = "telegram.AddAccountingHandler()"
	t.log.With("method", method, "userId", c.Sender().ID).Info("started...")
	defer func() {
		t.logDefer(method, err, c)
	}()

	return c.Send("Accounting added", emptyOpt)
}
