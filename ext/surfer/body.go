// Copyright 2015 henrylee2cn Author. All Rights Reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package surfer

import (
	"bytes"
	"encoding/json"
	"encoding/xml"
	"io"
	"io/ioutil"
	"log"
	"mime/multipart"
	"net/url"
	"strings"
)

// body set request body
type body interface {
	SetBody(*Request) error
}

// Content bytes type of body content
type Content struct {
	ContentType string
	Bytes       []byte
}

var _ body = new(Content)

// SetBody sets request body
func (c *Content) SetBody(r *Request) error {
	r.Header.Set("Content-Type", c.ContentType)
	r.body = bytes.NewReader(c.Bytes)
	return nil
}

// Bytes bytes type of body content, without content type
type Bytes []byte

var _ body = Bytes("")

// SetBody sets request body
func (b Bytes) SetBody(r *Request) error {
	r.body = bytes.NewReader(b)
	return nil
}

type (
	// Form impletes body interface
	Form struct {
		// Values [field name]-[]value
		Values map[string][]string
		// Files [field name]-[]File
		Files map[string][]File
	}
	// File post form file
	File struct {
		Filename string
		Bytes    []byte
	}
)

var _ body = new(Form)

// SetBody sets request body
func (f Form) SetBody(r *Request) error {
	if len(f.Files) > 0 {
		pr, pw := io.Pipe()
		bodyWriter := multipart.NewWriter(pw)
		go func() {
			for fieldname, postfiles := range f.Files {
				for _, postfile := range postfiles {
					fileWriter, err := bodyWriter.CreateFormFile(fieldname, postfile.Filename)
					if err != nil {
						log.Println("[E] Surfer: Multipart:", err)
					}
					_, err = fileWriter.Write(postfile.Bytes)
					if err != nil {
						log.Println("[E] Surfer: Multipart:", err)
					}
				}
			}
			for k, v := range f.Values {
				for _, vv := range v {
					bodyWriter.WriteField(k, vv)
				}
			}
			bodyWriter.Close()
			pw.Close()
		}()
		r.Header.Set("Content-Type", bodyWriter.FormDataContentType())
		r.body = ioutil.NopCloser(pr)
	} else {
		r.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		r.body = strings.NewReader(url.Values(f.Values).Encode())
	}
	return nil
}

// JSONObj JSON type of body content
type JSONObj struct{ Data interface{} }

var _ body = new(JSONObj)

// SetBody sets request body
func (obj *JSONObj) SetBody(r *Request) error {
	r.Header.Set("Content-Type", "application/json;charset=utf-8")
	if obj == nil || obj.Data == nil {
		return nil
	}
	b, err := json.Marshal(obj.Data)
	if err != nil {
		return err
	}
	r.body = bytes.NewReader(b)
	return nil
}

// XMLObj XML type of body content
type XMLObj struct{ Data interface{} }

var _ body = new(XMLObj)

// SetBody sets request body
func (obj *XMLObj) SetBody(r *Request) error {
	r.Header.Set("Content-Type", "application/xml;charset=utf-8")
	if obj == nil || obj.Data == nil {
		return nil
	}
	b, err := xml.Marshal(obj.Data)
	if err != nil {
		return err
	}
	r.body = bytes.NewReader(b)
	return nil
}
