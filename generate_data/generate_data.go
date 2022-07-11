package generate_data

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/project-safari/zebra"
	"github.com/project-safari/zebra/auth"
)

var resourceTypes = []string{"VLANPool", "Switch", "IPAddressPool", "Datacenter", "Lab", "Rack", "Server", "ESX", "VM", " ", "BaseResource", "NamedResource", "Credentials"}

// some usernames
func User() string {

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

	username := RandData(name_list)
	return username
}

//some passwords
func Password() string {
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
	pwd := RandData(password_list)

	return pwd
}

//random selection from lists
func RandData(res []string) string {
	rand.Seed(time.Now().UnixNano())
	var length = len(res)
	var ind int = rand.Intn(length - 1)
	typ := res[ind]
	return typ
}

//sample labels
func CreateLabels() map[string]string {
	codes := make(map[string]string)
	col := " "
	let := " "
	colors := []string{"red", "yellow", "green", "blue", "white", "magenta", "black", "purple", "brown", "orange", "pink", "grey"}
	letters := []string{"alpha", "beta", "gamma", "delta", "epsilon", "eta", "theta", "Iota", "Kappa", "Lambda", "Mu", "Nu"}

	for i := 0; i < len(colors); i++ {
		col = RandData(colors)
		let = RandData(letters)
		codes[let] = col
	}

	return codes
}

//put it all together
func Generate_Data() {

	creds := auth.User
	for each := 0; each < 100; each++ {

		theType := RandData(resourceTypes)
		theLabels := CreateLabels()

		//theData := zebra.CreateBaseResource(theType, theLabels)

		creds.zebra.NamedResource = zebra.CreateBaseResource(theType, theLabels)

		creds.PasswordHash = Password()

		creds.Role = User()

		creds.Key = CreateLabels()

		fmt.Println("The updated data with username, password, corresponding resource and its respective labels: ", creds)

	}
}
