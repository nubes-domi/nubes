package db

import (
	"database/sql/driver"
	"strings"
)

type pipeStringArray []string

func (n *pipeStringArray) Scan(value interface{}) error {
	*n = strings.Split(string(value.(string)), "|")
	return nil
}

func (n *pipeStringArray) Value() (driver.Value, error) {
	return driver.Value(strings.Join(*n, "|")), nil
}
