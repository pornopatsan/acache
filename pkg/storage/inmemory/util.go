package inmemory

import "fmt"

type KeyAlreadyExistsError struct {
	Key string
}

func (e KeyAlreadyExistsError) Error() string {
	return fmt.Sprintf("Key `%s` already exists", e.Key)
}

type KeyNotFoundError struct {
	Key string
}

func (e KeyNotFoundError) Error() string {
	return fmt.Sprintf("Key `%s` not found", e.Key)
}
