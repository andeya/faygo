package utils

import (
	"strings"
)

type Paths struct {
	list []string
}

func NewPaths(list []string) *Paths {
	return &Paths{
		list: list,
	}
}
func (this *Paths) Add(p ...string) *Paths {
	this.list = append(this.list, p...)
	return this
}
func (this *Paths) List() []string {
	return this.list
}
func (this *Paths) Filter(substr string) (paths []string) {
	for _, s := range this.list {
		if strings.Contains(s, substr) {
			paths = append(paths, s)
		}
	}
	return
}
func (this *Paths) PrefixFilter(prefix string) (paths []string) {
	for _, s := range this.list {
		if strings.HasPrefix(s, prefix) {
			paths = append(paths, s)
		}
	}
	return
}
func (this *Paths) SuffixFilter(suffix string) (paths []string) {
	for _, s := range this.list {
		if strings.HasSuffix(s, suffix) {
			paths = append(paths, s)
		}
	}
	return
}

func (this *Paths) TrimPrefix(prefix ...string) *Paths {
	if len(prefix) == 0 {
		return this
	}
	for i, s := range this.list {
		for _, substr := range prefix {
			this.list[i] = strings.TrimPrefix(s, substr)
		}
	}
	return this
}

func (this *Paths) TrimSuffix(suffix ...string) *Paths {
	if len(suffix) == 0 {
		return this
	}
	for i, s := range this.list {
		for _, substr := range suffix {
			this.list[i] = strings.TrimSuffix(s, substr)
		}
	}
	return this
}
