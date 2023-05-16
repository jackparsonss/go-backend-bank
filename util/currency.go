package util

const (
	USD = "USD"
	EUR = "EUR"
	CAD = "CAD"
)

// The function checks if a given currency is supported and returns a boolean value.
func IsSupportedCurrency(currency string) bool {
	switch currency {
	case USD, EUR, CAD:
		return true
	}
	return false
}
