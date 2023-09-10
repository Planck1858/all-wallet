package telegram

import (
	"gopkg.in/telebot.v3"
)

func (t *Telegram) AccountsHandler(c telebot.Context) (err error) {
	userId := c.Sender().ID

	method := "telegram.AccountsHandler()"
	t.log.With("method", method, "userId", userId).Info("started...")
	defer func() {
		t.logDefer(method, err, c)
	}()

	return c.Send("Accounts", emptyOpt)
}
