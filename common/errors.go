package common

import "errors"

var (
	ErrEmptyOrderProducts          = errors.New("required at least one product")
	ErrInvalidOrderProductQuantity = errors.New("invalid order product quantity must greater than zero")
)
