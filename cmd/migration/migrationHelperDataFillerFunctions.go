//nolint:gomnd, goconst, funlen, unconvert, tagliatelle // Using json that is appropriate to the server.
package main

import (
	"net"
	"strings"
	"time"

	"github.com/project-safari/zebra/model/compute"
	"github.com/project-safari/zebra/model/network"
)

func checkIPAddress(ip string) error {
	if ip == "" || net.ParseIP(ip) == nil || ip == "<nil>" || ip == "N/A" {
		return network.ErrIPEmpty
	}

	return nil
}

func esxFiller(rt Racktables) *ESXData {
	user := ""
	thePassword := strings.ToUpper(user) + "123" + "%"

	theIP := net.IP("1.1.1.1")

	if checkIPAddress(rt.IP) == nil {
		resIP := net.ParseIP(rt.IP)

		if resIP != nil {
			theIP = resIP
		} else {
			theIP = net.ParseIP("0.0.0.0")
		}
	}

	if rt.Owner == "" {
		user = "admin"
	} else {
		user = rt.Owner
	}

	if len(rt.ID) < minIDlength {
		rt.ID = addChars("id", rt.ID)
	}

	if len(thePassword) < 12 {
		thePassword = addChars("password", thePassword)
	}

	loginSession := user + "-" + string(time.Now().Format("2006-01-02T15:04:05-07:00"))

	theSSHkey := "secret-ssh-for" + user

	myType := &TheType{
		Name:        "compute.esx",
		Description: "data center esx",
	}

	myLabels := &TheLabels{
		SystemGroup: "Compute",
	}

	myMeta := &MetaData{
		ID: rt.ID, TheType: *myType,
		CreationTime:       time.Now().Format("2006-01-02T15:04:05-07:00"),
		ModifificationTime: time.Now().Format("2006-01-02T15:04:05-07:00"),
		Owner:              "admin", TheLabels: *myLabels,
		Name: rt.Name,
	}

	myKeys := &TheKeys{
		Password: thePassword,
		SSHKey:   theSSHkey,
	}

	myCreds := &TheCredentials{TheKeys: *myKeys, LoginID: loginSession}

	myESX := &ComputedESX{
		MetaData:       *myMeta,
		TheCredentials: *myCreds, TheServerID: "ESX#" + rt.ID,
		TheIP: theIP.String(),
	}

	myESXArr := make([]ComputedESX, 0)
	myESXArr = append(myESXArr, *myESX)

	myData := &ESXData{TheESX: myESXArr}

	return myData
}

type ServerData struct {
	TheServer []ComputedServer `json:"compute.server"`
}

type TheType struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

type MetaData struct {
	ID                 string `json:"id"`
	TheType            `json:"type"`
	CreationTime       string `json:"creationTime"`
	ModifificationTime string `json:"modifificationTime"`
	Owner              string `json:"owner"`
	TheLabels          `json:"labels"`
	Name               string `json:"name"`
}

type TheLabels struct {
	SystemGroup string `json:"system.group"`
}
type TheStatus struct {
	State string `json:"state"`
}

type TheKeys struct {
	Password string `json:"password"`
	SSHKey   string `json:"ssh-key"`
}

type TheCredentials struct {
	TheKeys `json:"keys"`
	LoginID string `json:"loginId"`
}

type ComputedServer struct {
	MetaData       `json:"meta"`
	TheStatus      `json:"status"`
	TheCredentials `json:"credentials"`
	SerialNumber   string `json:"serialNumber"`
	BoardIP        string `json:"boardIP"`
	Model          string `json:"model"`
}

type ComputedESX struct {
	MetaData       `json:"meta"`
	TheCredentials `json:"credentials"`
	TheServerID    string `json:"serverId"`
	TheIP          string `json:"ip"`
}

type ESXData struct {
	TheESX []ComputedESX `json:"compute.esx"`
}

type ComputedVC struct {
	MetaData       `json:"meta"`
	TheCredentials `json:"credentials"`
	TheIP          string `json:"ip"`
}

type VCData struct {
	TheVC []ComputedVC `json:"compute.vcenter"`
}

type ComputedVM struct {
	MetaData        `json:"meta"`
	TheCredentials  `json:"credentials"`
	TheESXID        string `json:"esxId"`
	TheManagementIP string `json:"managementIp"`
	TheVCenterID    string `json:"vCenterId"`
}

