package telegram

import (
	"strings"

	"gopkg.in/telebot.v3"
)

func (t *Telegram) CustomTextHandler(c telebot.Context) (err error) {
	text := strings.TrimSpace(c.Text())
	t.log.Info(text)

	if strings.HasPrefix(text, "/") {
		return nil
	}

	// TODO:
	// if strings.HasPrefix(strings.TrimSpace(text), addAccounting) {
	// 	return t.AddAccountingHandler(c)
	// } else if strings.HasPrefix(text, subtractAccounting) {
	// 	return t.SubAccountingHandler(c)
	// }

	if strings.HasPrefix(strings.TrimSpace(text), addSpendings) {
		return t.AddSpendingHandler(c)
	}
	// } else if strings.HasPrefix(text, subtractSpendings) {
	// return t.SubSpendingHandler(c)
	// }

	return nil
}
