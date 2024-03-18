package postgresql

import (
	"errors"
	"fmt"
)

var (
	ErrConnectionFailed = errors.New("connection failed")
)

type ErrRecordNotFound struct {
	tableName string
	identity  string
}

func (e *ErrRecordNotFound) Error() string {
	return fmt.Sprintf("record not found in %s with %s", e.tableName, e.identity)
}
