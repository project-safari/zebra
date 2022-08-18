package pkg

import (
	"crypto/rand"
	"encoding/binary"
)

// Function that provides random selection from lists.
//
// Takes in a string array to select from.
//
// Returns a randomly selscted string from the array.
func RandData(res []string) string {
	length := len(res)

	ind := RandNum(length)

	typ := res[ind]

	return typ
}

// Function for random numbers.
//
// Takes in a integer that represents the maximum number to be genertaed.
//
// Returns a random uint32 pseudo-random number.
func RandNum(length int) uint32 {
	four := 4

	// For safe, pseudo-random numbers.
	b := make([]byte, four)
	n, err := rand.Read(b)

	// Check panic.
	if n != four {
		panic(n)
	} else if err != nil {
		panic(err)
	}

	v := binary.BigEndian.Uint32(b)

	return v % uint32(length)
}
