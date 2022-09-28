package zebra

import (
	"database/sql"
	"fmt"
	"log"
)

// Racktables struct is the struct that contains info from the racktables table in the mysql db.
//
// It contains the id, ip, name, and type of the item in the racktable.
type Racktables struct {
	ID   int     `json:"object_id"` //nolint:tagliatelle
	IP   int     `json:"ip"`
	Name string  `json:"name"`
	Type TheType `json:"type"`
}

type TheType struct {
	Category string `json:"type"`
}

// The getRack function gets data from the data base.
//
// It prepares the necessary data from the data base to be used in the zebra tool.
//
// getRack takes in a Type struct and returns an error or nil in the absence thereof.
//
// It executes the function newMeta with the data it extracted from the db.
func GetRack() (Racktables, error) {
	var rt Racktables

	// to be filled in with appropriate user, password, and db name.
	db, err := sql.Open("mysql", "eachim:password@/racktables")
	// if there is an error opening the connection, handle it
	if err != nil {
		log.Print(err.Error())
	}

	defer db.Close()

	// Execute the query
	results, err := db.Query("SELECT object_id, ip, name, type FROM Racktables")
	if err != nil {
		panic(err.Error()) // proper error handling instead of panic in your app
	}

	for results.Next() {
		// for each row, scan the result into our tag composite object
		err = db.QueryRow("SELECT object_id, ip, name, type WHERE Type = ?", "").Scan(&rt.ID, &rt.IP, &rt.Name, &rt.Type)

		// err = results.Scan(&tag.ID, &tag.Name)
		if err != nil {
			panic(err.Error()) // proper error handling instead of panic in your app
		}
		// and then print out the tag's Name attribute
		log.Printf(rt.Name)
	}

	err = db.QueryRow("SELECT object_id, ip, name, type WHERE Type = ?", "resType").Scan(&rt.ID, &rt.IP, &rt.Name, &rt.Type) //nolint:lll

	if err != nil {
		fmt.Println(err.Error())
	}

	// Save a copy of each piece of data.
	id := rt.ID
	ip := rt.IP

	name := rt.Name
	rackType := rt.Type

	// Return the data retreived from the db to use in the zebra tool.
	return Racktables{
		ID:   id,
		IP:   ip,
		Name: name,
		Type: rackType,
	}, err
}
