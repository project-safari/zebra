/*
create 100 instances of each resource for some users
program to display results for each respective type
*/

package pkg

import (
	"math/rand"
)

func CreateLabels() map[string]string {
	many := rand.Int()%10 + 1 //nolint

	codes := make(map[string]string, many)

	colors := []string{
		"red", "yellow", "green",
		"blue", "white", "magenta",
		"black", "purple", "brown",
		"orange", "pink", "grey",
	}

	letters := []string{
		"alpha", "beta", "gamma",
		"delta", "epsilon", "eta",
		"theta", "Iota", "Kappa",
		"Lambda", "Mu", "Nu",
	}

	for let := 0; let < many; let++ {
		col := RandData(colors)
		let := RandData(letters)

		codes[let] = col
	}

	return codes
}
