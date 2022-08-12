package zebra

import (
	"github.com/google/uuid"
)

// the default maximum duration, 4 hours.
const DefaultMaxDuration = 4

// function that creates a new base resource, which is needed in all resource types.
// returns a pointer to a BaseResource.
func NewBaseResource(resType string, labels Labels) *BaseResource {
	id := uuid.New().String()

	if resType == "" {
		resType = "BaseResource"
	}

	if labels == nil {
		labels = Labels{"system.group": "default"}
	} else if labels["system.group"] == "" {
		labels.Add("system.group", "default")
	}

	return &BaseResource{
		ID:     id,
		Type:   resType,
		Labels: labels,
		Status: DefaultStatus(),
	}
}

// Return if val is in string list.
func IsIn(val string, list []string) bool {
	for _, v := range list {
		if val == v {
			return true
		}
	}

	return false
}
