package logging

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"testing"
	"time"
)

func TestFile1(t *testing.T) {
	log := NewLogger(10000)
	log.AddAdapter("file", `{"filename":"test.log"}`)
	log.Debug("debug")
	log.Info("info")
	log.Notice("notice")
	log.Warn("warning")
	log.Error("error")
	log.Alert("alert")
	log.Critical("critical")
	log.Emergency("emergency")
	f, err := os.Open("test.log")
	if err != nil {
		t.Fatal(err)
	}
	b := bufio.NewReader(f)
	lineNum := 0
	for {
		line, _, err := b.ReadLine()
		if err != nil {
			break
		}
		if len(line) > 0 {
			lineNum++
		}
	}
	var expected = LevelDebug + 1
	if lineNum != expected {
		t.Fatal(lineNum, "not "+strconv.Itoa(expected)+" lines")
	}
	os.Remove("test.log")
}

func TestFile2(t *testing.T) {
	log := NewLogger(10000)
	log.AddAdapter("file", fmt.Sprintf(`{"filename":"test2.log","level":%d}`, LevelError))
	log.Debug("debug")
	log.Info("info")
	log.Notice("notice")
	log.Warn("warning")
	log.Error("error")
	log.Alert("alert")
	log.Critical("critical")
	log.Emergency("emergency")
	f, err := os.Open("test2.log")
	if err != nil {
		t.Fatal(err)
	}
	b := bufio.NewReader(f)
	lineNum := 0
	for {
		line, _, err := b.ReadLine()
		if err != nil {
			break
		}
		if len(line) > 0 {
			lineNum++
		}
	}
	var expected = LevelError + 1
	if lineNum != expected {
		t.Fatal(lineNum, "not "+strconv.Itoa(expected)+" lines")
	}
	os.Remove("test2.log")
}

func TestFileRotate(t *testing.T) {
	log := NewLogger(10000)
	log.AddAdapter("file", `{"filename":"test3.log","maxlines":4}`)
	log.Debug("debug")
	log.Info("info")
	log.Notice("notice")
	log.Warn("warning")
	log.Error("error")
	log.Alert("alert")
	log.Critical("critical")
	log.Emergency("emergency")
	rotateName := "test3" + fmt.Sprintf(".%s.%03d", time.Now().Format("2006-01-02"), 1) + ".log"
	b, err := exists(rotateName)
	if !b || err != nil {
		os.Remove("test3.log")
		t.Fatal("rotate not generated")
	}
	os.Remove(rotateName)
	os.Remove("test3.log")
}

func exists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}

func BenchmarkFile(b *testing.B) {
	log := NewLogger(100000)
	log.AddAdapter("file", `{"filename":"test4.log"}`)
	for i := 0; i < b.N; i++ {
		log.Debug("debug")
	}
	os.Remove("test4.log")
}

func BenchmarkFileAsynchronous(b *testing.B) {
	log := NewLogger(100000)
	log.AddAdapter("file", `{"filename":"test4.log"}`)
	// log.Async()
	for i := 0; i < b.N; i++ {
		log.Debug("debug")
	}
	os.Remove("test4.log")
}

func BenchmarkFileCallDepth(b *testing.B) {
	log := NewLogger(100000)
	log.AddAdapter("file", `{"filename":"test4.log"}`)
	log.EnableFuncCallDepth(true)
	log.SetLogFuncCallDepth(2)
	for i := 0; i < b.N; i++ {
		log.Debug("debug")
	}
	os.Remove("test4.log")
}

func BenchmarkFileAsynchronousCallDepth(b *testing.B) {
	log := NewLogger(100000)
	log.AddAdapter("file", `{"filename":"test4.log"}`)
	log.EnableFuncCallDepth(true)
	log.SetLogFuncCallDepth(2)
	// log.Async()
	for i := 0; i < b.N; i++ {
		log.Debug("debug")
	}
	os.Remove("test4.log")
}

func BenchmarkFileOnGoroutine(b *testing.B) {
	log := NewLogger(100000)
	log.AddAdapter("file", `{"filename":"test4.log"}`)
	for i := 0; i < b.N; i++ {
		go log.Debug("debug")
	}
	os.Remove("test4.log")
}
