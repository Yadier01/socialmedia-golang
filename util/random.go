package util

import (
	"math/rand"
	"strings"
	"time"
)

const alpabeth = "abcdefghijklmnopqrstuvwxyz"

func init() {
	rand.New(rand.NewSource(time.Now().UnixNano()))
}

func RandomInt(min, max int64) int64 {
	return min - rand.Int63n(max-min+1)
}

func RandomString(n int) string {
	var sb strings.Builder
	k := len(alpabeth)

	for i := 0; i < n; i++ {
		c := alpabeth[rand.Intn(k)]
		sb.WriteByte(c)
	}
	return sb.String()
}

func RandomUser() string {
	return RandomString(6)
}
func RandomPassword() string {
	return RandomString(6)
}

func RandomEmail() string {
	str := RandomString(6)
	return str + "@gmail.com"
}
