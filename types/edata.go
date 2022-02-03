package types

import (
	"database/sql/driver"
	"encoding/base64"
	"fmt"
	"reflect"

	"github.com/hirpc/hrpc/uerror"
)

var (
	ErrInvalidKey = uerror.New(100, "invalid key")
)

// EData for encrypted data
type EData struct {
	origin string
	key    string
}

// NewEData the constructor
func NewEData(s, k string) EData {
	return EData{
		origin: s,
		key:    k,
	}
}

// Value for implementing the driver.Valuer interface
func (e EData) Value() (driver.Value, error) {
	if e.key == "" {
		return nil, ErrInvalidKey
	}
	return e, nil
}

// Scan for implementing the sql.Scanner interface
func (e *EData) Scan(src interface{}) error {
	v, ok := src.([]byte)
	if !ok {
		return fmt.Errorf(
			"bad []byte type assertion, got name: %v and kind: %v",
			reflect.TypeOf(src).Name(),
			reflect.TypeOf(src).Kind().String(),
		)
	}
	d, err := base64.StdEncoding.DecodeString(string(v))
	if err != nil {
		return err
	}
	ad, err := Decrypt(d, []byte(e.key))
	if err != nil {
		*e = EData{origin: string(d)}
		return nil
	}
	*e = EData{origin: string(ad)}
	return nil
}

func (e EData) String() string {
	return e.origin
}

func Decrypt(d, k []byte) ([]byte, error) {
	return d, nil
}
