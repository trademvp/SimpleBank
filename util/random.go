package util

import (
	"fmt"
	"math/rand"
	"strings"
	"time"
)

const alphabet = "abcdefghijkmlnopqrstuvwxyz"

func init() {
	rand.NewSource(time.Now().UnixNano())
}

// RandomInt 生成介于min和max之间的随机整数
func RandomInt(min, max int64) int64 {
	return min + rand.Int63n(max-min+1)
}

// RandomString  生成一个长度为n的随机字符串
func RandomString(n int) string {
	var sb strings.Builder
	k := len(alphabet)

	for i := 0; i < n; i++ {
		c := alphabet[rand.Intn(k)]
		sb.WriteByte(c)
	}
	return sb.String()
}

// RandomOwner 随机owner
func RandomOwner() string {
	return RandomString(6)
}

// RandomMoney 随机money
func RandomMoney() int64 {
	return RandomInt(0, 1000)
}

// RandomCurrency  随机currency
func RandomCurrency() string {
	currencies := []string{EUR, USD, CAD, RMB}
	n := len(currencies)
	return currencies[rand.Intn(n)]
}

// RandomEmail 随机Email
func RandomEmail() string {
	return fmt.Sprintf("%s@email.com", RandomString(6))
}
