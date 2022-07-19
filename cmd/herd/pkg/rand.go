package pkg

import (
	"crypto/rand"
	"encoding/binary"
)

// random selection from lists.
func RandData(res []string) string {
	length := len(res)

	ind := RandNum(length)

	typ := res[ind]

	return typ
}

// random numbers.
func RandNum(length int) uint32 {
	eight := 8

	// for safe, pseudo-random numbers.
	b := make([]byte, eight)
	n, err := rand.Read(b)

	// check panic.
	if n != eight {
		panic(n)
	} else if err != nil {
		panic(err)
	}

	ind := binary.BigEndian.Uint32(b) % uint32(length)

	return ind
}
