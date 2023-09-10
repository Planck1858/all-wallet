package telegram

import (
	"gopkg.in/telebot.v3"
)

func (t *Telegram) InfoHandler(c telebot.Context) (err error) {
	userId := c.Sender().ID

	method := "telegram.InfoHandler()"
	t.log.With("method", method, "userId", userId).Info("started...")
	defer func() {
		t.logDefer(method, err, c)
	}()

	btnMenu.Reply(
		btnMenu.Row(btnExchangeRates),
		btnMenu.Row(btnAccounts),
		btnMenu.Row(btnSpendings),
		btnMenu.Row(btnInfo),
	)

	return c.Send("info", btnMenu)
}
