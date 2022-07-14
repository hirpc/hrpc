package types

import (
	"bytes"
	"database/sql/driver"
	"errors"
	"fmt"
	"reflect"
	"strconv"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/bsontype"
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

func (t *Timestamp) UnmarshalJSON(data []byte) error {
	s := string(bytes.Trim(data, "\""))
	v, err := strconv.ParseInt(s, 10, 64)
	if err != nil {
		return err
	}
	var nt = Timestamp{
		Time: time.Unix(v, 0),
	}
	*t = nt
	return nil
}

func (t Timestamp) MarshalBSONValue() (bsontype.Type, []byte, error) {
	return bson.MarshalValue(t.Time)
}

func (t *Timestamp) UnmarshalBSONValue(bType bsontype.Type, data []byte) error {
	rv := bson.RawValue{Type: bType, Value: data}
	return rv.Unmarshal(&t.Time)
}
