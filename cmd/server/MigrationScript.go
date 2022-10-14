// Migration Script for data - it can be used to fetch,
// add, and use the data inside zebra.
// This is the Golang version of the same python script.
package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/julienschmidt/httprouter" //nolint:gci
	"github.com/project-safari/zebra"
	"github.com/project-safari/zebra/model/compute"
	"github.com/project-safari/zebra/model/dc"
	"github.com/project-safari/zebra/model/network"

	"github.com/go-logr/logr"
)

// var ErrInvalidImport = errors.New("the import of data from the data base failed")

// Racktables struct is the struct that contains info from the racktables table in the mysql db.
//
// It contains the id, ip, name, and type of the item in the racktable.
type Racktables struct {
	//  id, name, label, objtype_id, asset_no, has_problems, comment
	ID   string `json:"object_id"` //nolint:tagliatelle
	Name string `json:"name"`
	// Type TheType `json:"type"`
	Label     string `json:"label"`
	ObjtypeID string `json:"objtypeId"`
	AssetNo   string `json:"assetNo"`
	Problems  string `json:"hasProblems"`
	Comments  string `json:"comment"`
	IP        string `json:"ip"`
	Type      string `json:"type"`
	Port      int    `json:"port"`
	RackID    string `json:"rackId"`
	RowName   string `json:"rowName"`
	Owner     string `json:"owner"`
	RowID     string `json:"rowId"`
	Location  string `json:"locationName"`
}

// Determine the specific type of a resource.
//nolint
func determineType(means string, resName string) string {
	name := strings.ToLower(resName)
	typ := ""

	if means == "Shelf" {
		typ = "Rack"
	} else if means == "Compute" {
		if strings.Contains(name, "esx") {
			typ = "compute.esx"
		} else if strings.Contains(name, "jenkins") || strings.Contains(name, "server") || strings.Contains(name, "srv") || strings.Contains(name, "vintella") {
			typ = "compute.server"
		} else if strings.Contains(name, "datacenter") || strings.Contains(name, "dc") || strings.Contains(name, "bld") {
			typ = "dc.datacenter"
		} else if strings.Contains(name, "dmz") || strings.Contains(name, "vlan") || strings.Contains(name, "asa") || strings.Contains(name, "bridge") {
			typ = "network.vlanPool"
		} else if strings.Contains(name, "vleaf") || strings.Contains(name, "switch") || strings.Contains(name, "sw") || strings.Contains(name, "aci") {
			typ = "network.switch"
		} else if strings.Contains(name, "vm") || strings.Contains(name, "capic") || strings.Contains(name, "frodo") {
			typ = "compute.vm"
		} else if strings.Contains(name, "vapic") || strings.Contains(name, "vpod") {
			typ = "compute.vcenter"
		} else if strings.Contains(name, "ipc") {
			typ = "network.ipAddressPool"
		}
	} else if means == "Other" {
		if strings.Contains(name, "chasis") || strings.Contains(name, "ixia") || strings.Contains(name, "rack") {
			typ = "dc.rack"
		} else if strings.Contains(name, "nexus") || strings.Contains(name, "sw") || strings.Contains(name, "switch") || strings.Contains(name, "n3k") {
			typ = "network.switch"
		}
	} else {
		typ = means
	}

	return typ
}

/*
   ### still need something for vapic* and vpod*, as well as FRODO*,  APPLIANCE-HOME1, CAPIC*,
       # aci-github.cisco.com*, DMASHAL-VINTELLA*.
       # RESOLVED AS EXPLAINED BELOW:
   ### vpod uses VMware ESXi hosts, VMware vCenter, storage, networking and a Windows Console VM.
           # => vcenter.
   ### vAPIC virtual machines use VMware vCenter ==> vcenter
       # Cisco ACI vCenter plug-in.
       # BUT also uses Cisco ACI Virtual Edge VM.
   ### Cisco Cloud APIC on Microsoft Azure is deployed and runs as an
       # Microsoft Azure Virtual Machine => capic => VM.
   ### Frodo is enabled by default on VMs powered on after AOS 5.5.X => frodo => VM.
       # About frodo - VMware Technology Network VMTN
   ### Vintela -> VAS is Vintela's flagship product in a line that includes Vintela Management
       # eXtensions (VMX), which extends Microsoft Systems Management Server => server.
   ### apic uses controllers and so does cisco aci but it is similar to switches => switch.

   ## vcenter is a management interface type => management interface type = vcenter
*/

