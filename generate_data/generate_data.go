package generate_data

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/google/uuid"
	"github.com/project-safari/zebra"
)

func User() map[string]string {
	fmt.Println("Test usernames")
	people := make(map[string]string)

	// take care of usernames
	name_list := []string{
		"user1",
		"user2",
		"user3",
		"manager",
		"admin",
		"guest",
		"another-user",
		" ",
	}

	password_list := []string{
		"pass1",
		"pass2",
		"pass3",
		"manager",
		"admin-pass",
		" ",
		"another-password",
		" ",
	}

	for i := 0; i < len(name_list); i++ {

		if name_list[i] != " " {
			name := name_list[i]
			password := password_list[i]
			//usernames = append(usernames, h.name)
			people[name] = password
		}
	}
	
	return people
}

func RandType(res []string) string {
	rand.Seed(time.Now().UnixNano())
	var length = len(res)
	var ind int = rand.Intn(length - 1)
	typ := res[ind]
	return typ
}

func CreateLabels() map[string]string {
	codes := make(map[string]string)
	col := " "
	let := " "
	colors := []string{"red", "yellow", "green", "blue", "white", "magenta", "black", "purple", "brown", "orange", "pink", "grey"}
	letters := []string{"alpha", "beta", "gamma", "delta", "epsilon", "eta", "theta", "Iota", "Kappa", "Lambda", "Mu", "Nu"}

	for i := 0; i < len(colors); i++ {
		col = RandType(colors)
		let = RandType(letters)
		codes[let] = col
	}

	return codes
}

func CreateBaseResource(resType string, labels zebra.Labels) *zebra.BaseResource {
	id := uuid.New().String()

	if resType == "" {
		resType = "BaseResource"
	}

	return &zebra.BaseResource{
		ID:     id,
		Type:   resType,
		Labels: labels,
	}

}

var resourceTypes = []string{"VLANPool", "Switch", "IPAddressPool", "Datacenter", "Lab", "Rack", "Server", "ESX", "VM", " ", "BaseResource", "NamedResource", "Credentials"}

func Generate_Data() {

	credentials := User()
	all := make(map[interface{}]interface{})
	//all := make(map[string]*zebra.BaseResource)
	for each := 0; each < 100; each++ {

		theType := RandType(resourceTypes)
		theLabels := CreateLabels()
		theData := CreateBaseResource(theType, theLabels)
		fmt.Println("The username, password, corresponding resource and its respective labels: ", theData)
		all[credentials] = theData
	}
}
