package zebra

import "github.com/google/uuid"

func NewBaseResource(resType string, labels Labels) *BaseResource {
	id := uuid.New().String()

	if resType == "" {
		resType = "BaseResource"
	}

	return &BaseResource{
		ID:     id,
		Type:   resType,
		Labels: labels,
	}
}
