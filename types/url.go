package types

import (
	"bytes"
	"database/sql/driver"
	"fmt"
	"net/url"
	"reflect"
)

type URL struct {
	*url.URL
}

func (u URL) Value() (driver.Value, error) {
	return u.String(), nil
}

func (u *URL) Scan(src interface{}) error {
	v, ok := src.([]byte)
	if !ok {
		return fmt.Errorf(
			"bad []byte type assertion, got name: %v and kind: %v",
			reflect.TypeOf(src).Name(),
			reflect.TypeOf(src).Kind().String(),
		)
	}
	t, err := url.Parse(string(v))
	if err != nil {
		return err
	}
	*u = URL{t}
	return nil
}

func (u URL) MarshalJSON() ([]byte, error) {
	return []byte("\"" + u.String() + "\""), nil
}

func (u *URL) UnmarshalJSON(data []byte) error {
	s := string(bytes.Trim(data, "\""))
	t, err := url.Parse(s)
	if err != nil {
		return err
	}
	*u = URL{t}
	return nil
}
