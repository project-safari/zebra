/*
create 100 instances of each resource for some users
program to display results for each respective type
*/

package generate_data //nolint // just don't lint the package name.

import (
	"math/rand"
	"net"
	"time"

	"github.com/project-safari/zebra"
	"github.com/project-safari/zebra/auth"
	"github.com/project-safari/zebra/compute"
	"github.com/project-safari/zebra/dc"
	"github.com/project-safari/zebra/network"
)

// this array wil help set the type for each iteration.
func AllResourceTypes() []string {
	resourceTypes := []string{
		"VLANPool", "Switch", "IPAddressPool", "Datacenter", "Lab",
		"Rack", "Server", "ESX", "VM", "VCenter", " ",
	}

	return resourceTypes
}

// this array will help set possible sample IP addresses for each iteration.
func IPsamples() []string {
	SampleIPAddr := []string{
		"192.332.11.05", "192.232.11.37", "192.232.22.05", "192.225.11.05",
		"192.0.0.0", "192.192.192.192", "225.225.225.225", "192.192.64.08",
	}

	return SampleIPAddr
}

// create user roles.
func User() string {
	nameList := []string{
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
	} // some usernames.

	username := RandData(nameList)

	return username
}

// create names.
func Name() string {
	nameList := []string{
		"Marie", "Jack", "Clare",
		"James", "Erika", "Frank",
		"Donna", "John", "Jane",
		"Louis", "Eliza", "Phelippe",
	}

	theName := RandData(nameList)

	return theName
}

