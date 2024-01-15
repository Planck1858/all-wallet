package telegram

const startMsg = `🎉 <b>Welcome to all-wallet!</b> 🎉

This bot helps you to manage your finances. It can help you to:
- Keep an eye on your accounts (cards, cash, etc.)
- Keep your spendings and watch some statistics
- Watch Exchange rates (default is USD to EUR)

<b>Features</b>
- Bot supports over 450 currencies (fiat currencies and crypto currencies)
- Add records in any currency and watch your numbers in one currency
- Convert your accounts or spendings currency anytime
- And more (in the future :Ь)!

/help <b>to get started!</b>
/menu to see the main menu
/settings to set up your account
/stop to stop the bot and delete all your data

<b>TODO</b>
- Create sepparate accounts
- Show statistics for a day/week, etc.
- Send signals for exchange rates
- Download records in CSV/XLS format
`

const helpMsg = `ℹ️ <b>Help</b> ℹ️
- /start - start the bot
- /menu - see the main menu
- /settings - set up your account
- /stop - stop the bot and delete all your data

<i>Spendings</i>
- Add spendings record
- Clenup account records
- Set default spendings currency
- Show spendings total

How to add records to spendings
Template: <code>+/- [value] [currency(optional)] [date(optional)]</code>
<code>+ 100 eur</code> with current time in UTC
or
<code>+ 100 eur 25.12</code>
or
<code>- 100 eur 25.12.2020</code>
also you can send message with multiple values
<code>+ 100 eur 25.12
- 50 eur 25.12</code>

<i>Accounts - in progress...</i>
- Add account record
- Clenup account records
- Set default account currency
- Show account total

How to add records to account
Template: <code>++/-- [value] [currency(optional)] [date(optional)]</code>
<code>++ 100 eur</code> with current time in UTC
or
<code>++ 100 eur 25.12</code>
or
<code>-- 100 eur 25.12.2020</code>
`

const infoMsg = `💰 <b>Info</b> 💰
- Default currency: <b>%v</b>
- Total spendings: <b>%v</b>
- Total accounts: <b>%v</b>`
