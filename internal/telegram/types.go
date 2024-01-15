package telegram

import (
	"gopkg.in/telebot.v3"
)

const (
	replyCacheKey = "reply"

	// Commands
	startUrl    = "/start"
	menuUrl     = "/menu"
	settingsUrl = "/settings"
	stopUrl     = "/stop"
	helpUrl     = "/help"

	// CustomMessagePrefixes
	addSpendingsPref       = "+ "
	subtractSpendingsPref  = "- "
	addAccountingPref      = "++ "
	subtractAccountingPref = "-- "

	// Currencies
	usd = "usd"
	eur = "eur"
	gbp = "gbp"
	rub = "rub"
	kzt = "kzt"
	btc = "btc"
)

var emptyOpt = &telebot.SendOptions{}

// Buttons
var (
	button = &telebot.ReplyMarkup{
		ResizeKeyboard:  true,
		OneTimeKeyboard: true,
		RemoveKeyboard:  true,
	}

	// Main menu
	menuBtn = menuButton{
		Info: button.Text("ℹ️ Info"),
		// Accounts:      button.Text("💰 Accounts"),
		Spendings:     button.Text("🧾 Spendings"),
		ExchangeRates: button.Text("💹 Exchange Rates"),
	}

	// Spendings
	spendingsBtn = spendingsButton{
		ChangeSpendingsCurrency: telebot.InlineButton{
			Unique: "change_spendings_currency", Text: "💸 Change spendings currency",
		},
		SelectUsd: telebot.InlineButton{Unique: "select_USD", Text: "$"},
		SelectEur: telebot.InlineButton{Unique: "select_EUR", Text: "€"},
		SelectGbp: telebot.InlineButton{Unique: "select_GBP", Text: "£"},
		SelectRub: telebot.InlineButton{Unique: "select_RUB", Text: "₽"},
		SelectKzt: telebot.InlineButton{Unique: "select_KZT", Text: "₸"},
		SelectBtc: telebot.InlineButton{Unique: "select_BTC", Text: "₿"},
		SelectAnotherCurrency: telebot.InlineButton{
			Unique: "select_another_currency", Text: "🔁 Select another currency",
		},
	}

	// Accounts

	// Exchange rates
)

type menuButton struct {
	Info          telebot.Btn
	Accounts      telebot.Btn
	Spendings     telebot.Btn
	ExchangeRates telebot.Btn
}

func (b menuButton) ButtonPerRow() []telebot.Row {
	return []telebot.Row{
		button.Row(b.ExchangeRates),
		button.Row(b.Accounts),
		button.Row(b.Spendings),
		button.Row(b.Info),
	}
}

type spendingsButton struct {
	ChangeSpendingsCurrency telebot.InlineButton
	SelectUsd               telebot.InlineButton
	SelectEur               telebot.InlineButton
	SelectGbp               telebot.InlineButton
	SelectRub               telebot.InlineButton
	SelectKzt               telebot.InlineButton
	SelectBtc               telebot.InlineButton
	SelectAnotherCurrency   telebot.InlineButton
}

func (b spendingsButton) Inline() []telebot.InlineButton {
	return []telebot.InlineButton{
		b.ChangeSpendingsCurrency,
		b.SelectUsd,
		b.SelectEur,
		b.SelectGbp,
		b.SelectRub,
		b.SelectKzt,
		b.SelectBtc,
		b.SelectAnotherCurrency,
	}
}
