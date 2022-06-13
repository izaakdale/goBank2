package util

import (
	"math/rand"
	"strings"
	"time"
)

var alphabet = "abcdefghijklmnopqrstuvwxyz"

func init() {
	rand.Seed(time.Now().UnixNano())
}

func RandomInt(min, max int64) int64 {
	return min + rand.Int63n(max-min+1)
}

func RandomString(n int) string {
	var sb strings.Builder

	k := len(alphabet)
	for i := 0; i < n; i++ {
		c := alphabet[rand.Intn(k)]
		sb.WriteByte(c)
	}
	return sb.String()
}

func RandomName() string {
	return RandomString(6) + " " + RandomString(8)
}

func RandomBalance() int64 {
	return RandomInt(1000, 5000)
}

func RandomCurrency() string {
	return SupportedCurrencies[rand.Intn(len(SupportedCurrencies))]
}

func RandomEmail() string {
	return RandomString(6) + "@" + RandomString(4) + ".com"
}
