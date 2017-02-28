/*
Package Faygo is a fast and concise Go Web framework that can be used to develop high-performance web app(especially API) with fewer codes;
Just define a struct Handler, Faygo will automatically bind/verify the request parameters and generate the online API doc.

Copyright 2016 HenryLee. All Rights Reserved.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

     http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.

A trivial example is:

    package main

    import (
        // "mime/multipart"
        "time"
        "github.com/henrylee2cn/faygo"
    )

    type Index struct {
        Id        int      `param:"<in:path> <required> <desc:ID> <range: 0:10>"`
        Title     string   `param:"<in:query> <nonzero>"`
        Paragraph []string `param:"<in:query> <name:p> <len: 1:10> <regexp: ^[\\w]*$>"`
        Cookie    string   `param:"<in:cookie> <name:faygoID>"`
        // Picture         *multipart.FileHeader `param:"<in:formData> <name:pic> <maxmb:30>"`
    }

    func (i *Index) Serve(ctx *faygo.Context) error {
        if ctx.CookieParam("faygoID") == "" {
            ctx.SetCookie("faygoID", time.Now().String())
        }
        return ctx.JSON(200, i)
    }

    func main() {
        app := faygo.New("myapp", "0.1")

        // Register the route in a chain style
        app.GET("/index/:id", new(Index))

        // Register the route in a tree style
        // app.Route(
        //     app.NewGET("/index/:id", new(Index)),
        // )

        // Start the service
        faygo.Run()
    }

run result:
    http GET:
    http://localhost:8080/index/1?title=test&p=abc&p=xyz
    response:
        {
            "Id": 1,
            "Title": "test",
            "Paragraph": [
                "abc",
                "xyz"
            ],
            "Cookie": "2016-11-13 01:14:40.9038005 +0800 CST"
        }


StructHandler tag value description:
    tag   |   key    | required |     value     |   desc
    ------|----------|----------|---------------|----------------------------------
    param |    in    | only one |     path      | (position of param) if `required` is unsetted, auto set it. e.g. url: "http://www.abc.com/a/{path}"
    param |    in    | only one |     query     | (position of param) e.g. url: "http://www.abc.com/a?b={query}"
    param |    in    | only one |     formData  | (position of param) e.g. "request body: a=123&b={formData}"
    param |    in    | only one |     body      | (position of param) request body can be any content
    param |    in    | only one |     header    | (position of param) request header info
    param |    in    | only one |     cookie    | (position of param) request cookie info, support: `*http.Cookie`,`http.Cookie`,`string`,`[]byte`
    param |   name   |    no    |   (e.g.`id`)   | specify request param`s name
    param | required |    no    |               | request param is required
    param |   desc   |    no    |   (e.g.`id`)   | request param description
    param |   len    |    no    | (e.g.`3:6` `3`) | length range of param's value
    param |   range  |    no    |  (e.g.`0:10`)  | numerical range of param's value
    param |  nonzero |    no    |               | param`s value can not be zero
    param |   maxmb  |    no    |   (e.g.`32`)   | when request Content-Type is multipart/form-data, the max memory for body.(multi-param, whichever is greater)
    param |  regexp  |    no    | (e.g.`^\\w+$`) | verify the value of the param with a regular expression(param value can not be null)
    param |   err    |    no    |(e.g.`incorrect password format`)| the custom error for binding or validating

    NOTES:
        1. the binding object must be a struct pointer
        2. in addition to `*multipart.FileHeader`, the binding struct's field can not be a pointer
        3. `regexp` or `param` tag is only usable when `param:"type(xxx)"` is exist
        4. if the `param` tag is not exist, anonymous field will be parsed
        5. when the param's position(`in`) is `formData` and the field's type is `*multipart.FileHeader`, `multipart.FileHeader`, `[]*multipart.FileHeader` or `[]multipart.FileHeader`, the param receives file uploaded
        6. if param's position(`in`) is `cookie`, field's type must be `*http.Cookie` or `http.Cookie`
        7. param tags `in(formData)` and `in(body)` can not exist at the same time
        8. there should not be more than one `in(body)` param tag

List of supported param value types:
    base    |   slice    | special
    --------|------------|-------------------------------------------------------
    string  |  []string  | [][]byte
    byte    |  []byte    | [][]uint8
    uint8   |  []uint8   | *multipart.FileHeader (only for `formData` param)
    bool    |  []bool    | []*multipart.FileHeader (only for `formData` param)
    int     |  []int     | *http.Cookie (only for `net/http`'s `cookie` param)
    int8    |  []int8    | http.Cookie (only for `net/http`'s `cookie` param)
    int16   |  []int16   | struct (struct type only for `body` param or as an anonymous field to extend params)
    int32   |  []int32   |
    int64   |  []int64   |
    uint8   |  []uint8   |
    uint16  |  []uint16  |
    uint32  |  []uint32  |
    uint64  |  []uint64  |
    float32 |  []float32 |
    float64 |  []float64 |
*/
package faygo