// Get resource type by id.
//nolint
func determineIDMeaning(id string, name string) string {
	means := ""
	final := ""
	this := ""

	if id == "2" || id == "27" {
		means = "compute.vm"
	} else if id == "30" || id == "31" || id == "34" {
		means = "dc.rack"
	} else if id == "3" {
		means = "dc.rack"
	} else if id == "38" {
		means = "compute.vcenter"
	} else if id == "4" || id == "13" || id == "36" {
		means = "compute.server"
	} else if id == "8" || id == "12" || id == "14" || id == "21" || id == "26" || id == "32" || id == "33" {
		means = "network.switch"
	} else if id == "1504" {
		means = "Compute"
	} else if id == "1503" {
		means = "Other"
	} else {
		means = "/"
	}

	final = determineType(means, name)

	if final == "/" {
		final = "unclassified"
	}

	this = final

	return this
}

// The getRack function gets data from the data base.
//
// It prepares the necessary data from the data base to be used in the zebra tool.
//
// getRack takes in a Type struct and returns an error or nil in the absence thereof.
//
// It executes the function newMeta with the data it extracted from the db.
//nolint:funlen
func getData() ([]Racktables, error) {
	var rt Racktables

	RackArr := []Racktables{}

	// Statement to query the db - currently only one rack, 76.
	statement := "SELECT rack_id, object_id  FROM rackspace WHERE rack_id = 76"
	// to be filled in with appropriate user, password, and db name.
	db, err := sql.Open("mysql", "eachim:1234@/racktables")
	// if there is an error opening the connection, handle it
	if err != nil {
		log.Print(err.Error())
	}

	defer db.Close()

	// Execute the query
	results, err := db.Query(statement)
	if err != nil {
		panic(err.Error()) // proper error handling instead of panic in your app
	}

	for results.Next() {
		// for each row, scan the result into our tag composite object
		err = results.Scan(&rt.RackID, &rt.ID)

		rt.ID, rt.Name, rt.Label, rt.ObjtypeID, rt.AssetNo, rt.Problems, rt.Comments = getMoreDetails(rt.ID)

		typeID := rt.ObjtypeID
		resType := determineIDMeaning(typeID, rt.Name)
		rt.Type = resType

		if strings.Contains(resType, "compute") || resType == "network.switch" {
			rt.IP = getIPDetaiLs(rt.ID)

			ownedBy := getUserDetails(rt.IP)
			rt.Owner = ownedBy
		} else {
			rt.IP = "null"

			ownedBy := "null"
			rt.Owner = ownedBy
		}

		if resType == "network.switch" {
			portInfo := getPortDetails(rt.ID)
			rt.Port = portInfo
		} else {
			portID := -1
			rt.Port = portID
		}

		rowName, rowID, rowLocation := getRowDetails(rt.ID)
		rt.RowName = rowName
		rt.RowID = rowID
		rt.Location = rowLocation

		assetNumber := rt.AssetNo
		rt.AssetNo = assetNumber
		probs := rt.Problems
		rt.Problems = probs

		notes := rt.Comments
		rt.Comments = notes

		RackArr = append(RackArr, rt)

		if err != nil {
			panic(err.Error())
		}
	}

	if err != nil {
		fmt.Println(err.Error())
	}

	return RackArr, err
}