// create some passwords.
func Password() string {
	Plist := []string{
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
	pwd := RandData(Plist)

	return pwd
}

// random selection from lists.
func RandData(res []string) string {
	length := len(res)

	rand.Seed(time.Now().UnixNano())

	var ind int = rand.Intn(length - 1) //nolint // random selection for data sampling
	typ := res[ind]

	return typ
}

// creating random starting and end points.
func Range() uint16 {
	nums := []uint16{0, 11, 121, 44, 32, 234, 9, 64, 33, 2, 5, 8, 16, 200, 14, 77}

	rand.Seed(time.Now().UnixNano())

	length := len(nums)

	var ind int = rand.Intn(length - 1) //nolint // sample range numbers for resources

	num := nums[ind]

	return num
}

// creating random ports.
func Ports() uint32 {
	nums := []uint32{1, 2, 3, 4, 6, 8, 9, 16, 27, 32, 36, 54, 72, 64, 81, 128, 162, 216, 256, 512}

	rand.Seed(time.Now().UnixNano())

	length := len(nums)

	var ind int = rand.Intn(length - 1) //nolint // sample port numbers

	port := nums[ind]

	return port
}

// create random sample model names.
func Models() string {
	models := []string{
		"model A", "modelB", "modelC", "modelD",
		"modelA.1", "modelA.1.1", "modelE", "modelD.010",
		"modelF", "modelG", "modelG.1", "modelG.1.1",
	}

	model := RandData(models)

	return model
}

// create random serial codes.
func Serials() string {
	nums := []string{
		"00000", "00001", "00002", "00003",
		"00004", "00005", "00006", "00007",
		"00008", "00009", "00010", "00020",
		"00030", "00040", "00050", "00060",
		"00070", "00080", "00090", "00100",
		"00200", "00300", "00400", "00500",
		"01000", "02000", "03000", "04000",
	}

	ser := RandData(nums)

	return ser
}

// some random rows.
func Rows() string {
	theRows := []string{
		"00A", "00B", "00C", "00D", "00E",
		"0A0", "0B0", "0C0", "0D0", "0E0",
		"A00", "B00", "C00", "D00", "E00",
	}

	oneRow := RandData(theRows)

	return oneRow
}

// create IP arr.
func CreateIPArr(ipNum int) []net.IPNet {
	nets := net.IPNet{} //nolint // just for sample data

	netArr := []net.IPNet{}
	SampleIPAddr := IPsamples()

	for i := 0; i < ipNum; i++ {
		ip := RandData(SampleIPAddr)
		nets.IP = net.IP(ip)
		netArr = append(netArr, nets)
	}

	return netArr
}

// some random address.
func Addresses() string {
	possibleAdr := []string{
		"NYC", "Dallas", "Seattle", "Ottawa", "Paris",
		"London", "Athens", "Milan", "Philadelphia", "Ann Arbor",
		"DC", "Ankara", "Cape Verde", "LA", "Perth",
	}

	theAdr := RandData(possibleAdr)

	return theAdr
}

func Order(start uint16, end uint16) (uint16, uint16) {
	if start > end {
		end, start = start, end
	}

	return start, end
}

// creating new instances of various resource types.
func NewVlanPool(theType string) *network.VLANPool {
	start := Range() // this is the range's start point

	end := Range() // this is the range's end point

	theLabels := CreateLabels() // these are the labels to be used here, in this iteration

	theRes := zebra.NewBaseResource(theType, theLabels) // this is the result of a new base resource creation

	start, end = Order(start, end)

	ret := &network.VLANPool{
		BaseResource: *theRes,
		RangeStart:   start,
		RangeEnd:     end,
	}

	return ret
}

func NewSwitch(theType string, ip net.IP) *network.Switch {
	serial := Serials()

	model := Models()

	ports := Ports()

	theLabels := CreateLabels()

	theRes := zebra.NewBaseResource(theType, theLabels)

	cred := new(zebra.Credentials)

	// add some info to these sample credentials
	named := new(zebra.NamedResource)

	named.BaseResource = *theRes

	named.Name = Name()

	cred.NamedResource = *named

	cred.Keys = CreateLabels()

	ret := &network.Switch{
		BaseResource: *theRes,
		ManagementIP: ip,
		SerialNumber: serial,
		Model:        model,
		NumPorts:     ports,
		Credentials:  *cred,
	}

	return ret
}

func NewIPAddressPool(theType string, netArr []net.IPNet) *network.IPAddressPool {
	theLabels := CreateLabels()

	theRes := zebra.NewBaseResource(theType, theLabels)

	ret := &network.IPAddressPool{
		BaseResource: *theRes,
		Subnets:      netArr,
	}

	return ret
}

func NewDatacenter(theType string) *dc.Datacenter {
	named := new(zebra.NamedResource)

	named.BaseResource = *zebra.NewBaseResource(theType, CreateLabels())

	named.Name = Name()

	ret := &dc.Datacenter{
		NamedResource: *named,
		// some addr.
		Address: Addresses(),
	}

	return ret
}

func NewVCenter(theType string, ip net.IP) *compute.VCenter {
	namedRes := new(zebra.NamedResource)

	namedRes.BaseResource = *zebra.NewBaseResource(theType, CreateLabels())

	namedRes.Name = Name()

	cred := new(zebra.Credentials)

	namedRes.Name = Name()

	cred.NamedResource = *namedRes

	cred.Keys = CreateLabels()

	ret := &compute.VCenter{
		NamedResource: *namedRes,
		IP:            ip,
		Credentials:   *cred,
	}

	return ret
}

func NewLab(theType string) *dc.Lab {
	namedR := new(zebra.NamedResource)

	theLabels := CreateLabels()

	namedR.BaseResource = *zebra.NewBaseResource(theType, theLabels)

	namedR.Name = Name()

	ret := &dc.Lab{
		NamedResource: *namedR,
	}

	return ret
}

func NewRack(theType string) *dc.Rack {
	theLabels := CreateLabels()

	namedRes := new(zebra.NamedResource)

	namedRes.BaseResource = *zebra.NewBaseResource(theType, theLabels)

	namedRes.Name = Name()

	ret := &dc.Rack{
		NamedResource: *namedRes,
		// some row.
		Row: Rows(),
	}

	return ret
}

// sample labels.
func CreateLabels() map[string]string {
	codes := make(map[string]string)

	many := rand.Int() //nolint

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

// put it all together.
func IsGood(manyRes int) bool {
	err := false
	resourceTypes := AllResourceTypes()

	if len(resourceTypes) == 0 || manyRes == 0 {
		err = true
	}

	return err
}

func GenerateData(isGood bool, manyRes int) (*auth.User, []zebra.Resource) {
	ipNum := 10 // number of ip's to have in the []net.IPNet array.

	resourceTypes := AllResourceTypes()
	SampleIPAddr := IPsamples()

	creds := new(auth.User)

	allResources := make([]zebra.Resource, manyRes)

	// go through each resource type.
	for i := 0; i < len(resourceTypes); i++ {
		theType := resourceTypes[i]

		// 100 resources of each type.
		for each := 0; each < manyRes; each++ {
			creds = new(auth.User)
			rsa := new(auth.RsaIdentity)

			sampleIP := RandData(SampleIPAddr)
			IPArr := CreateIPArr(ipNum)

			creds.PasswordHash = Password()
			creds.Role = new(auth.Role)
			creds.Role.Name = User()
			creds.Key = rsa

			switch theType {
			case "VLANPool":
				Res := NewVlanPool(theType)
				allResources[each] = Res

			case "Switch":
				Res := NewSwitch(theType, net.IP(sampleIP))
				allResources[each] = Res

			case "IPAddressPool":
				Res := NewIPAddressPool(theType, IPArr)
				allResources[each] = Res

			case "Datacenter":
				Res := NewDatacenter(theType)
				allResources[each] = Res

			case "VCenter":
				Res := NewVCenter(theType, net.IP(sampleIP))
				allResources[each] = Res

			case "Lab":
				Res := NewLab(theType)
				allResources[each] = Res

			case "Rack":
				Res := NewRack(theType)
				allResources[each] = Res
			}
		}
	}

	return creds, allResources
}
