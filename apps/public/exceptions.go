package public

import (
	"errors"
)

var (
	ErrIdParamIsMissing error = errors.New("id param is missing")
	ErrInvalidUniqueId  error = errors.New("invalid unique id")
)
