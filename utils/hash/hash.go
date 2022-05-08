package hash

import (
	"crypto/md5"
	"crypto/sha1"
	"crypto/sha256"
	"fmt"
	"io"
)

func MD5(vs ...interface{}) string {
	h := md5.New()
	for _, v := range vs {
		switch val := v.(type) {
		case string:
			io.WriteString(h, val)
		case []byte:
			h.Write(val)
		default:
			io.WriteString(h, fmt.Sprintf("%v", val))
		}
	}
	return fmt.Sprintf("%x", h.Sum(nil))
}

func SHA1(vs ...interface{}) string {
	h := sha1.New()
	for _, v := range vs {
		switch val := v.(type) {
		case string:
			io.WriteString(h, val)
		case []byte:
			h.Write(val)
		default:
			io.WriteString(h, fmt.Sprintf("%v", val))
		}
	}
	return fmt.Sprintf("%x", h.Sum(nil))
}

func SHA256(vs ...interface{}) string {
	h := sha256.New()
	for _, v := range vs {
		switch val := v.(type) {
		case string:
			io.WriteString(h, val)
		case []byte:
			h.Write(val)
		default:
			io.WriteString(h, fmt.Sprintf("%v", val))
		}
	}
	return fmt.Sprintf("%x", h.Sum(nil))
}
