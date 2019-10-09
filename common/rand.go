package common

import (
	"crypto/rand"
	mrand "math/rand"
	"time"
)

var alphaNum = []byte(`0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz`)

// RandomBytes generate random []byte by specify chars.
func RandomBytes(n int, alphabets ...byte) []byte {
	if len(alphabets) == 0 {
		alphabets = alphaNum
	}
	var bytes = make([]byte, n)
	var randBy bool
	if num, err := rand.Read(bytes); num != n || err != nil {
		mrand.Seed(time.Now().UnixNano())
		randBy = true
	}
	for i, b := range bytes {
		if randBy {
			bytes[i] = alphabets[mrand.Intn(len(alphabets))]
		} else {
			bytes[i] = alphabets[b%byte(len(alphabets))]
		}
	}
	return bytes
}