type VMData struct {
	TheVM []ComputedVM `json:"compute.vm"`
}

type ComputedDC struct {
	MetaData `json:"meta"`
	Address  string `json:"address"`
}

type DCData struct {
	TheDC []ComputedDC `json:"dc.datacenter"`
}

type ComputedLab struct {
	MetaData       `json:"meta"`
	TheStatus      `json:"status"`
	TheCredentials `json:"credentials"`
}

type LabData struct {
	TheLab []ComputedLab `json:"dc.lab"`
}

type ComputedRack struct {
	MetaData `json:"meta"`
	Row      string `json:"row"`
}

type RackData struct {
	TheRack []ComputedRack `json:"dc.rack"`
}

type ComputedSwitch struct {
	MetaData        `json:"meta"`
	TheCredentials  `json:"credentials"`
	TheManagementIP string `json:"managementIp"`
	TheSerialNumber string `json:"serialNumber"`
	TheModel        string `json:"model"`
	TheNumPorts     uint32 `json:"numPorts"`
}

type SwitchData struct {
	TheSwitch []ComputedSwitch `json:"network.switch"`
}

type ComputedVLAN struct {
	MetaData      `json:"meta"`
	TheRangeStart uint16 `json:"rangeStart"`
	TheRangeEnd   uint16 `json:"rangeEnd"`
}

type VLANData struct {
	TheVLAN []ComputedVLAN `json:"network.vlanPool"`
}

type ComputedAddressPool struct {
	MetaData   `json:"meta"`
	TheSubnets []net.IPNet `json:"subnets"`
}

type AddressPoolData struct {
	TheAddressPool []ComputedAddressPool `json:"network.ipAddressPool"`
}

func vcenterFiller(rt Racktables) *VCData {
	if rt.Owner == "" {
		rt.Owner = "admin"
	}

	theIP := net.IP("1.1.1.1")

	user := rt.Owner

	thePassword := strings.ToUpper(user) + "123" + "!"

	if checkIPAddress(rt.IP) == nil {
		resIP := net.ParseIP(rt.IP)

		if resIP != nil {
			theIP = resIP
		} else {
			theIP = net.ParseIP("0.0.0.0")
		}
	}

	if len(rt.ID) < minIDlength {
		rt.ID = addChars("id", rt.ID)
	}

	if len(thePassword) < 12 {
		thePassword = addChars("password", thePassword)
	}

	loginSession := user + "-" + string(time.Now().Format("2006-01-02T15:04:05-07:00"))
	theSSHkey := "secret-ssh-for" + user

	myType := &TheType{
		Name:        "compute.vcenter",
		Description: "VMWare vcenter",
	}

	myLabels := &TheLabels{
		SystemGroup: "Compute",
	}

	myMeta := &MetaData{
		ID:                 rt.ID,
		TheType:            *myType,
		CreationTime:       time.Now().Format("2006-01-02T15:04:05-07:00"),
		ModifificationTime: time.Now().Format("2006-01-02T15:04:05-07:00"),
		Owner:              "admin", TheLabels: *myLabels,
		Name: rt.Name,
	}

	myKeys := &TheKeys{
		Password: thePassword,
		SSHKey:   theSSHkey,
	}

	myCreds := &TheCredentials{TheKeys: *myKeys, LoginID: loginSession}

	myVC := &ComputedVC{
		MetaData:       *myMeta,
		TheCredentials: *myCreds,
		TheIP:          theIP.String(),
	}

	myVCArr := make([]ComputedVC, 0)
	myVCArr = append(myVCArr, *myVC)

	myData := &VCData{TheVC: myVCArr}

	return myData
}

