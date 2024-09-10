package utils

import (
	"github.com/Ritwiksrivastava0809/go-bank/pkg/constants"
)

func ISSupportedCurrency(currency string) bool {
	// check if the currency is valid
	switch currency {
	case constants.USD, constants.EUR, constants.INR, constants.CAD, constants.YEN, constants.YUAN, constants.RUB, constants.PESO:
		return true
	}
	return false
}
