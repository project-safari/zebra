package zebra

import "errors"

var ErrLabel = errors.New("missing mandatory system label")

// Labels are key/value pairs that can be attached to any resource. They are
// used to identify attributes of a resource that are added or modified during
// the runtime operations on a resource. Resources can be matched using labels
// to select or not select for a particular operation. For example a workflow
// can select all matching resources that has a label color=red and display
// them for user inspection.
type Labels map[string]string

// Add adds an label with specified key/value to this label set. If the label
// is already present, Add will simply update the value.
func (l Labels) Add(key, value string) Labels {
	l[key] = value

	return l
}

// HasKey returns if the label key is present in the label set.
func (l Labels) HasKey(key string) bool {
	_, ok := l[key]

	return ok
}

// MatchEqual returns true if the specified label is present in the label set,
// and the given value is equal to the label's value, else returns false.
func (l Labels) MatchEqual(key, value string) bool {
	return l[key] == value
}

// MatchNotEqual returns true if the specified label key is present and the
// value does not match the specified value for that key, else returns false.
func (l Labels) MatchNotEqual(key, value string) bool {
	v, ok := l[key]

	return ok && v != value
}

// MatchIn returns true if the specified label key is present and the value
// matches at least one of the values specified, else returns false.
func (l Labels) MatchIn(key string, values ...string) bool {
	v, ok := l[key]
	if ok {
		return IsIn(v, values)
	}

	return false
}

// MatchNotIn returns true if the specified label key is present and the value
// does not match any specified value for that key, else returns false.
func (l Labels) MatchNotIn(key string, values ...string) bool {
	return l.HasKey(key) && !l.MatchIn(key, values...)
}

func (l Labels) Validate() error {
	v, ok := l["system.group"]
	if !ok || v == "" {
		return ErrLabel
	}

	return nil
}

// IsIn returns if val is in string list.
func IsIn(val string, list []string) bool {
	for _, v := range list {
		if val == v {
			return true
		}
	}

	return false
}
