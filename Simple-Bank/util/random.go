package util

import (
	"math/rand"
	"strconv"
	"strings"
	"time"
)

const alphabet = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

func init() {
	rand.Seed(time.Now().UnixNano())
}

func RandomInt(min, max int64) int64 {
	return min + rand.Int63n(max-min+1)
}

func RandomPhoneNumber() string {
	st := "77"
	for i := 0; i < 9; i++ {
		st += strconv.Itoa(int(RandomInt(0, 9)))
	}
	return st
}

func RandomFullName() string {
	return RandomOwner() + " " + RandomOwner()
}

func RandomEmail() string {
	return RandomOwner() + "@" + RandomOwner() + "." + RandomOwner()
}

func RandomString(n int) string {
	var st strings.Builder
	k := len(alphabet)
	for i := 0; i < n; i++ {
		c := alphabet[rand.Intn(k)]
		st.WriteByte(c)
	}
	return st.String()
}

func RandomOwner() string {
	return RandomString(6)
}

func RandomMoney() int64 {
	return RandomInt(0, 1000)
}

func RandomCurrency() string {
	currencies := []string{"USD", "EUR", "KZT"}
	n := len(currencies)
	return currencies[rand.Intn(n)]
}
