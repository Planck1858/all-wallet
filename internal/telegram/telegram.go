package telegram

import (
	"github.com/Planck1858/all-wallet/internal/currency_api"
	"github.com/Planck1858/all-wallet/internal/store"
	"go.uber.org/zap"
	"gopkg.in/telebot.v3"
)

type Telegram struct {
	log            *zap.SugaredLogger
	cache          Cache
	store          store.Store
	bot            *telebot.Bot
	currencyClient *currency_api.Client
}

func New(log *zap.SugaredLogger, store store.Store, bot *telebot.Bot, currencyClient *currency_api.Client) *Telegram {
	return &Telegram{
		log:            log,
		cache:          newCache(),
		store:          store,
		bot:            bot,
		currencyClient: currencyClient,
	}
}

func (t *Telegram) RegisterAllHandlers() {
	t.bot.Handle(startUrl, t.StartHandler)
	t.bot.Handle(stopUrl, t.StopHandler)
	t.bot.Handle(menuUrl, t.MenuHandler)
	t.bot.Handle(helpUrl, t.HelpHandler)
	t.bot.Handle(settingsUrl, t.SettingsHandler)

	t.bot.Handle(telebot.OnText, t.CustomTextHandler)

	t.bot.Handle(&menuBtn.Info, t.InfoHandler)
	// t.bot.Handle(&menuBtn.Accounts, t.AccountsHandler)
	t.bot.Handle(&menuBtn.Spendings, t.SpendingsHandler)
	t.bot.Handle(&menuBtn.ExchangeRates, t.ExchangeRatesHandler)

	t.bot.Handle(&spendingsBtn.ChangeSpendingsCurrency, t.ChangeSpendingsCurrencyHandler)
	t.bot.Handle(&spendingsBtn.SelectUsd, t.SelectUsdHandler)
	t.bot.Handle(&spendingsBtn.SelectEur, t.SelectEurHandler)
	t.bot.Handle(&spendingsBtn.SelectGbp, t.SelectGbpHandler)
	t.bot.Handle(&spendingsBtn.SelectRub, t.SelectRubHandler)
	t.bot.Handle(&spendingsBtn.SelectKzt, t.SelectKztHandler)
	t.bot.Handle(&spendingsBtn.SelectBtc, t.SelectBtcHandler)
	t.bot.Handle(&spendingsBtn.SelectAnotherCurrency, t.SelectAnotherCurrencyRequestHandler)
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

type Cache map[string]interface{}

func newCache() Cache {
	return map[string]interface{}{}
}

func (c *Cache) add(k string, v interface{}) {
	(*c)[k] = v
}

func (c *Cache) get(k string) (interface{}, bool) {
	v, f := (*c)[k]

	return v, f
}

func (c *Cache) delete(k string) {
	delete(*c, k)
}
