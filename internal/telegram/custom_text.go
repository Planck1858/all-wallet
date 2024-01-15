package telegram

import (
	"fmt"
	"strings"

	"gopkg.in/telebot.v3"
)

func (t *Telegram) CustomTextHandler(c telebot.Context) (err error) {
	const method = "telegram.CustomTextHandler()"

	text := strings.TrimSpace(c.Text())
	t.log.Info(method + ": " + text)

	if strings.HasPrefix(text, "/") {
		return nil
	}

	if c.Message().IsReply() {
		key := fmt.Sprintf("%v-%v", replyCacheKey, c.Sender().ID)

		_, found := t.cache.get(key)
		if found {
			defer t.cache.delete(key)

			return t.SelectAnotherCurrencyResponseHandler(c)
		}
	}

	// TODO
	// if strings.HasPrefix(strings.TrimSpace(text), addAccounting) || strings.HasPrefix(text, subtractAccounting) {
	// 	return t.AddAccountingHandler(c)
	// }

	if strings.HasPrefix(strings.TrimSpace(text), addSpendingsPref) || strings.HasPrefix(text, subtractSpendingsPref) {
		return t.AddSpendingHandler(c)
	}

	return nil
}
