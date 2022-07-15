package pkg

import (
	"math/rand"
	"time"
)

// random selection from lists.
func RandData(res []string) string {
	length := len(res)

	rand.Seed(time.Now().UnixNano())

	var ind int = rand.Intn(length - 1) //nolint // random selection for data sampling
	typ := res[ind]

	return typ
}
