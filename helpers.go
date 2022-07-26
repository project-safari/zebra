package zebra

import (
	"github.com/google/uuid"
)

const DefaultMaxDuration = 4

func NewBaseResource(resType Type, labels Labels) *BaseResource {
	id := uuid.New().String()

	if resType.Name == "" {
		resType.Name = "BaseResource"
		resType.Description = "Base Resource"
		resType.Constructor = func() Resource { return new(BaseResource) }
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
