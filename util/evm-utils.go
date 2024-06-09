package util

import (
	"regexp"
	"strconv"
	"strings"
)

func IsValidAddress(address string) bool {
	re := regexp.MustCompile("^0x[0-9a-fA-F]{40}$")
	return re.MatchString(address)
}

func HexToDecimal(hex string) (int64, error) {
	hex = strings.Replace(hex, "0x", "", -1)
	return strconv.ParseInt(hex, 16, 64)
}
