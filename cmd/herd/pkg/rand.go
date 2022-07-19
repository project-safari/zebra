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
	four := 4

	// for safe, pseudo-random numbers.
	b := make([]byte, four)
	n, err := rand.Read(b)

	// check panic.
	if n != four {
		panic(n)
	} else if err != nil {
		panic(err)
	}

	v := binary.BigEndian.Uint32(b)

	return v % uint32(length)
}
