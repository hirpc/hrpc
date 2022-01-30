package types

import (
	"database/sql/driver"
	"encoding/binary"
	"errors"
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
	ts := binary.BigEndian.Uint64(v)
	*t = Timestamp{time.Unix(int64(ts), 0)}
	return nil
}
