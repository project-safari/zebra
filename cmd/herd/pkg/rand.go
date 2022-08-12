package pkg

import (
	"crypto/rand"
	"encoding/binary"
)

// function that provides random selection from lists.
//
// takes in a string array to select from.
//
// returns a randomly selscted string from the array.
func RandData(res []string) string {
	length := len(res)

	ind := RandNum(length)

	typ := res[ind]

	return typ
}

// function for random numbers.
//
// takes in a integer that represents the maximum number to be genertaed.
//
// returns a random uint32 pseudo-random number.
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
