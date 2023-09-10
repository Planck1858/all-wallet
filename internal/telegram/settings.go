package telegram

import (
	"gopkg.in/telebot.v3"
)

// SettingsHandler
// enable/disable write spendings
func (t *Telegram) SettingsHandler(c telebot.Context) (err error) {
	userId := c.Sender().ID

	method := "telegram.SettingsHandler()"
	t.log.With("method", method, "userId", userId).Info("started...")
	defer func() {
		t.logDefer(method, err, c)
	}()

	return c.Send("Settings", emptyOpt)
}
