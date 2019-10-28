package post

import (
	"errors"
)

var (
	errAffectNoRows     = errors.New("Affect No Rows")
	errContentNotChange = errors.New("Content Not Change")
)
