package zebra

import "errors"

type Operator uint8

// Constants defined for QueryOperator type.
const (
	MatchEqual Operator = iota
	MatchNotEqual
	MatchIn
	MatchNotIn
)

// Command struct for label queries.
type Query struct {
	Op     Operator
	Key    string
	Values []string
}

var (
	ErrNotFound        = errors.New("resource not found in store")
	ErrInvalidResource = errors.New("create/delete on invalid resource")
	ErrInvalidQuery    = errors.New("invalid query")
)

// Store interface requires basic store functionalities.
type Store interface {
	Initialize() error
	Wipe() error
	Clear() error
	Load() (*ResourceMap, error)
	Create(res Resource) error
	Delete(res Resource) error
	Query() *ResourceMap
	QueryUUID(uuids []string) *ResourceMap
	QueryType(types []string) *ResourceMap
	QueryLabel(query Query) (*ResourceMap, error)
	QueryProperty(query Query) (*ResourceMap, error)
}

func (q *Query) Validate() error {
	if (q.Op == MatchEqual || q.Op == MatchNotEqual) && len(q.Values) != 1 {
		return ErrInvalidQuery
	}

	if q.Op > MatchNotIn {
		return ErrInvalidQuery
	}

	return nil
}
