package db

import (
	"crypto/rand"
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/lestrrat-go/jwx/jwk"
)

type Model struct {
	ID        string `gorm:"primaryKey"`
	CreatedAt time.Time
	UpdatedAt time.Time
}

type pipeStringArray []string

func (n *pipeStringArray) Scan(value interface{}) error {
	*n = strings.Split(value.(string), "|")
	return nil
}

func (n pipeStringArray) Value() (driver.Value, error) {
	return driver.Value(strings.Join(n, "|")), nil
}

type jwkSet struct {
	Set jwk.Set
}

func (n jwkSet) MarshalJSON() ([]byte, error) {
	if n.Set != nil {
		val, err := n.Value()
		return []byte(val.(string)), err
	} else {
		return []byte("null"), nil
	}
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

func GenID(t string) string {
	alphabet := "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	id := make([]byte, 16)

	randBuf := make([]byte, 16)
	rand.Read(randBuf)

	i := 0
	for i < 16 {
		id[i] = alphabet[int(randBuf[i])%len(alphabet)]
		i += 1
	}

	fmt.Printf("ID: %s", t+"_"+string(id))

	return t + "_" + string(id)
}
