package telegram

var startMsg = `🎉 <b>Welcome to all-wallet!</b> 🎉

This bot helps you to manage your finances. It can help you to:
- Keep an eye on your accounts (cards, cash, etc.)
- Keep your spendings and watch some statistics
- Watch Exchange rates (default is USD)

<b>Features</b>
- Bot supports over 450 currencies (fiat currencies, crypto currencies, etc.)
- Add records in any currency and watch your numbers in one currency
- Convert your account or spendings currency anytime
- Watch statistics for a day/week and etc.
- and more!

Type /help to get started
Type /menu to see the main menu
Type /settings to set up your account
Type /stop to stop the bot
`

var helpMsg = `ℹ️ <b>Help</b> ℹ️
- /start - start the bot
- /menu - see the main menu
- /settings - set up your account
- /stop - stop the bot and delete all your data

<i>Accounts</i>
- Add new account
- Delete account
- Set default account currency
- Show account info

How to add records to account
<code>++ 100 card1</code>
or
<code>++ 100 eur card1</code>
Template: <code>++ [value] [currency(optional)] [account name]</code>

<i>Spendings</i>
- Add new spendings
- Delete spendings
- Set default spendings currency
- Show spendings info

<code>+ 100 tag1</code>
or
<code>+ 100 usd tag1</code>
Template: <code>++ [value] [currency(optional)] [tag name]</code>
`
