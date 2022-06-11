package util

const (
	USD = "USD"
	GBP = "GBP"
	EUR = "EUR"
)

var SupportedCurrencies = []string{
	USD, GBP, EUR,
}

// returns true if the currency is supported and false otherwise
func IsSupportedCurrency(currency string) bool {
	for _, k := range SupportedCurrencies {
		if currency == k {
			return true
		}
	}
	return false
}
