package utils

import (
	"crypto/rand"
	"io"
	"sync"
	"testing"
)

const tokenLength = 32

// shortReader provides a broken implementation of io.Reader for testing.
type shortReader struct{}

func (sr shortReader) Read(p []byte) (int, error) {
	return len(p) % 2, io.ErrUnexpectedEOF
}

// TestRandomBytes tests the (extremely rare) case that crypto/rand does
// not return the expected number of bytes.
func TestRandomBytes(t *testing.T) {
	// Pioneered by https://github.com/justinas/nosurf
	original := rand.Reader
	rand.Reader = shortReader{}
	defer func() {
		rand.Reader = original
	}()

	var b = make([]byte, tokenLength)
	defer func() {
		if err := recover(); err == nil {
			t.Fatalf("RandomBytes did not report a short read: only read %d bytes", len(b))
		}
	}()

	b = RandomBytes(tokenLength)
}

func TestRandomString(t *testing.T) {
	m := map[string]bool{}
	var lock sync.Mutex
	var group sync.WaitGroup
	count := 10000
	group.Add(count)
	for i := 0; i < count; i++ {
		go func() {
			id := RandomString(10)
			lock.Lock()
			m[id] = true
			lock.Unlock()
			group.Done()
		}()
	}
	group.Wait()
	if len(m) != count {
		t.Fail()
	}
	var i int
	t.Log("print the top ten...")
	for id := range m {
		i++
		if i > 10 {
			break
		}
		t.Log(id)
	}
}
