package service

import "errors"

var (
	ErrEntryCreateFailed = errors.New("failed to create a new entry")
	ErrEntryNotFound     = errors.New("entry not found")
)