func serverFiller(rt Racktables) *ServerData {
	theIP := net.IP("1.1.1.1")

	user := ""

	thePassword := strings.ToUpper(user) + "123" + "%"

	if len(rt.ID) < minIDlength {
		rt.ID = addChars("id", rt.ID)
	}

	if len(thePassword) < 12 {
		thePassword = addChars("password", thePassword)
	}

	if rt.Owner == "" {
		user = "admin"
	} else {
		user = rt.Owner
	}

	if checkIPAddress(rt.IP) == nil {
		resIP := net.ParseIP(rt.IP)

		if resIP != nil {
			theIP = resIP
		} else {
			theIP = net.ParseIP("0.0.0.0") // some default IP since db is empty.
		}
	}

	serverModel := "db-compute-server" + rt.Name + rt.ID

	theSSHkey := "secret-ssh-for" + user

	myType := &TheType{
		Name:        "compute.server",
		Description: "data center server",
	}

	myLabels := &TheLabels{
		SystemGroup: "Server",
	}

	myStatus := &TheStatus{State: "active"}

	myMeta := &MetaData{
		ID: "0198300027", TheType: *myType,
		CreationTime:       time.Now().Format("2006-01-02T15:04:05-07:00"),
		ModifificationTime: time.Now().Format("2006-01-02T15:04:05-07:00"),
		Owner:              "admin", TheLabels: *myLabels,
		Name: rt.Name,
	}

	loginSession := user + "-" + string(time.Now().Format("2006-01-02T15:04:05-07:00"))

	myKeys := &TheKeys{
		Password: thePassword,
		SSHKey:   theSSHkey,
	}

	myCreds := &TheCredentials{TheKeys: *myKeys, LoginID: loginSession}

	myServer := &ComputedServer{
		MetaData: *myMeta, TheStatus: *myStatus,
		TheCredentials: *myCreds, SerialNumber: rt.Serial,
		BoardIP: theIP.String(), Model: serverModel,
	}

	myServerArr := make([]ComputedServer, 0)
	myServerArr = append(myServerArr, *myServer)

	myData := &ServerData{TheServer: myServerArr}

	return myData
}

func vmFiller(rt Racktables) *VMData {
	res := compute.NewVM("DB-esx", rt.Name, rt.Owner, "system.group-server-vcenter-vm")

	if rt.Owner == "" {
		rt.Owner = "admin"
	}

	theIP := "1.1.1.1"
	user := rt.Owner

	if checkIPAddress(rt.IP) == nil {
		resIP := net.ParseIP(rt.IP)

		if resIP != nil {
			res.ManagementIP = resIP
			theIP = resIP.String()
		} else {
			res.ManagementIP = net.ParseIP("0.0.0.0")
			theIP = res.ManagementIP.String()
		}
	}

	thePassword := strings.ToUpper(user) + "123" + "!"

	if len(rt.ID) < minIDlength {
		rt.ID = addChars("id", rt.ID)
	}

	if len(thePassword) < 12 {
		thePassword = addChars("password", thePassword)
	}

	myType := &TheType{
		Name:        "compute.vm",
		Description: "virtual machine",
	}

	myLabels := &TheLabels{
		SystemGroup: "Compute",
	}

	myMeta := &MetaData{
		ID: "0198300027", TheType: *myType,
		CreationTime:       time.Now().Format("2006-01-02T15:04:05-07:00"),
		ModifificationTime: time.Now().Format("2006-01-02T15:04:05-07:00"),
		Owner:              user, TheLabels: *myLabels,
		Name: rt.Name,
	}

	theSSHkey := "secret-ssh-for" + user

	loginSession := user + "-" + string(time.Now().Format("2006-01-02T15:04:05-07:00"))

	myKeys := &TheKeys{
		Password: thePassword,
		SSHKey:   theSSHkey,
	}

	myCreds := &TheCredentials{TheKeys: *myKeys, LoginID: loginSession}

	myVM := &ComputedVM{
		MetaData:        *myMeta,
		TheCredentials:  *myCreds,
		TheESXID:        "ESX-" + rt.Name + rt.ID,
		TheManagementIP: theIP,
		TheVCenterID:    "VM-" + rt.Name + rt.ID,
	}

	myVMArr := make([]ComputedVM, 0)
	myVMArr = append(myVMArr, *myVM)

	myData := &VMData{TheVM: myVMArr}

	return myData
}

