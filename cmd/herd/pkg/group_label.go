package pkg

import "github.com/project-safari/zebra"

// helper function to add mandatory group label.
// call this function if the given resoure does not have a group label.
func GroupLabels(l zebra.Labels, groupValue string) zebra.Labels {
	groupLabel := l.Add("group", groupValue)

	return groupLabel
}

// group Value based on resource type.
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
