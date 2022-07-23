package zebra

import (
	"errors"
)

// ErrLabel is an error that occurs if labels do not have the default.
// mandatory group label.
// A mandatory group label is a key-value pair that has the key as "group".
// The value of a group label can be the resource type, location, access, or othes.
var ErrLabel = errors.New("missing default mandatory group label")

// The Labels.Validate function ensures that labels have the mandatory default group label.
func (l Labels) Validate() error {
	if _, found := l["group"]; found {
		return nil
	}

	return ErrLabel
}
