package store

import "errors"

var (
	ErrCatalogNotFound        = errors.New("catalog item not found")
	ErrAnthologyNotFound      = errors.New("anthology item not found")
	ErrAnthologyNotInCatalog  = errors.New("anthology item not found")
	ErrAnthologyAlreadyExists = errors.New("anthology item already exists")
	ErrUserDoesNotExist       = errors.New("user not found")
	ErrUserAlreadyExists      = errors.New("user id already exists")
)
