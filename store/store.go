package store

import (
	"sync"

	"github.com/project-safari/zebra/filestore"
	"github.com/project-safari/zebra/idstore"
	"github.com/project-safari/zebra/labelstore"
	"github.com/project-safari/zebra/typestore"
)

//nolint:unused
type ResourceStore struct {
	lock sync.RWMutex
	fs   *filestore.FileStore
	ids  *idstore.IDStore
	ls   *labelstore.LabelStore
	ts   *typestore.TypeStore
}
