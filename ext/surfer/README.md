# Surfer    [![GoDoc](https://godoc.org/github.com/tsuna/gohbase?status.png)](https://godoc.org/github.com/andeya/surfer) [![GitHub release](https://img.shields.io/github/release/andeya/surfer.svg)](https://github.com/andeya/surfer/releases)

Package surfer is a high level concurrency http client.
It has `surf` and` phantom` download engines, highly simulated browser behavior, the function of analog login and so on.

[简体中文](https://github.com/andeya/surfer/blob/master/README_ZH.md)

## Features
- Both `surf` and `phantomjs` engines are supported
- Support random User-Agent
- Support cache cookie
- Support http/https

## Usage
```
package main

import (
    "github.com/andeya/surfer"
    "io/ioutil"
    "log"
)

func main() {
    // Use surf engine
    resp, err := surfer.Download(&surfer.Request{
        Url: "http://github.com/andeya/surfer",
    })
    if err != nil {
        log.Fatal(err)
    }
    b, err := ioutil.ReadAll(resp.Body)
    log.Println(string(b), err)

    // Use phantomjs engine
    resp, err = surfer.Download(&surfer.Request{
        Url:          "http://github.com/andeya",
        DownloaderID: 1,
    })
    if err != nil {
        log.Fatal(err)
    }
    b, err = ioutil.ReadAll(resp.Body)
    log.Println(string(b), err)

    resp.Body.Close()
    surfer.DestroyJsFiles()
}
```
[Full example](https://github.com/andeya/faygo/raw/master/samples)

## License

Surfer is under Apache v2 License. See the [LICENSE](https://github.com/andeya/faygo/raw/master/LICENSE) file for the full license text
