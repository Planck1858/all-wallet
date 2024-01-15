package telegram

import (
	"context"
	"strconv"
	"time"

	"gopkg.in/telebot.v3"

	"github.com/Planck1858/all-wallet/internal/store"
)

func (t *Telegram) StartHandler(c telebot.Context) (err error) {
	ctx := context.Background()
	userId := c.Sender().ID

	const method = "telegram.StartHandler()"
	t.log.With("method", method, "userId", userId).Info("started...")
	defer func() {
		t.logDefer(method, err, c)
	}()

	var found bool
	found, err = t.store.ExistUser(ctx, int32(userId))
	if err != nil {
		return err
	}

	if !found {
		err = t.store.CreateUser(ctx, &store.User{
			Id:         0,
			TelegramId: strconv.Itoa(int(userId)),
			TotalMoney: 0,
			CreatedAt:  time.Now().UTC(),
		})
		if err != nil {
			return err
		}
	}

	button.Reply(menuBtn.ButtonPerRow()...)

	return c.Send(startMsg, button, telebot.ModeHTML)
}
