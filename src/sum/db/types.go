package db

import (
	"database/sql/driver"
	"encoding/json"
	"strings"

	"github.com/lestrrat-go/jwx/jwk"
)

type pipeStringArray []string

func (n *pipeStringArray) Scan(value interface{}) error {
	*n = strings.Split(value.(string), "|")
	return nil
}

func (n *pipeStringArray) Value() (driver.Value, error) {
	return driver.Value(strings.Join(*n, "|")), nil
}

type jwkSet struct {
	Set jwk.Set
}

func (n jwkSet) MarshalJSON() ([]byte, error) {
	val, err := n.Value()
	return []byte(val.(string)), err
}

func (n *jwkSet) UnmarshalJSON(data []byte) error {
	set, err := jwk.Parse(data)
	if err != nil {
		return err
	}

	*n = jwkSet{set}
	return nil
}

func (n *jwkSet) Scan(value interface{}) error {
	val := value.(string)
	if val != "" {
		set, err := jwk.Parse([]byte(val))
		if err != nil {
			return err
		}
		*n = jwkSet{set}
		return nil
	} else {
		*n = jwkSet{nil}
		return nil
	}
}

func (n jwkSet) Value() (driver.Value, error) {
	if n.Set == nil {
		return "", nil
	}

	stream, err := json.Marshal(n.Set)
	return string(stream), err
}
