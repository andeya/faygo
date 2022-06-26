package goutil

import (
	"bytes"
	"crypto/rand"
	"encoding/base64"
	"errors"
	mrand "math/rand"

	"github.com/henrylee2cn/ameda"
)

// NewRandom creates a new padded Encoding defined by the given alphabet string.
func NewRandom(alphabet string) *Random {
	r := new(Random)
	diff := 64 - len(alphabet)
	if diff < 0 {
		r.substitute = []byte(alphabet[64:])
		r.substituteLen = len(r.substitute)
		alphabet = alphabet[:64]
	} else {
		r.substitute = []byte(alphabet)
		r.substituteLen = len(r.substitute)
		if diff > 0 {
			alphabet += string(bytes.Repeat([]byte{0x00}, diff))
		}
	}
	r.encoding = base64.NewEncoding(alphabet).WithPadding(base64.NoPadding)
	return r
}

// Random random string creator.
type Random struct {
	encoding      *base64.Encoding
	substitute    []byte
	substituteLen int
}

// RandomString returns a base64 encoded securely generated
// random string. It will panic if the system's secure random number generator
// fails to function correctly.
// The length n must be an integer multiple of 4, otherwise the last character will be padded with `=`.
func (r *Random) RandomString(length int) string {
	d := r.encoding.DecodedLen(length)
	buf := make([]byte, length)
	r.encoding.Encode(buf, RandomBytes(d))
	for k, v := range buf {
		if v == 0x00 {
			buf[k] = r.substitute[mrand.Intn(r.substituteLen)]
		}
	}
	return ameda.UnsafeBytesToString(buf)
}

const tsLen = 6 // base62=ZZZZZZ, unix=56800235583, time=3769-12-05 11:13:03 +0800 CST

// RandomStringWithTime returns a random string with UNIX timestamp(in second).
// unixTs: range [0,56800235583], that is 56800235583 3769-12-05 11:13:03 +0800 CST to 3769-12-05 11:13:03 +0800 CST
func (r *Random) RandomStringWithTime(length int, unixTs int64) (string, error) {
	if length <= tsLen {
		return "", errors.New("length is less than 7")
	}
	if unixTs < 0 || unixTs > 56800235583 {
		return "", errors.New("unixTs is out of range [0,56800235583]")
	}
	return r.RandomString(length-tsLen) + ameda.FormatInt(unixTs, 62), nil
}

// ParseTime parses UNIX timestamp(in second) from stringWithTime.
func (r *Random) ParseTime(stringWithTime string) (unixTs int64, err error) {
	length := len(stringWithTime)
	if length <= tsLen {
		return 0, errors.New("stringWithTime length is less than 7")
	}
	return ameda.ParseInt(stringWithTime[length-6:], 62, 64)
}

const urlEncoder = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789-_"

var urlRandom = &Random{
	encoding:      base64.URLEncoding,
	substitute:    []byte(urlEncoder),
	substituteLen: len(urlEncoder),
}

// URLRandom returns Random object with URL encoder.
func URLRandom() *Random {
	return urlRandom
}

// URLRandomString returns a URL-safe, base64 encoded securely generated
// random string. It will panic if the system's secure random number generator
// fails to function correctly.
// The length n must be an integer multiple of 4, otherwise the last character will be padded with `=`.
func URLRandomString(n int) string {
	return urlRandom.RandomString(n)
}

// RandomBytes returns securely generated random bytes. It will panic
// if the system's secure random number generator fails to function correctly.
func RandomBytes(n int) []byte {
	b := make([]byte, n)
	_, err := rand.Read(b)
	// Note that err == nil only if we read len(b) bytes.
	if err != nil {
		panic(err)
	}
	return b
}
