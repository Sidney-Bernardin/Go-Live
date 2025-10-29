package src

import (
	"crypto/rand"
	"encoding/hex"
)

func MustRandomString(length int) string {
	b := make([]byte, length/2)
	_, err := rand.Read(b)
	if err != nil {
		panic(err)
	}
	return hex.EncodeToString(b)
}

func MapSlice[A, B any](aa []A, cb func(A) B) []B {
	bb := make([]B, len(aa))
	for i, a := range aa {
		bb[i] = cb(a)
	}
	return bb
}
