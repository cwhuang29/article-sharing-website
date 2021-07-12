package utils

import (
	"fmt"
	"strconv"
)

const (
	parseIntErr    = "Input cannot be parsed to an integer."
	parsePosIntErr = "Input cannot be parsed to a positive integer."
)

func Str2Int(s string) (int, error) {
	num, err := strconv.Atoi(s)
	if err != nil {
		return 0, fmt.Errorf(parseIntErr)
	}

	return num, nil
}

func Str2PosInt(s string) (int, error) {
	num, err := strconv.Atoi(s)
	if err != nil || num <= 0 {
		return 0, fmt.Errorf(parsePosIntErr)
	}

	return num, nil
}
