package telegram

import (
	"context"
	"fmt"

	"github.com/Planck1858/all-wallet/internal/consts"
	"github.com/Planck1858/all-wallet/internal/utils"
	"github.com/pkg/errors"
	"gopkg.in/telebot.v3"
)

func (t *Telegram) InfoHandler(c telebot.Context) (err error) {
	ctx := context.Background()
	userId := c.Sender().ID

	const method = "telegram.InfoHandler()"
	t.log.With("method", method, "userId", userId).Info("started...")
	defer func() {
		t.logDefer(method, err, c)
	}()

	total, currency, err := t.store.GetSpendingTotal(ctx, int32(userId))
	if err != nil {
		return errors.Wrap(err, "get spending total")
	}

	info := fmt.Sprintf(infoMsg, currency, utils.FormatFloatToStr(total, consts.TotalPrecision)+" "+currency, 0)

	return c.Send(info, button, telebot.ModeHTML)
}
