package error

import "errors"

var (
	ErrDBError        = errors.New("something went wrong with the database")
	ErrRecordNotFound = errors.New("record not found")
)
