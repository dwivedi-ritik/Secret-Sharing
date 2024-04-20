package public

import (
	"errors"
)

var (
	ErrIdParamIsMissing error = errors.New("id param is missing")
)
