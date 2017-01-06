// Copyright 2013, Örjan Persson. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package logging

// TODO share more code between these tests
func MemoryRecordN(b *MemoryBackend, n int) *Record {
	node := b.Head()
	for i := 0; i < n; i++ {
		if node == nil {
			break
		}
		node = node.Next()
	}
	if node == nil {
		return nil
	}
	return node.Record
}

func ChannelMemoryRecordN(b *ChannelMemoryBackend, n int) *Record {
	b.Flush()
	node := b.Head()
	for i := 0; i < n; i++ {
		if node == nil {
			break
		}
		node = node.Next()
	}
	if node == nil {
		return nil
	}
	return node.Record
}
