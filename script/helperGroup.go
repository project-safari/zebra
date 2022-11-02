package script

import (
	"database/sql"
	"strings"
)

// Determine the specific type of a resource.
// nolint
func DetermineType(means string, resName string) string {
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
		} else if strings.Contains(name, "lab") {
			typ = "dc.lab"
		} else if strings.Contains(name, "infra") || strings.Contains(name, "vshield") || strings.Contains(name, "firewall") {
			typ = "infrastructure"
		}
	} else if means == "Other" {
		if strings.Contains(name, "chasis") || strings.Contains(name, "ixia") || strings.Contains(name, "rack") {
			typ = "dc.rack"
		} else if strings.Contains(name, "nexus") || strings.Contains(name, "sw") || strings.Contains(name, "switch") || strings.Contains(name, "n3k") {
			typ = "network.switch"
		} else if strings.Contains(name, "lab") {
			typ = "dc.lab"
		} else if strings.Contains(name, "vcenter") {
			typ = "compute.vcenter"
		} else if strings.Contains(name, "infra") || strings.Contains(name, "vshield") || strings.Contains(name, "firewall") {
			typ = "infrastructure"
		} else {
			typ = "unclassified"
		}
	} else {
		typ = means
	}

	return typ
}

// Function that gets the parent type based on the type of the child.
//
//nolint:ineffassign
func GetParent(childStrType string) string {
	means := "/"

	switch childStrType {
	case "compute.esx":
		means = "compute.server"
	case "compute.vm":
		means = "compute.esx"
	case "compute.vcenter":
		means = "compute.esx"
	case "dc.rack":
		means = "dc.lab"
	case "dc.lab":
		means = "dc.datacenter"
	case "compute.server", "compute.switch":
		means = "dc.rack"
	default:
		means = childStrType
	}

	return means
}

// Determine type of a parent based on the parent and child IDs.
//
//nolint:cyclop
func DetermineParentType(parentID string, childID string, childName string) string {
	means := "/"

	switch parentID {
	case "4":
		if childID == "1504" {
			typeMeans := DetermineType("Compute", childName)
			means = GetParent(typeMeans)
		} else if childID == "1507" {
			typeMeans := DetermineType("Other", childName)
			means = GetParent(typeMeans)
		}
	case "1503":
		if childID == "1504" {
			typeMeans := DetermineType("Compute", childName)
			means = GetParent(typeMeans)
		} else if childID == "1507" {
			if DetermineType("Other", childName) == "network.switch" || childID == "8" {
				means = "dc.rack"
			} else if DetermineType("Other", childName) == "dc.rack" || childID == "4" {
				means = "dc.lab"
			}
		}
	default:
		typeMeans := DetermineType("Other", childName)
		means = GetParent(typeMeans)
	}

	return means
}

// Get the parent ID information from the database.
func ProvideParent(objtypeID string, db *sql.DB) string {
	parentType := ""
	statement := "SELECT parent_objtype_id  FROM ObjectParentCompat WHERE  child_objtype_id = ?"

	// Execute the query
	results, err := db.Query(statement, objtypeID)
	if err != nil {
		panic(err.Error())
	}
	defer results.Close()

	for results.Next() {
		err = results.Scan(&parentType)

		if err != nil {
			panic(err.Error())
		}
	}

	return parentType
}

// Get the parent's type ID from the database.
func ParentWhatTypeID(childType string, db *sql.DB) string {
	parent := ""

	statement := "SELECT parent_objtype_id FROM ObjectParentCompat WHERE child_objtype_id = ?"
	// to be filled in with appropriate user, password, and db name.

	// Execute the query
	results, err := db.Query(statement, childType)
	if err != nil {
		panic(err.Error())
	}

	defer results.Close()

	for results.Next() {
		err = results.Scan(&parent)

		if err != nil {
			panic(err.Error())
		}
	}

	return parent
}

// To be implemented once data is uploaded to the TagTree table.
// This will get the individual ID of each parent for a given individual resource ID.
func ParentIndividualID(childIndividualID string, db *sql.DB) string {
	parIndividualID := ""

	statement := "SELECT parent_id FROM TagTree WHERE id = ?"

	// Execute the query
	results, err := db.Query(statement, childIndividualID)
	if err != nil {
		panic(err.Error())
	}

	defer results.Close()

	for results.Next() {
		err = results.Scan(&parIndividualID)

		if err != nil {
			panic(err.Error())
		}
	}

	return parIndividualID
}

// Get the full group/relation between a resource and its parent based on all the obtained data.
func ParentChildRelation(childentityID string, db *sql.DB) (string, string) {
	parEntityID := ""
	parEntityDescription := ""

	statement := "SELECT parent_entity_id, parent_entity_type  FROM EntityLink WHERE child_entity_id  = ?"

	// Execute the query
	results, err := db.Query(statement, childentityID)
	if err != nil {
		panic(err.Error())
	}
	defer results.Close()

	for results.Next() {
		err = results.Scan(&parEntityID, &parEntityDescription)

		if err != nil {
			panic(err.Error())
		}
	}

	if parEntityID == "" {
		parEntityID = "NA"
	} else if parEntityDescription == "" || len(parEntityDescription) == 0 {
		parEntityDescription = "NA"
	}

	return parEntityID, parEntityDescription
}
