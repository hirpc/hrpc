package types

import (
	"database/sql/driver"
	"errors"
	"fmt"
	"reflect"
	"time"
)

type Timestamp struct {
	time.Time
}

func (t Timestamp) Value() (driver.Value, error) {
	return t.Unix(), nil
}

func (t *Timestamp) Scan(src interface{}) error {
	v, ok := src.(int64)
	if !ok {
		return errors.New(
			"bad int64 type assertion, got name: " + reflect.TypeOf(src).Name() + " kind: " + reflect.TypeOf(src).Kind().String(),
		)
	}
	*t = Timestamp{time.Unix(v, 0)}
	return nil
}

func (t Timestamp) MarshalJSON() ([]byte, error) {
	return []byte(fmt.Sprintf("%d", t.Unix())), nil
}
