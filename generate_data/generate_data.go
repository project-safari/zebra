package generate_data

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/google/uuid"
	"github.com/project-safari/zebra"
)

/*type Human struct {
	name     string
	password string
}
*/

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
	//usernames := []string{}

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

	/*
		//passwords := []string{}
		for i := 0; i < len(password_list); i++ {

			if password_list[i] != " " {
				h.password = password_list[i]
				//passwords = append(passwords, h.password)
			}
		}
	*/
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

func CreateResource(resType string, labels zebra.Labels) *zebra.BaseResource {
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

//var Ids = []string{"01000000", "01000001", "01000002", "01000003", "01000004", "01000005", "01000006", "01000007", "02000000", "02000001", "02000002", "02000003", "02000004", "02000005", "02000006", "02000007", "03000000", "03000001", "03000002", "03000003", "03000004", "03000005", "03000006", "03000007", "10000000", "10000001", "10000003", "10000004", "10000005", "10000006", "10000007", "20000000", "20000001", "20000002", "20000003", "20000004", "20000005", "20000006", "20000007", "30000000", "30000001", "30000002", "30000003", "30000004", "30000005", "30000006", "30000007"}

func Generate_Data() {

	credentials := User()
	all := make(map[interface{}]interface{})
	//all := make(map[string]*zebra.BaseResource)
	for each := 0; each < 100; each++ {

		theType := RandType(resourceTypes)
		theLabels := CreateLabels()
		theData := CreateResource(theType, theLabels)
		fmt.Println("The username, password, corresponding resource and its respective labels: ", theData)
		all[credentials] = theData
	}
}
