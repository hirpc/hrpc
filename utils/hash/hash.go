package hash

import (
	"crypto/md5"
	"crypto/sha1"
	"crypto/sha256"
	"fmt"
	"io"
	"strings"
)

func MD5(d ...string) string {
	t := md5.New()
	io.WriteString(t, strings.Join(d, ""))
	return fmt.Sprintf("%x", t.Sum(nil))
}

func SHA1(d ...string) string {
	t := sha1.New()
	io.WriteString(t, strings.Join(d, ""))
	return fmt.Sprintf("%x", t.Sum(nil))
}

func SHA256(d ...string) string {
	t := sha256.New()
	io.WriteString(t, strings.Join(d, ""))
	return fmt.Sprintf("%x", t.Sum(nil))
}