func labFiller(rt Racktables) *LabData {
	var user string

	if rt.Owner != "" {
		user = rt.Owner
	} else {
		user = "admin"
	}

	if len(rt.ID) < minIDlength {
		rt.ID = addChars("id", rt.ID)
	}

	myType := &TheType{
		Name:        "dc.lab",
		Description: "data center lab",
	}

	myLabels := &TheLabels{
		SystemGroup: "Lab",
	}

	myMeta := &MetaData{
		ID:                 rt.ID,
		TheType:            *myType,
		CreationTime:       time.Now().Format("2006-01-02T15:04:05-07:00"),
		ModifificationTime: time.Now().Format("2006-01-02T15:04:05-07:00"),
		Owner:              rt.Owner,
		TheLabels:          *myLabels,
		Name:               rt.Name,
	}

	thisPassword := strings.ToUpper(user) + "123" + "!"
	if len(thisPassword) < 12 {
		thisPassword = addChars("password", thisPassword)
	}

	theSSHkey := "secret-ssh-for" + user

	loginSession := user + "-" + string(time.Now().Format("2006-01-02T15:04:05-07:00"))

	myKeys := &TheKeys{
		Password: thisPassword,
		SSHKey:   theSSHkey,
	}

	myCreds := &TheCredentials{TheKeys: *myKeys, LoginID: loginSession}

	myStatus := &TheStatus{State: "active"}

	myLab := &ComputedLab{
		MetaData:       *myMeta,
		TheStatus:      *myStatus,
		TheCredentials: *myCreds,
	}

	myLabArr := make([]ComputedLab, 0)
	myLabArr = append(myLabArr, *myLab)

	myData := &LabData{TheLab: myLabArr}

	return myData
}

func dcFiller(rt Racktables) *DCData {
	theID := rt.ID

	myType := &TheType{
		Name:        "dc.datacenter",
		Description: "data center",
	}

	myLabels := &TheLabels{
		SystemGroup: "Datacenter",
	}

	if len(theID) < 7 {
		theID = addChars("id", theID)
	}

	myDatacenterArr := make([]ComputedDC, 0)

	myMeta := &MetaData{
		ID:                 theID,
		TheType:            *myType,
		CreationTime:       time.Now().Format("2006-01-02T15:04:05-07:00"),
		ModifificationTime: time.Now().Format("2006-01-02T15:04:05-07:00"),
		Owner:              rt.Owner,
		TheLabels:          *myLabels,
		Name:               rt.Name,
	}

	myDC := &ComputedDC{
		MetaData: *myMeta,
		Address:  rt.Location,
	}

	myDatacenterArr = append(myDatacenterArr, *myDC)

	myData := &DCData{TheDC: myDatacenterArr}

	return myData
}

func rackFiller(rt Racktables) *RackData {
	myType := &TheType{
		Name:        "dc.rack",
		Description: "data center rack",
	}

	theID := rt.ID
	if len(theID) < 7 {
		theID = addChars("id", theID)
	}

	myLabels := &TheLabels{
		SystemGroup: "Rack",
	}

	myMeta := &MetaData{
		ID:                 theID,
		TheType:            *myType,
		CreationTime:       time.Now().Format("2006-01-02T15:04:05-07:00"),
		ModifificationTime: time.Now().Format("2006-01-02T15:04:05-07:00"),
		Owner:              rt.Owner,
		TheLabels:          *myLabels,
		Name:               rt.Name,
	}

	myRack := &ComputedRack{
		MetaData: *myMeta,
		Row:      rt.RowName,
	}

	myRackArr := make([]ComputedRack, 0)
	myRackArr = append(myRackArr, *myRack)

	myData := &RackData{TheRack: myRackArr}

	return myData
}

