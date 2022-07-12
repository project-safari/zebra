package generate_data
/*
generate sample data and users
*/

import (
	"fmt"
	"math/rand"
	"net"
	"time"

	"github.com/google/uuid"
	"github.com/project-safari/zebra"
	"github.com/project-safari/zebra/auth"
	"github.com/project-safari/zebra/dc"
	"github.com/project-safari/zebra/network"
)

var resourceTypes = []string{"VLANPool", "Switch", "IPAddressPool", "Datacenter", "Lab", "Rack", "Server", "ESX", "VM", " "}

// create user roles
func User() string {

	// some usernames
	name_list := []string{
		"user1",
		"user2",
		"user3",
		"manager",
		"admin",
		"guest",
		"another-user",
		" ",
		"engineer1",
		"designer",
		"private_user",
		"ceo",
		"staff",
	}

	username := RandData(name_list)

	return username

}

//create some passwords
func Password() string {
	
	password_list := []string{
		"pass123",
		"pass222",
		"pass321",
		"m3g2r1",
		"passthrough",
		" ",
		"another",
		"no_pass",
		"somepass",
		"v)entry2",
		"p%3eqr3",
		"random_pass",
		"validate000",
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

// creating random starting and end points
func Range() uint16 {

	nums := []uint16{0, 11, 121, 44, 32, 234, 9, 64, 33, 2, 5, 8, 16, 200, 14, 77}
	rand.Seed(time.Now().UnixNano())
	var length = len(nums)
	var ind int = rand.Intn(length - 1)
	num := nums[ind]
	
	return num
}

// creating random ports
func Ports() uint32 {
	
	nums := []uint32{1, 2, 3, 4, 6, 8, 9, 16, 27, 32, 36, 54, 72, 64, 81, 128, 162, 216, 256, 512}
	rand.Seed(time.Now().UnixNano())
	var length = len(nums)
	var ind int = rand.Intn(length - 1)
	port := nums[ind]
	
	return port
	
}

// create random sample model names
func Models() string {
	
	models := []string{"model A", "modelB", "modelC", "modelD", "modelA.1", "modelA.1.1", "modelE", "modelD.010", "modelF", "modelG", "modelG.1", "modelG.1.1"}
	model := RandData(models)
	
	return model
	
}

// create random serial codes
func Serials() string {
	
	nums := []string{"00000", "00001", "00002", "00003", "00004", "00005", "00006", "00007", "00008", "00009", "00010", "00020", "00030", "00040", "00050", "00060", "00070", "00080", "00090", "00100", "00200", "00300", "00400", "00500", "01000", "02000", "03000"}
	return RandData(nums)
	
}

//creating various resource types
func CreateResource(resType string, labels zebra.Labels, start uint16, theRes zebra.BaseResource, ip net.IP, serial string, model string, ports uint32, namedRes zebra.NamedResource) zebra.Resource {

	id := uuid.New().String()

	// base case (i.e. nothing exists for resource type)
	if resType == " " {
		resType = "BaseResource"
	}

	if resType == "VLANPool" {
		return &network.VLANPool{

			BaseResource: theRes,
			RangeStart:   start,
			RangeEnd:     (start + 100),
		}
	} else if resType == "Switch" {
		return &network.Switch{

			BaseResource: theRes,
			ManagementIP: ip,
			SerialNumber: serial,
			Model:        model,
			NumPorts:     ports,
		}
	} else if resType == "IPAddressPool" {
		return &network.IPAddressPool{
			BaseResource: theRes,
			//Subnets:      []net.IPNet,
		}

	} else if resType == "Datacenter" {
		return &dc.Datacenter{
			NamedResource: namedRes,
			Address:       "sample address",
		}
	} else if resType == "Lab" {
		return &dc.Lab{
			NamedResource: namedRes,
		}
	} else if resType == "Rack" {
		return &dc.Rack{
			NamedResource: namedRes,
			Row:           "sample row",
		}
	}
	// if none of those apply, just use base resource info
	return &zebra.BaseResource{
		ID:     id,
		Type:   resType,
		Labels: labels,
	}

}

//sample labels
func CreateLabels() (string, map[string]string) {
	
	codes := make(map[string]string)
	col := " "
	let := " "
	colors := []string{"red", "yellow", "green", "blue", "white", "magenta", "black", "purple", "brown", "orange", "pink", "grey"}
	letters := []string{"alpha", "beta", "gamma", "delta", "epsilon", "eta", "theta", "Iota", "Kappa", "Lambda", "Mu", "Nu"}

	col = RandData(colors)
	let = RandData(letters)
	codes[let] = col

	// get the key and the label pair

	return codes[let], codes
	
}


//put it all together
func Generate_Data() {

	// 100 resources
	for each := 0; each < 100; each++ {
		
		//generate new user
		
		creds := new(auth.User)

		// info to be used in the resources
		
		theType := RandData(resourceTypes)
		keys, theLabels := CreateLabels()
		start := Range()
		base := zebra.NewBaseResource(theType, theLabels)
		serial := Serials()
		port := Ports()
		model := Models()
		sampleIP := "192.232.11.05"
		named := new(zebra.NamedResource)
		
		// update info
		
		creds.zebra.NamedResource = CreateResource(theType, theLabels, start, base, net.IP(sampleIP), serial, model, port, named)
		creds.PasswordHash = Password()
		creds.Role = User()
		creds.Key = keys

		// display info
		fmt.Println("Information: ", each, "\nThe data with username, password, corresponding resource and its respective labels: ", creds)

	}
}
