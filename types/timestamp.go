package types

import (
	"database/sql/driver"
	"errors"
	"fmt"
	"strconv"
	"time"
)

type Timestamp struct {
	time.Time
}

func (t Timestamp) Value() (driver.Value, error) {
	return t.Unix(), nil
}

func (t *Timestamp) Scan(src interface{}) error {
	v, ok := src.([]byte)
	if !ok {
		return errors.New("bad []byte type assertion")
	}
	ts, err := strconv.ParseInt(string(v), 10, 64)
	if err != nil {
		return err
	}
	*t = Timestamp{time.Unix(ts, 0)}
	return nil
}

func (t Timestamp) MarshalJSON() ([]byte, error) {
	return []byte(fmt.Sprintf("%d", t.Unix())), nil
}
