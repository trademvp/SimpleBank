package util

const (
	USD = "USD"
	EUR = "EUR"
	CAD = "CAD"
	RMB = "RMB"
)

// IsSupportedCurrency 支持的货币
func IsSupportedCurrency(currency string) bool {
	switch currency {
	case USD,
		EUR,
		RMB,
		CAD:
		return true
	}
	return false
}
