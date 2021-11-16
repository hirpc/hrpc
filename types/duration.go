package types

import (
	"bytes"
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
