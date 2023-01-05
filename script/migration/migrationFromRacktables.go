// Migration Script for data - it can be used to fetch,
// add, and use the data inside zebra.
package migration

import (
	"bytes"
	"context"
	"database/sql"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"net/http/httptest"
	"strings"

	//nolint:gci

	"github.com/project-safari/zebra"
	"github.com/project-safari/zebra/model/dc"
	"github.com/project-safari/zebra/script"

	// this is needed for mysql access.
	_ "github.com/go-sql-driver/mysql"
)

// Racktables struct is the struct that contains info from the racktables table in the mysql db.
//
// It contains the id, ip, name, and type of the item in the racktable.
type Racktables struct {
	ID        string `json:"object_id"` //nolint:tagliatelle
	Name      string `json:"name"`
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
	Group     string `json:"systemGroup"`
}

type ResourceAPI struct {
	factory zebra.ResourceFactory
	Store   zebra.Store
}

type CtxKey string

const (
	ResourcesCtxKey = CtxKey("resources")
)

func NewResourceAPI(factory zebra.ResourceFactory) *ResourceAPI {
	return &ResourceAPI{
		factory: factory,
		Store:   nil,
	}
}

// Determine the specific type of a resource.
// nolint
func determineType(means string, resName string) string {
	name := strings.ToLower(resName)
	typ := ""

	if means == "Shelf" {
		typ = "dc.rack"
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

// Get resource type by id.
// nolint
func determineIDMeaning(id string, name string) string {
	means := ""
	final := ""
	this := ""

	if id == "2" || id == "27" {
		means = "compute.vm"
	} else if id == "30" || id == "31" || id == "34" || id == "3" {
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

//nolint:funlen, cyclop
func Do() []Racktables {
	var rt Racktables

	RackArr := []Racktables{}

	// Statement to query the db - currently only one rack, 76.
	statement := "SELECT rack_id, object_id  FROM rackspace WHERE rack_id = 76"
	// to be filled in with appropriate user, password, and db name.
	db, err := sql.Open("mysql", "username:1234@/racktables")
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

		rt.ID, rt.Name, rt.Label, rt.ObjtypeID, rt.AssetNo, rt.Problems, rt.Comments = getMoreDetails(rt.ID, db)

		typeID := rt.ObjtypeID
		resType := determineIDMeaning(typeID, rt.Name)
		rt.Type = resType

		if strings.Contains(resType, "compute") || resType == "network.switch" {
			rt.IP = getIPDetaiLs(rt.ID, db)

			ownedBy := getUserDetails(rt.IP, db)
			rt.Owner = ownedBy

			if rt.IP == "" || net.ParseIP(rt.IP) == nil {
				rt.IP = "127.0.0.1"
			}
		} else {
			rt.IP = "127.0.0.1"

			ownedBy := "null"
			rt.Owner = ownedBy
		}

		if resType == "network.switch" {
			portInfo := getPortDetails(rt.ID, db)
			rt.Port = portInfo
		} else {
			portID := -1
			rt.Port = portID
		}

		rackID := getRackDetails(rt.ID, db)
		rowName, rowID, rowLocation := getRowDetails(rackID, db)

		rt.RowName = rowName
		rt.RowID = rowID
		rt.Location = rowLocation

		assetNumber := rt.AssetNo
		rt.AssetNo = assetNumber

		probs := rt.Problems
		rt.Problems = probs

		notes := rt.Comments
		rt.Comments = notes

		parentTypeID := script.ParentWhatTypeID(typeID, db)
		// parentName := parentWhatName(parentID)

		parID, parDescription := script.ParentChildRelation(rt.ID, db)

		// Currently returning an empt string since table TagTree is still empty.
		// parentID := parentIndividualID(rt.ID)

		parentType := script.DetermineParentType(parentTypeID, typeID, rt.Name)

		group := parentTypeID + "-" + parID + "-" + parentType + "-" + parDescription + "-" + rt.Location
		print("GROUP!!!! ", group, " !!\nn")
		rt.Group = group

		RackArr = append(RackArr, rt)

		if err != nil {
			panic(err.Error())
		}
	}

	if err != nil {
		fmt.Println(err.Error())
	}

	// all data gets posted by reusing handle from API post resources.
	allData(RackArr)

	return RackArr
}

// used to post data.
func Post() {
	postData := Do()
	allData(postData)
}

func allData(rackArr []Racktables) {
	factory := zebra.Factory()

	myAPI := NewResourceAPI(factory)

	h := script.HandlePost()
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		h(w, r, nil)
	})

	for i := 0; i < (len(rackArr)); i++ {
		res := rackArr[i]

		_, _, eachRes := createResFromData(res)

		// Create new resource on zebra with post request.
		req := createRequests("POST", "/resources", eachRes, myAPI)
		rr := httptest.NewRecorder()
		handler.ServeHTTP(rr, req)
	}
}

func createRequests(method string, url string,
	body string, api *ResourceAPI,
) *http.Request {
	ctx := context.WithValue(context.Background(), ResourcesCtxKey, api)
	req, _ := http.NewRequestWithContext(ctx, method, url, nil)

	if body != "" {
		req.Body = ioutil.NopCloser(bytes.NewBufferString(body))
		print("Posted   ", body, "  successfully!\n")
	}

	return req
}

// Get IPs from db based on type id.
func getIPDetaiLs(objectID string, db *sql.DB) string {
	var rt Racktables

	var IPnum string

	statement := "SELECT ip FROM IPv4Allocation WHERE object_id = ?"

	// Execute the query
	results, err := db.Query(statement, objectID)
	if err != nil {
		panic(err.Error()) // proper error handling instead of panic in your app
	}

	defer results.Close()

	for results.Next() {
		// for each row, scan the result into our tag composite object
		err = results.Scan(&IPnum)

		rt.IP = IPnum

		if err != nil {
			panic(err.Error()) // proper error handling instead of panic in your app
		}
	}

	return rt.IP
}

// Get port IDs from db based on type ID.
func getPortDetails(objectID string, db *sql.DB) int {
	var rt Racktables

	numPort := 0

	statement := "SELECT id FROM Port WHERE object_id = ?"

	// Execute the query
	results, err := db.Query(statement, objectID)
	if err != nil {
		panic(err.Error())
	}

	// defer results.Close()

	for results.Next() {
		err = results.Scan(&rt.Port)

		if err != nil {
			panic(err.Error())
		}

		numPort++
	}

	return numPort
}

// Get rack details using the resource's specific ID.
func getMoreDetails(objectID string, db *sql.DB) (string, string, string, string, string, string, string) {
	var rt Racktables

	var label sql.NullString

	var assetNo sql.NullString

	var comment sql.NullString

	statement := "SELECT id, name, label, objtype_id, asset_no, has_problems, comment FROM rackobject WHERE id = ?"

	// Execute the query
	results, err := db.Query(statement, objectID)
	if err != nil {
		panic(err.Error())
	}

	defer results.Close()

	for results.Next() {
		err = results.Scan(&rt.ID, &rt.Name, &label, &rt.ObjtypeID, &assetNo, &rt.Problems, &comment)

		rt.Label = label.String

		rt.AssetNo = assetNo.String

		rt.Comments = comment.String

		if err != nil {
			panic(err.Error())
		}

		log.Print(rt.RackID)
	}

	return rt.ID, rt.Name, rt.Label, rt.ObjtypeID, rt.AssetNo, rt.Problems, rt.Comments
}

func getRackDetails(objID string, db *sql.DB) string {
	var rt Racktables

	statement := "SELECT rack_id FROM RackSpace WHERE object_id = ?"

	// Execute the query
	results, err := db.Query(statement, objID)
	if err != nil {
		panic(err.Error())
	}

	defer results.Close()

	for results.Next() {
		err = results.Scan(&rt.RackID)

		if err != nil {
			panic(err.Error())
		}
	}

	return rt.RackID
}

// Get row and location details based on rack info (rack ID).
func getRowDetails(id string, db *sql.DB) (string, string, string) {
	var rt Racktables

	statement := "SELECT row_id, row_name, location_name FROM Rack WHERE id = ?"

	// Execute the query
	results, err := db.Query(statement, id)
	if err != nil {
		panic(err.Error())
	}
	defer results.Close()

	for results.Next() {
		err = results.Scan(&rt.RowID, &rt.RowName, &rt.Location)

		if err != nil {
			panic(err.Error())
		}
	}

	return rt.RowID, rt.RowName, rt.Location
}

// Get owner / user details based on the resource's IP.
func getUserDetails(resIP string, db *sql.DB) string {
	var rt Racktables

	statement := "SELECT user FROM IPv4Log WHERE ip = ?"

	// Execute the query
	results, err := db.Query(statement, resIP)
	if err != nil {
		panic(err.Error())
	}

	defer results.Close()

	for results.Next() {
		err = results.Scan(&rt.Owner)

		if err != nil {
			panic(err.Error())
		}
	}

	return rt.Owner
}

// Function to create a resource given data obtained from db, guven a certain type.
//
// Returns a zebra.Resource and a string version of the resource struct to be used with APIs.
//
//nolint:cyclop, funlen, lll
func createResFromData(res Racktables) (zebra.Resource, string, string) {
	resType := res.Type

	switch resType {
	case "dc.dataceneter":
		addR := dc.NewDatacenter(res.Location, res.Name, res.Owner, "system.group-datacenter"+"-"+res.Group)

		this := fmt.Sprintf("%v", addR)

		return addR, "dc", this

	case "dc.lab":
		addR := dc.NewLab(res.Name, res.Owner, "system.group-datacenter-lab"+"-"+res.Group)

		this := fmt.Sprintf("%v", addR)

		return addR, "lab", this

	case "dc.rack", "dc.shelf":
		addR := dc.NewRack(res.RowName, res.RowID, res.Name, res.Location, res.Owner, "system.group-datacenter-lab-rack"+"-"+res.Group)

		this := fmt.Sprintf("%v", addR)

		return addR, "rack", this

	case "compute.server":
		addR := serverFiller(res)

		this := fmt.Sprintf("%v", addR)

		return addR, "server", this

	case "compute.esx":
		addR := esxFiller(res)

		this := fmt.Sprintf("%v", addR)

		return addR, "esx", this

	case "compute.vm":
		addR := vmFiller(res)

		this := fmt.Sprintf("%v", addR)

		return addR, "vm", this

	case "compute.vceneter":
		addR := vcenterFiller(res)

		this := fmt.Sprintf("%v", addR)

		return addR, "vcenter", this

	case "network.switch":
		addR := switchFiller(res)

		this := fmt.Sprintf("%v", addR)

		return addR, "switch", this

	case "network.ipaddresspool":
		addR := addressPoolFiller(res)

		this := fmt.Sprintf("%v", addR)

		return addR, "ip-address-pool", this

	case "network.vlanpool":
		addR := vlanFiller(res)

		this := fmt.Sprintf("%v", addR)

		return addR, "vlan-pool", this
	}

	return nil, "", ""
}
