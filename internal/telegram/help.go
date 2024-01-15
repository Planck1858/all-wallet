package telegram

import "gopkg.in/telebot.v3"

func (t *Telegram) HelpHandler(c telebot.Context) (err error) {
	userId := c.Sender().ID
	const method = "telegram.InfoHandler()"
	t.log.With("method", method, "userId", userId).Info("started...")
	defer func() {
		t.logDefer(method, err, c)
	}()

	button.Reply(menuBtn.ButtonPerRow()...)

	return c.Send(helpMsg, button, telebot.ModeHTML)
}
