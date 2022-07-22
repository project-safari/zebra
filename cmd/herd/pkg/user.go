package pkg

import (
	"github.com/project-safari/zebra"
	"github.com/project-safari/zebra/auth"
)

// generate some user info.
func GenerateUser(numUsr int) []zebra.Resource {
	users := make([]zebra.Resource, 0, numUsr)

	// Generate only one key its too costly an operation
	key, _ := auth.Generate()
	key = key.Public()

	for i := 0; i < numUsr; i++ {
		name := Name()
		email := Email(name)
		pwd := Password(name)
		labels := CreateLabels()

		usr := auth.NewUser(name, email, pwd, key, labels)
		users = append(users, usr)
	}

	return users
}
