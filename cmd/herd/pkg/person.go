package pkg

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
func Password(user string) string {
	pwd := user + "123"

	return pwd
}
