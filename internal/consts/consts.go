package consts

const (
	DefaultCurrency = "usd"

	FloatPrecision = 10
	TotalPrecision = 6
)

type AccountType string

const (
	AccountTypeCash AccountType = "cash"
	AccountTypeCard AccountType = "card"
)
