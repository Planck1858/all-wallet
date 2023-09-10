package telegram

import (
	"github.com/Planck1858/all-wallet/internal/currency_api"
	"github.com/Planck1858/all-wallet/internal/store"
	"go.uber.org/zap"
	"gopkg.in/telebot.v3"
)

type Telegram struct {
	log            *zap.SugaredLogger
	store          store.Store
	bot            *telebot.Bot
	currencyClient *currency_api.Client
}

func New(log *zap.SugaredLogger, store store.Store, bot *telebot.Bot, currencyClient *currency_api.Client) *Telegram {
	return &Telegram{
		log:            log,
		store:          store,
		bot:            bot,
		currencyClient: currencyClient,
	}
}

func (t *Telegram) RegisterAllHandlers() {
	t.bot.Handle(startUrl, t.StartHandler)
	t.bot.Handle(stopUrl, t.StopHandler)
	t.bot.Handle(settingsUrl, t.SettingsHandler)
	t.bot.Handle(helpUrl, t.HelpHandler)

	t.bot.Handle(menuUrl, t.MenuHandler)
	t.bot.Handle(&btnInfo, t.InfoHandler)
	t.bot.Handle(&btnAccounts, t.AccountsHandler)
	t.bot.Handle(&btnSpendings, t.SpendingsHandler)
	t.bot.Handle(&btnExchangeRates, t.ExchangeRatesHandler)

	t.bot.Handle(telebot.OnText, t.CustomTextHandler)
}

func (t *Telegram) Run() {
	t.bot.Start()
}

func (t *Telegram) Stop() {
	t.bot.Stop()
}

func (t *Telegram) logDefer(method string, err error, c telebot.Context) {
	if err != nil {
		t.log.With("method", method, "err", err).Error("finished with error")

		if c != nil {
			c.Send("Oops... error", emptyOpt)
		}
	} else {
		t.log.With("method", method).Info("finished successfully")
	}
}
