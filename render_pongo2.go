package faygo

import (
	"bytes"
	"errors"
	"io/ioutil"
	"net/http"
	"os"
	"sync"
	"time"

	"github.com/henrylee2cn/faygo/pongo2"
)

type (
	// Tpl is template with modfied time.
	Tpl struct {
		template *pongo2.Template
		modTime  time.Time
	}
	// Render is a custom faygo template renderer using pongo2.
	Render struct {
		set           *pongo2.TemplateSet
		tplCache      map[string]*Tpl
		tplContext    pongo2.Context // Context hold globle func for tpl
		openCacheFile func(name string) (http.File, error)
		caching       bool // false=disable caching, true=enable caching
		sync.RWMutex
	}
)

// New creates a new Render instance with custom Options.
func newRender(openCacheFile func(name string) (http.File, error)) *Render {
	return &Render{
		set:           pongo2.NewSet("faygo", pongo2.DefaultLoader),
		tplCache:      make(map[string]*Tpl),
		tplContext:    make(pongo2.Context),
		openCacheFile: openCacheFile,
		caching:       openCacheFile != nil,
	}
}

// TemplateVar sets the global template variable or function
func (render *Render) TemplateVar(name string, v interface{}) {
	switch d := v.(type) {
	case func(in *pongo2.Value, param *pongo2.Value) (out *pongo2.Value, err *pongo2.Error):
		pongo2.RegisterFilter(name, d)
	case pongo2.FilterFunction:
		pongo2.RegisterFilter(name, d)
	default:
		render.tplContext[name] = d
	}
}

// Render should render the template to the io.Writer.
func (render *Render) Render(filename string, data Map) ([]byte, error) {
	if render.caching {
		b, _, err := render.fromCache(filename, data, false)
		return b, err
	}
	f, err := os.Open(filename)
	if err != nil {
		return nil, err
	}

	var fbytes []byte
	fbytes, err = ioutil.ReadAll(f)
	if err != nil {
		return nil, err
	}
	fbytes, err = render.RenderFromBytes(fbytes, data)
	return fbytes, err

}

// RenderFromBytes should render the template to the io.Writer.
func (render *Render) RenderFromBytes(fbytes []byte, data Map) ([]byte, error) {
	template, err := render.set.FromBytes(fbytes)
	if err != nil {
		return nil, err
	}

	var data2 pongo2.Context
	if data == nil {
		data2 = render.tplContext

	} else {
		data2 = pongo2.Context(data)
		for k, v := range render.tplContext {
			if _, ok := data2[k]; !ok {
				data2[k] = v
			}
		}
	}

	var b bytes.Buffer
	err = template.ExecuteWriter(data2, &b)
	return b.Bytes(), err
}

func (render *Render) fromCache(fname string, data Map, withInfo bool) ([]byte, os.FileInfo, error) {
	// Get file content from the file system cache
	f, err := render.openCacheFile(fname)
	if err != nil {
		render.Lock()
		_, has := render.tplCache[fname]
		if has {
			delete(render.tplCache, fname)
		}
		render.Unlock()
		return nil, nil, errors.New(fname + "is not find.")
	}
	defer f.Close()
	fileInfo, err := f.Stat()
	if err != nil {
		return nil, nil, err
	}

	if fileInfo.IsDir() {
		return nil, fileInfo, nil
	}

	render.RLock()
	tplObj, has := render.tplCache[fname]
	render.RUnlock()

	var tpl *pongo2.Template

	// When the template cache exists and the file is not updated
	if has && tplObj.modTime.Equal(fileInfo.ModTime()) {
		tpl = tplObj.template

	} else {
		// The cache template does not exist or the file is updated
		render.Lock()
		defer render.Unlock()

		// Create a new template and cache it
		fbytes, _ := ioutil.ReadAll(f)
		tpl, err = render.set.FromBytesWithName(fname, fbytes)
		if err != nil {
			return nil, nil, err
		}

		render.tplCache[fname] = &Tpl{template: tpl, modTime: fileInfo.ModTime()}
	}

	var data2 pongo2.Context
	if data == nil {
		data2 = render.tplContext
	} else {
		data2 = pongo2.Context(data)
		for k, v := range render.tplContext {
			if _, ok := data2[k]; !ok {
				data2[k] = v
			}
		}
	}
	var b bytes.Buffer
	err = tpl.ExecuteWriter(data2, &b)
	if withInfo {
		return b.Bytes(), newNowFileInfo(fileInfo, int64(b.Len())), err
	}
	return b.Bytes(), nil, err
}

type nowFileInfo struct {
	os.FileInfo
	size    int64
	modTime time.Time
}

// Size returns the size in bytes for regular files; system-dependent for others
func (info *nowFileInfo) Size() int64 {
	return info.size
}

// Mode returns file mode bits
func (info *nowFileInfo) ModTime() time.Time {
	return info.modTime
}

func newNowFileInfo(info os.FileInfo, size int64) os.FileInfo {
	return &nowFileInfo{
		FileInfo: info,
		size:     size,
		modTime:  time.Now(),
	}
}

func (render *Render) renderForFS(filename string, data Map) ([]byte, os.FileInfo, error) {
	if render.caching {
		return render.fromCache(filename, data, true)
	}
	f, err := os.Open(filename)
	if err != nil {
		return nil, nil, err
	}

	fileInfo, err := f.Stat()
	if err != nil {
		return nil, nil, err
	}

	if fileInfo.IsDir() {
		return nil, fileInfo, nil
	}

	var fbytes []byte
	fbytes, err = ioutil.ReadAll(f)
	if err != nil {
		return nil, nil, err
	}
	fbytes, err = render.RenderFromBytes(fbytes, data)
	return fbytes, newNowFileInfo(fileInfo, int64(len(fbytes))), err
}
