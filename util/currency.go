package util

// Constants all supported currencies
const (
	USD = "USD"
	EUR = "EUR"
	CAD = "CAD"
)

func IsSupported(currency string) bool {
	switch currency {
	case USD, CAD, EUR:
		return true
	}
	return false
}
