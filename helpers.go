package zebra

import (
	"github.com/google/uuid"
)

func NewBaseResource(resType string, labels Labels) *BaseResource {
	id := uuid.New().String()

	if resType == "" {
		resType = "BaseResource"
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
