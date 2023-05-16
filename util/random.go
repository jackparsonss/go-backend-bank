package util

import (
	"math/rand"
	"strings"
	"time"
)

const alphabet = "abcdefghijklmnopqrstuvwxyz"

func init() {
	rand.Seed(time.Now().UnixNano())
}

// returns random number between min and max
func RandomInt(min, max int64) int64 {
	return min + rand.Int63n(max-min+1)
}

// returns random string of n characters
func RandomString(n int) string {
	var sb strings.Builder
	k := len(alphabet)

	for i := 0; i < n; i++ {
		c := alphabet[rand.Intn(k)]
		sb.WriteByte(c)
	}

	return sb.String()
}

// returns a random owner name
func RandomOwner() string {
	return RandomString(6)
}

// returns a random amount of money
func RandomMoney() int64 {
	return RandomInt(0, 1000)
}

// returns a random currency
func RandomCurrency() string {
	currencies := []string{CAD, USD, EUR}
	n := len(currencies)
	return currencies[rand.Intn(n)]
}
