package main

import (
	"github.com/henrylee2cn/thinkgo/ext/surfer"
	"io/ioutil"
	"log"
	"net/url"
	"time"
)

func main() {
	var values, _ = url.ParseQuery("username=123456@qq.com&password=123456&login_btn=login_btn&submit=login_btn")
	log.Println("values:", values)
	var files = []surfer.PostFile{
		{
			Fieldname: "abc",
			Filename:  "filename.txt",
			Bytes:     []byte("files test."),
		},
	}

	// 默认使用surf内核下载
	log.Println("********************************************* surf内核GET下载测试开始 *********************************************")
	resp, err := surfer.Download(&surfer.Request{
		Url: "http://www.baidu.com/",
	})
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("baidu resp.Status: %s\nresp.Header: %#v\n", resp.Status, resp.Header)

	b, err := ioutil.ReadAll(resp.Body)
	resp.Body.Close()
	log.Printf("baidu resp.Body: %s\nerr: %v", b, err)

	// 默认使用surf内核下载
	log.Println("********************************************* surf内核POST下载测试开始 *********************************************")
	resp, err = surfer.Download(&surfer.Request{
		Url:    "http://accounts.lewaos.com/",
		Method: "POST",
		Values: values,
		Files:  files,
	})
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("lewaos resp.Status: %s\nresp.Header: %#v\n", resp.Status, resp.Header)

	b, err = ioutil.ReadAll(resp.Body)
	resp.Body.Close()
	log.Printf("lewaos resp.Body: %s\nerr: %v", b, err)

	log.Println("********************************************* phantomjs内核GET下载测试开始 *********************************************")

	// 指定使用phantomjs内核下载
	resp, err = surfer.Download(&surfer.Request{
		Url:          "http://www.baidu.com/",
		DownloaderID: 1,
	})
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("baidu resp.Status: %s\nresp.Header: %#v\n", resp.Status, resp.Header)

	b, err = ioutil.ReadAll(resp.Body)
	resp.Body.Close()
	log.Printf("baidu resp.Body: %s\nerr: %v", b, err)

	log.Println("********************************************* phantomjs内核POST下载测试开始 *********************************************")

	// 指定使用phantomjs内核下载
	resp, err = surfer.Download(&surfer.Request{
		DownloaderID: 1,
		Url:          "http://accounts.lewaos.com/",
		Method:       "POST",
		Values:       values,
	})
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("lewaos resp.Status: %s\nresp.Header: %#v\n", resp.Status, resp.Header)

	b, err = ioutil.ReadAll(resp.Body)
	resp.Body.Close()
	log.Printf("lewaos resp.Body: %s\nerr: %v", b, err)

	surfer.DestroyJsFiles()

	time.Sleep(10e9)
}
