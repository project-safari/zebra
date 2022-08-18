package pkg

import "github.com/project-safari/zebra"

// Helper function to add mandatory group label.
// Call this function if the given resoure does not have a group label.
func GroupLabels(l zebra.Labels, groupValue string) zebra.Labels {
	groupLabel := l.Add("system.group", groupValue)

	return groupLabel
}

// Function to generate group Value based on resource type.
//
// This function will be used for gereration of sample.group labels for resources that lack such labels.
//
// Returns a string that contains a group value.
//
// Group value could be a geographic location, a building name/number, a user type or role.
func GroupVal(resource zebra.Resource) string {
	groupSamples := []string{
		"Americas", "admins", "Building15",
		"Oceania", "engineers", "Building2",
		"designers", "Europe", "leadership",
		"Building7", "Asia", "users",
	}

	groupValue := RandData(groupSamples)

	return groupValue
}
