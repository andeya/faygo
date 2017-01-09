package utils

import (
	"sync"
	"testing"
)

func TestRandomBytes(t *testing.T) {
	m := map[string]bool{}
	var lock sync.Mutex
	var group sync.WaitGroup
	count := 5000
	group.Add(count)
	for i := 0; i < count; i++ {
		go func() {
			id := string(RandomBytes(6))
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
	for id := range m {
		i++
		if i > 10 {
			break
		}
		t.Log(id)
	}
}
