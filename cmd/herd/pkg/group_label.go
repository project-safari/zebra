package pkg

import "github.com/project-safari/zebra"

// helper function to add mandatory group label.
func GroupLabels(l zebra.Labels, groupValue string) zebra.Labels {
	groupLabel := l.Add("group", groupValue)

	return groupLabel
}

// group Value based on resource type.
func GroupVal(resource zebra.Resource) string {
	groupValue := resource.GetType()

	return groupValue
}