// Get IPs from db based on type id.
func getIPDetaiLs(objectID string) string {
	var rt Racktables

	statement := "SELECT ip FROM IPv4Allocation WHERE object_id = ?"

	// to be filled in with appropriate user, password, and db name.
	db, err := sql.Open("mysql", "eachim:1234@/racktables")
	// if there is an error opening the connection, handle it
	if err != nil {
		log.Print(err.Error())
	}

	defer db.Close()

	// Execute the query
	results, err := db.Query(statement, objectID)
	if err != nil {
		panic(err.Error()) // proper error handling instead of panic in your app
	}
	defer results.Close()

	for results.Next() {
		// for each row, scan the result into our tag composite object
		err = results.Scan(&rt.IP)

		if err != nil {
			panic(err.Error()) // proper error handling instead of panic in your app
		}
	}

	return rt.IP
}

// Get port IDs from db based on type ID.
func getPortDetails(objectID string) int {
	var rt Racktables

	numPort := 0

	statement := "SELECT id FROM Port WHERE object_id = ?"

	// to be filled in with appropriate user, password, and db name.
	db, err := sql.Open("mysql", "eachim:1234@/racktables")
	// if there is an error opening the connection, handle it
	if err != nil {
		log.Print(err.Error())
	}

	defer db.Close()

	// Execute the query
	results, err := db.Query(statement, objectID)
	if err != nil {
		panic(err.Error()) // proper error handling instead of panic in your app
	}
	defer results.Close()

	for results.Next() {
		// for each row, scan the result into our tag composite object
		err = results.Scan(&rt.Port)

		if err != nil {
			panic(err.Error()) // proper error handling instead of panic in your app
		}

		numPort++
	}

	return numPort
}

// Get rack details using the resource's specific ID.
func getMoreDetails(objectID string) (string, string, string, string, string, string, string) {
	var rt Racktables

	statement := "SELECT id, name, label, objtype_id, asset_no, has_problems, comment FROM rackobject WHERE id = ?"
	// to be filled in with appropriate user, password, and db name.
	db, err := sql.Open("mysql", "eachim:1234@/racktables")
	// if there is an error opening the connection, handle it
	if err != nil {
		log.Print(err.Error())
	}

	defer db.Close()

	// Execute the query
	results, err := db.Query(statement, objectID)
	if err != nil {
		panic(err.Error()) // proper error handling instead of panic in your app
	}
	defer results.Close()

	for results.Next() {
		// for each row, scan the result into our tag composite object, GIVEN THE ID.
		err = results.Scan(&rt.ID, &rt.Name, &rt.Label, &rt.ObjtypeID, &rt.AssetNo, &rt.Problems, &rt.Comments)

		if err != nil {
			panic(err.Error()) // proper error handling instead of panic in your app
		}

		log.Print(rt.RackID)
	}

	return rt.ID, rt.Name, rt.Label, rt.ObjtypeID, rt.AssetNo, rt.Problems, rt.Comments
}

// Get row and location details based on rack info (rack ID).
func getRowDetails(id string) (string, string, string) {
	var rt Racktables

	statement := "SELECT row_id, row_name, location_name FROM Rack WHERE id = ?"

	// to be filled in with appropriate user, password, and db name.
	db, err := sql.Open("mysql", "eachim:1234@/racktables")
	// if there is an error opening the connection, handle it
	if err != nil {
		log.Print(err.Error())
	}

	defer db.Close()

	// Execute the query
	results, err := db.Query(statement, id)
	if err != nil {
		panic(err.Error()) // proper error handling instead of panic in your app
	}
	defer results.Close()

	for results.Next() {
		// for each row, scan the result into our tag composite object
		err = results.Scan(&rt.RowID, &rt.RowName, &rt.Location)

		if err != nil {
			panic(err.Error()) // proper error handling instead of panic in your app
		}
	}

	return rt.RowID, rt.RowName, rt.Location
}

// Get owner / user details based on the resource's IP.
func getUserDetails(resIP string) string {
	var rt Racktables

	statement := "SELECT user FROM IPv4Log WHERE ip = ?"

	// to be filled in with appropriate user, password, and db name.
	db, err := sql.Open("mysql", "eachim:1234@/racktables")
	// if there is an error opening the connection, handle it
	if err != nil {
		log.Print(err.Error())
	}

	defer db.Close()

	// Execute the query
	results, err := db.Query(statement, resIP)
	if err != nil {
		panic(err.Error()) // proper error handling instead of panic in your app
	}
	defer results.Close()

	for results.Next() {
		// for each row, scan the result into our tag composite object
		err = results.Scan(&rt.Owner)

		if err != nil {
			panic(err.Error()) // proper error handling instead of panic in your app
		}
	}

	return rt.Owner
}

