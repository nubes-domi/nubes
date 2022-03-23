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
	ID        string    `gorm:"primaryKey" json:"id" binding:"-"`
	CreatedAt time.Time `json:"created_at" binding:"-"`
	UpdatedAt time.Time `json:"updated_at" binding:"-"`
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

type JSONDate struct {
	time.Time
}

func (d JSONDate) MarshalJSON() ([]byte, error) {
	return []byte(`"` + d.Time.Format("2006-01-02") + `"`), nil
}

func (d *JSONDate) UnmarshalJSON(data []byte) error {
	t, err := time.Parse("\"2007-01-02\"", string(data))
	d.Time = t
	return err
}

func (d *JSONDate) Scan(value interface{}) error {
	date, ok := value.(time.Time)
	if ok {
		d.Time = date
		return nil
	}

	t, err := time.Parse("2007-01-02", value.(string))
	d.Time = t
	return err
}

func (d JSONDate) Value() (driver.Value, error) {
	return d.Time.Format("2007-01-02"), nil
}
