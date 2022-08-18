/*
create 100 instances of each resource for some users
program to display results for each respective type
labels file
*/

package pkg

import (
	"github.com/project-safari/zebra"
)

const NumLabels = 10

// Function to generate radom label pairs to be used for generation of resources.
//
// Returns zebra.Labels.
func CreateLabels() zebra.Labels {
	many := RandNum(NumLabels)
	codes := make(zebra.Labels, many+1)
	codes.Add("system.group", "default")

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

	for let := uint32(0); let < many; let++ {
		col := RandData(colors)
		let := RandData(letters)

		codes[let] = col
	}

	return codes
}
