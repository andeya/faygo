package thinkgo

import (
	"bytes"
	"errors"
	"io/ioutil"
	"net/http"
	"sync"
	"time"

	"github.com/henrylee2cn/thinkgo/pongo2"
)

type (
	Tpl struct {
		template *pongo2.Template
		modTime  time.Time
	}
	// Pongo2Render is a custom thinkgo template renderer using Pongo2.
	Pongo2Render struct {
		set           *pongo2.TemplateSet
		tplCache      map[string]*Tpl
		tplContext    pongo2.Context // Context hold globle func for tpl
		openCacheFile func(name string) (http.File, error)
		caching       bool // false=disable caching, true=enable caching
		sync.RWMutex
	}
)

// New creates a new Pongo2Render instance with custom Options.
func newPongo2Render(openCacheFile func(name string) (http.File, error)) *Pongo2Render {
	return &Pongo2Render{
		set:           pongo2.NewSet("thinkgo", pongo2.DefaultLoader),
		tplCache:      make(map[string]*Tpl),
		tplContext:    make(pongo2.Context),
		openCacheFile: openCacheFile,
		caching:       openCacheFile != nil,
	}
}

// Sets the global template variable or function
func (p *Pongo2Render) TemplateVariable(name string, v interface{}) {
	switch d := v.(type) {
	case func(in *pongo2.Value, param *pongo2.Value) (out *pongo2.Value, err *pongo2.Error):
		pongo2.RegisterFilter(name, d)
	case pongo2.FilterFunction:
		pongo2.RegisterFilter(name, d)
	default:
		p.tplContext[name] = d
	}
}

// Render should render the template to the io.Writer.
func (p *Pongo2Render) Render(filename string, data Map) ([]byte, error) {
	var data2 pongo2.Context
	if data == nil {
		data2 = p.tplContext

	} else {
		data2 = pongo2.Context(data)
		for k, v := range p.tplContext {
			if _, ok := data2[k]; !ok {
				data2[k] = v
			}
		}
	}

	var template *pongo2.Template

	if p.caching {
		template = pongo2.Must(p.FromCache(filename))
	} else {
		template = pongo2.Must(p.set.FromFile(filename))
	}
	var b bytes.Buffer
	err := template.ExecuteWriter(data2, &b)
	return b.Bytes(), err
}

func (p *Pongo2Render) FromCache(fname string) (*pongo2.Template, error) {
	// Get file content from the file system cache
	f, err := p.openCacheFile(fname)
	if err != nil {
		p.Lock()
		_, has := p.tplCache[fname]
		if has {
			delete(p.tplCache, fname)
		}
		p.Unlock()
		return nil, errors.New(fname + "is not found.")
	}

	p.RLock()
	tpl, has := p.tplCache[fname]
	p.RUnlock()
	stat, _ := f.Stat()

	// When the template cache exists and the file is not updated
	if has && p.tplCache[fname].modTime.Equal(stat.ModTime()) {
		return tpl.template, nil
	}

	// The cache template does not exist or the file is updated
	p.Lock()
	defer p.Unlock()

	// Create a new template and cache it
	fbytes, _ := ioutil.ReadAll(f)
	newtpl, err := p.set.FromBytesWithName(fname, fbytes)
	if err != nil {
		return nil, err
	}

	p.tplCache[fname] = &Tpl{template: newtpl, modTime: stat.ModTime()}
	return newtpl, nil
}
