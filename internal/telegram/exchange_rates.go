package telegram

import (
	"fmt"
	"time"

	"gopkg.in/telebot.v3"
)

func (t *Telegram) ExchangeRatesHandler(c telebot.Context) (err error) {
	userId := c.Sender().ID

	const method = "telegram.ExchangeRatesHandler()"
	t.log.With("method", method, "userId", userId).Info("started...")
	defer func() {
		t.logDefer(method, err, c)
	}()

	// TODO
	curr, err := t.currencyClient.GetCurrency("usd", "rub", time.Now().UTC())
	if err != nil {
		return err
	}

	return c.Send("USD to RUB exchange rate: "+fmt.Sprintf("%.2f", curr)+"", emptyOpt)
}
