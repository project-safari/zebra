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
	} // some passwords.

	pwd := RandData(Plist)

	return pwd
}
