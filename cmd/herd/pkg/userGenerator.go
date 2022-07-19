package pkg

import "github.com/project-safari/zebra/auth"

// generate some user info.
func GenerateUser(numUsr int) []*auth.User {
	users := make([]*auth.User, 0, numUsr)

	for i := 0; i < numUsr; i++ {
		role := User()
		name := Name()
		pwd := Password()
		labels := CreateLabels()

		usr := auth.NewUser(role, name, pwd, labels)
		users = append(users, usr)
	}

	return users
}
