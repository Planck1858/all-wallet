package utils

import (
	"fmt"
	"math"
	"strconv"
	"strings"
)

func RoundFloat(val float64, precision int) float64 {
	ratio := math.Pow(10, float64(precision))
	return math.Round(val*ratio) / ratio
}

func FormatFloatToStr(v float64, precision int) string {
	val := fmt.Sprintf("%f", RoundFloat(v, precision))

	parts := strings.Split(val, ".")
	if len(parts) != 2 {
		return val
	}

	intPart := parts[0]
	decimalPart := parts[1]

	trimmedDecimal := strings.TrimRight(decimalPart, "0")

	if len(trimmedDecimal) == 0 {
		trimmedDecimal = strings.Repeat("0", precision)
	} else {
		if len(trimmedDecimal) < precision {
			trimmedDecimal += strings.Repeat("0", precision-len(trimmedDecimal))
		} else {
			trimmedDecimal = trimmedDecimal[:precision]
		}
	}

	if trimmedDecimal == strings.Repeat("0", precision) {
		return intPart + ".00"
	}

	floatStr := intPart + "." + trimmedDecimal

	float, err := strconv.ParseFloat(floatStr, 64)
	if err != nil {
		return val
	}

	precisionStr := fmt.Sprintf("%v", precision)

	return fmt.Sprintf("%."+precisionStr+"f", float)
}