// Function to create a resource given data obtained from db, guven a certain type.
//
// Returns a zebra.Resource.
//nolint:cyclop
func createResFromData(res Racktables) zebra.Resource {
	// dbResources, err := GetData()
	// var resType = ""
	resType := res.Type

	switch resType {
	case "dc.datacenetr":
		return dc.NewDatacenter(res.Location, res.Name, res.Owner, "system.group-datacenter")

	case "dc.lab":
		return dc.NewLab(res.Name, res.Owner, "system.group-datacenter-lab")

	case "dc.rack", "dc.shelf":
		return dc.NewRack(res.RowName, res.RowID, res.Name, res.Location, res.Owner, "system.group-datacenter-lab-rack")

	case "compute.server":
		return compute.NewServer("serial", "model", res.Name, res.Owner, "system.group-server")

	case "compute.esx":
		return compute.NewESX(res.ID, res.Name, res.Owner, "system.group-server-esx")

	case "compute.vm":
		return compute.NewVM("esx??", res.Name, res.Owner, "system.group-server-vcenter-vm")

	case "compute.vcenetr":
		return compute.NewVCenter(res.Name, res.Owner, "system.group-server-vcenter")

	case "network.switch":
		return network.NewSwitch(res.Name, res.Owner, "system.group-vlan-switch")

	case "network.ipaddresspool":
		return network.NewIPAddressPool(res.Name, res.Owner, "system.group-vlan-ipaddrpool")

	case "network.vlanpool":
		return network.NewVLANPool(res.Name, res.ObjtypeID, "system.group-vlan")
	}

	return nil
}

/*
// Validate all resources in a resource map.
func validateResources(ctx context.Context, resMap *zebra.ResourceMap) error {
	// Check all resources to make sure they are valid
	for _, l := range resMap.Resources {
		for _, r := range l.Resources {
			if err := r.Validate(ctx); err != nil {
				return err
			}
		}
	}

	return nil
}

func readJSON(ctx context.Context, req *http.Request, data interface{}) error {
	log := logr.FromContextOrDiscard(ctx)

	body, err := ioutil.ReadAll(req.Body)
	if err != nil {
		return err
	}

	log.Info("request", "body", string(body))

	if len(body) > 0 {
		err = json.Unmarshal(body, data)
	} else {
		err = ErrEmptyBody
	}

	return err
}
*/

func postToZebra() httprouter.Handle {
	return func(res http.ResponseWriter, req *http.Request, params httprouter.Params) {
		ctx := req.Context()
		log := logr.FromContextOrDiscard(ctx)
		api, ok := ctx.Value(ResourcesCtxKey).(*ResourceAPI)

		if !ok {
			res.WriteHeader(http.StatusInternalServerError)

			return
		}

		factory := zebra.Factory()
		resMap := zebra.NewResourceMap(factory)

		dbResources, err := getData()
		if err != nil {
			log.Info("the import of data from the data base failed")
		}

		for i := 0; i < len(dbResources); i++ {
			zebraRes := dbResources[i]

			resource := createResFromData(zebraRes)

			resMap.Add(resource) //nolint:errcheck
		}

		// Read request, return error if applicable
		if err := readJSON(ctx, req, resMap); err != nil {
			res.WriteHeader(http.StatusBadRequest)
			log.Info("resources could not be created, could not read request")

			return
		}

		if validateResources(ctx, resMap) != nil {
			res.WriteHeader(http.StatusBadRequest)
			log.Info("resources could not be created, found invalid resource(s)")

			return
		}

		// Add all resources to store
		if applyFunc(resMap, api.Store.Create) != nil {
			res.WriteHeader(http.StatusInternalServerError)
			log.Info("internal server error while creating resources")

			return
		}

		log.Info("successfully created resources")

		res.WriteHeader(http.StatusOK)
	}
}
