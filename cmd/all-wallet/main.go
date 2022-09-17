package main

import (
	"context"
	"encoding/json"
	"log"
	"time"

	"github.com/tidwall/pretty"

	tele "gopkg.in/telebot.v3"

	cfg "github.com/Planck1858/all-wallet/internal/config"
	"github.com/Planck1858/all-wallet/internal/tinkoff"
	"github.com/Planck1858/all-wallet/pkg/config"
)

func main() {
	ctx := context.Background()

	conf := cfg.Config{}
	err := config.NewConfig("config.yaml", &conf)
	if err != nil {
		log.Fatal(err)
	}

	tc := tinkoff.New(conf.TinkoffToken)

	bot, err := tele.NewBot(tele.Settings{
		Token:  conf.BotToken,
		Poller: &tele.LongPoller{Timeout: 10 * time.Second},
	})
	if err != nil {
		log.Fatal(err)
	}

	bot.Handle("/getAccounts", func(c tele.Context) error {
		log.Println("User send /getAccounts")

		acc, err := tc.GetAccounts(ctx)
		if err != nil {
			return err
		}

		accRaw, err := json.Marshal(acc)
		if err != nil {
			return err
		}
		accRaw = pretty.Pretty(accRaw)

		return c.Send(string(accRaw))
	})

	bot.Start()
}
