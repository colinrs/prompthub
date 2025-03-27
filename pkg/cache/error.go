package cache

import (
	"errors"
)

var (
	ErrNotFound   = errors.New("not found")
	ErrFromSource = errors.New("from source err")
)