func switchFiller(rt Racktables) *SwitchData {
	myType := &TheType{
		Name:        "network.switch",
		Description: "network server",
	}

	user := rt.Owner

	if user == "" { // if user is unspecified, the default user is admin.
		user = "admin"
	}

	if len(rt.ID) < minIDlength {
		rt.ID = addChars("id", rt.ID)
	}

	myPassword := strings.ToUpper(user) + "123" + "%"

	if len(myPassword) < 12 {
		myPassword = addChars("password", myPassword)
	}

	myLabels := &TheLabels{
		SystemGroup: "Switch",
	}

	myMeta := &MetaData{
		ID:                 rt.ID,
		TheType:            *myType,
		CreationTime:       time.Now().Format("2006-01-02T15:04:05-07:00"),
		ModifificationTime: time.Now().Format("2006-01-02T15:04:05-07:00"),
		Owner:              rt.Owner,
		TheLabels:          *myLabels,
		Name:               rt.Name,
	}

	theSSHkey := "secret-ssh-for" + user

	loginSession := user + "-" + string(time.Now().Format("2006-01-02T15:04:05-07:00"))

	myKeys := &TheKeys{Password: myPassword, SSHKey: theSSHkey}

	myCreds := &TheCredentials{TheKeys: *myKeys, LoginID: loginSession}

	theIP := "0.0.0.0" // a sample IP address - could possibly add functionality where if IP is missing, more is done.

	if checkIPAddress(rt.IP) == nil {
		resIP := net.ParseIP(theIP)

		if resIP != nil && rt.IP != "" {
			theIP = net.ParseIP(rt.IP).String()
		} else {
			theIP = "0.0.0.0"
		}
	}

	mySwitch := &ComputedSwitch{
		MetaData:        *myMeta,
		TheCredentials:  *myCreds,
		TheManagementIP: theIP,
		TheSerialNumber: rt.Serial,
		TheModel:        "model",
		TheNumPorts:     uint32(rt.Port),
	}

	mySwitchArr := make([]ComputedSwitch, 0)
	mySwitchArr = append(mySwitchArr, *mySwitch)

	myData := &SwitchData{TheSwitch: mySwitchArr}

	return myData
}

func addressPoolFiller(rt Racktables) *AddressPoolData {
	var resIP net.IP

	if checkIPAddress(rt.IP) == nil {
		resIP = net.ParseIP(rt.IP)

		if resIP == nil && net.IPMask(resIP) == nil {
			resIP = net.ParseIP("0.0.0.0")
		}
	}

	if rt.Owner == "" {
		rt.Owner = "admin"
	}

	user := rt.Owner

	theSubnets := []net.IPNet{{
		IP:   resIP,
		Mask: net.IPMask(resIP),
	}}

	myType := &TheType{
		Name:        "network.ipAddressPool",
		Description: "a network ip address pool",
	}

	myLabels := &TheLabels{
		SystemGroup: "ipAddrressPool",
	}

	if len(rt.ID) < minIDlength {
		rt.ID = addChars("id", rt.ID)
	}

	myMeta := &MetaData{
		ID:                 rt.ID,
		TheType:            *myType,
		CreationTime:       time.Now().Format("2006-01-02T15:04:05-07:00"),
		ModifificationTime: time.Now().Format("2006-01-02T15:04:05-07:00"),
		Owner:              user,
		TheLabels:          *myLabels,
		Name:               rt.Name,
	}

	myAddressPool := &ComputedAddressPool{
		MetaData:   *myMeta,
		TheSubnets: theSubnets,
	}

	myAddrArr := make([]ComputedAddressPool, 0)
	myAddrArr = append(myAddrArr, *myAddressPool)

	myData := &AddressPoolData{TheAddressPool: myAddrArr}

	return myData
}

func vlanFiller(rt Racktables) *VLANData {
	user := "admin"

	if rt.Owner != "" {
		user = rt.Owner
	}

	rangeStart := 0 // curently no provision for this in the db.
	someEnd := 100  // curently no provision for this in the db.
	rangeEnd := uint16(someEnd)

	myType := &TheType{
		Name:        "network.vlanPool",
		Description: "network vlan pool",
	}

	myLabels := &TheLabels{
		SystemGroup: "Vlan",
	}

	if len(rt.ID) < minIDlength {
		rt.ID = addChars("id", rt.ID)
	}

	myMeta := &MetaData{
		ID:                 rt.ID,
		TheType:            *myType,
		CreationTime:       time.Now().Format("2006-01-02T15:04:05-07:00"),
		ModifificationTime: time.Now().Format("2006-01-02T15:04:05-07:00"),
		Owner:              user,
		TheLabels:          *myLabels,
		Name:               rt.Name,
	}

	myVLAN := &ComputedVLAN{
		MetaData:      *myMeta,
		TheRangeStart: uint16(rangeStart),
		TheRangeEnd:   rangeEnd,
	}

	myVLANArr := make([]ComputedVLAN, 0)
	myVLANArr = append(myVLANArr, *myVLAN)

	myData := &VLANData{TheVLAN: myVLANArr}

	return myData
}
