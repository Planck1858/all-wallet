package telegram

import (
	"context"
	"time"

	"gopkg.in/telebot.v3"
)

func (t *Telegram) StopHandler(c telebot.Context) (err error) {
	userId := c.Sender().ID
	const method = "telegram.StopHandler()"
	t.log.With("method", method, "userId", userId).Info("started...")
	defer func() {
		t.logDefer(method, err, c)
	}()

	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	err = t.store.DeleteUser(ctx, int32(userId))
	if err != nil {
		return err
	}

	return c.Send("All your data have been deleted", emptyOpt)
}
