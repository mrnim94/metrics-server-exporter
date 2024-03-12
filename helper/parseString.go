package helper

import (
	"metrics-server-exporter/log"
	"strconv"
	"strings"
)

// milli-unit or not
// convertStringToNumber converts a string to a float64.
// If the string ends with 'm', it removes 'm', converts the remaining part to a number, and divides by 1000.
func ConvertStringToNumber(input string) (float64, error) {
	var result float64
	var err error

	if strings.HasSuffix(input, "m") {
		numberPart := strings.TrimSuffix(input, "m")
		result, err = strconv.ParseFloat(numberPart, 64)
		if err != nil {
			log.Error(err)
			return 0, err
		}
		result /= 1000 // Convert to milli-unit
	} else {
		result, err = strconv.ParseFloat(input, 64)
		if err != nil {
			log.Error(err)
			return 0, err
		}
	}

	return result, nil
}
