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

// Generate a randomInt
func RandomInt(min, max int64) int64 {
	return min + rand.Int63n(max-min+1)
}

// Generate a randomString with length n
func RandomString(n int) string {
	var sb strings.Builder
	k := len(alphabet)

	for i := 0; i < n; i++ {
		c := alphabet[rand.Intn(k)]
		sb.WriteByte(c)
	}

	return sb.String()
}

// Generate a random owner string name
func RandomOwner() string {
	return RandomString(6)
}

// Generate a random amount of money
func RandomMoney() int64 {
	return RandomInt(0, 1000)
}

// Generate a random currency code
func RandomCurrency() string {
	code := []string{"EUR", "USD", "CAD", "MLC"}
	l := len(code)
	return code[rand.Intn(l)]
}
