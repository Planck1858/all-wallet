package telegram

import (
	"fmt"
	"time"

	"gopkg.in/telebot.v3"
)

func (t *Telegram) ExchangeRatesHandler(c telebot.Context) (err error) {
	userId := c.Sender().ID

	method := "telegram.ExchangeRatesHandler()"
	t.log.With("method", method, "userId", userId).Info("started...")
	defer func() {
		t.logDefer(method, err, c)
	}()

	curr, err := t.currencyClient.GetCurrency("usd", "rub", time.Now().UTC())
	if err != nil {
		return err
	}

	return c.Send("usd to rub exchange rate: "+fmt.Sprintf("%.2f", curr)+"", emptyOpt)
}
