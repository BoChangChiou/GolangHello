package util

import (
	"strconv"

	"github.com/shopspring/decimal"
)

func StringToInt64(s string) (int64, error) {
	return strconv.ParseInt(s, 0, 64)
}

func StringToFloat64(s string) (float64, error) {
	return strconv.ParseFloat(s, 64)
}

func StringToBool(s string) (bool, error) {
	return strconv.ParseBool(s)
}

func StringToDecimal(s string) (decimal.Decimal, error) {
	return decimal.NewFromString(s)
}

func Int64ToInt(i int64) int {
	return int(i)
}
