package telegram

import (
	"context"
	"time"

	"gopkg.in/telebot.v3"
)

// Commands
const (
	startUrl    = "/start"
	menuUrl     = "/menu"
	settingsUrl = "/settings"
	stopUrl     = "/stop"
	helpUrl     = "/help"
)

const (
	addAccounting      = "++ "
	subtractAccounting = "-- "

	addSpendings      = "+ "
	subtractSpendings = "- "
)

var emptyOpt = &telebot.SendOptions{}

// Buttons
var (
	// Main menu
	btnMenu = &telebot.ReplyMarkup{
		ResizeKeyboard:  true,
		OneTimeKeyboard: true,
		RemoveKeyboard:  true,
	}
	btnInfo          = btnMenu.Text("ℹ️ Info")
	btnAccounts      = btnMenu.Text("💰 Accounts")
	btnSpendings     = btnMenu.Text("🧾 Spendings")
	btnExchangeRates = btnMenu.Text("💹 Exchange Rates")

	// Exchange rates
	// Spendings
	// Accounts
)

func getTimeoutCtx() (context.Context, context.CancelFunc) {
	return context.WithTimeout(context.Background(), 5*time.Second)
}
