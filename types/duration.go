package types

import (
	"bytes"
	"database/sql/driver"
	"errors"
	"reflect"
	"time"
)

// Duration can be used for "ns", "us" (or "µs"), "ms", "s", "m", "h" suffixs
type Duration struct {
	time.Duration
}

// UnmarshalJSON for implementing the Unmarshaler interface
func (d *Duration) UnmarshalJSON(data []byte) error {
	duration, err := time.ParseDuration(string(bytes.Trim(data, "\"")))
	if err != nil {
		return err
	}
	*d = Duration{duration}
	return nil
}

func (d Duration) MarshalJSON() ([]byte, error) {
	return []byte("\"" + d.String() + "\""), nil
}

func (d Duration) Value() (driver.Value, error) {
	return d.String(), nil
}

func (d *Duration) Scan(src interface{}) error {
	v, ok := src.([]byte)
	if !ok {
		return errors.New(
			"bad []byte type assertion, got name: " + reflect.TypeOf(src).Name() + " kind: " + reflect.TypeOf(src).Kind().String(),
		)
	}
	d1, err := time.ParseDuration(string(v))
	if err != nil {
		return err
	}
	*d = Duration{d1}
	return nil
}
