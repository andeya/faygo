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
	log := NewLogger("TestFile1")
	fileBackend, err := NewDefaultFileBackend("test.log")
	if err != nil {
		panic(err)
	}
	var format = MustStringFormatter(
		`%{color}%{module} %{time:15:04:05.000} %{longfile} ▶ %{level:.4s} %{id:03x}%{color:reset} %{message}`,
	)
	log.SetBackend(AddModuleLevel(NewBackendFormatter(fileBackend, format)))
	log.Debug("debug\n")
	log.Info("info\n")
	log.Notice("notice\n")
	log.Warning("warning\n")
	log.Error("error\n")
	log.Critical("critical\n")
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
	f.Close()
	fileBackend.Close()
	var expected = int(DEBUG) + 1
	if lineNum != expected {
		t.Fatal(lineNum, "not "+strconv.Itoa(expected)+" lines")
	}
	os.Remove("test.log")
}

func TestFile2(t *testing.T) {
	log := NewLogger("TestFile2")
	fileBackend, err := NewDefaultFileBackend("test2.log", 1000)
	if err != nil {
		panic(err)
	}
	var format = MustStringFormatter(
		`%{color}%{module} %{time:15:04:05.000} %{longfile} ▶ %{level:.4s} %{id:03x}%{color:reset} %{message}`,
	)
	log.SetBackend(AddModuleLevel(NewBackendFormatter(fileBackend, format)))
	log.Debug("debug\n")
	log.Info("info\n")
	log.Notice("notice\n")
	log.Warning("warning\n")
	log.Error("error\n")
	log.Critical("critical\n")
	log.Close()
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
	f.Close()
	var expected = int(DEBUG) + 1
	if lineNum != expected {
		t.Fatal(lineNum, "not "+strconv.Itoa(expected)+" lines")
	}
	os.Remove("test2.log")
}

func TestFileRotate(t *testing.T) {
	log := NewLogger("TestFileRotate")
	fileBackend, err := NewDefaultFileBackend("test3.log")
	if err != nil {
		panic(err)
	}
	fileBackend.MaxLines = 4
	log.SetBackend(AddModuleLevel(fileBackend))
	log.Debug("debug")
	log.Info("info")
	log.Notice("notice")
	log.Warn("warning")
	log.Error("error")
	log.Critical("critical")
	rotateName := "test3" + fmt.Sprintf(".%s.%03d", time.Now().Format("2006-01-02"), 1) + ".log"
	b, err := exists(rotateName)
	if !b || err != nil {
		os.Remove("test3.log")
		t.Fatal("rotate not generated")
	}
	fileBackend.Close()
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
	log := NewLogger("BenchmarkFile")
	fileBackend, err := NewDefaultFileBackend("test4.log")
	if err != nil {
		panic(err)
	}
	log.SetBackend(AddModuleLevel(fileBackend))
	for i := 0; i < b.N; i++ {
		log.Debug("debug")
	}
	fileBackend.Close()
	os.Remove("test4.log")
}

func BenchmarkFileAsynchronous(b *testing.B) {
	log := NewLogger("BenchmarkFileAsynchronous")
	fileBackend, err := NewDefaultFileBackend("test4.log")
	if err != nil {
		panic(err)
	}
	log.SetBackend(AddModuleLevel(fileBackend))
	for i := 0; i < b.N; i++ {
		log.Debug("debug")
	}
	fileBackend.Close()
	os.Remove("test4.log")
}

func BenchmarkFileCallDepth(b *testing.B) {
	log := NewLogger("BenchmarkFileCallDepth")
	fileBackend, err := NewDefaultFileBackend("test4.log")
	if err != nil {
		panic(err)
	}
	log.SetBackend(AddModuleLevel(fileBackend))
	log.ExtraCalldepth = 2
	for i := 0; i < b.N; i++ {
		log.Debug("debug")
	}
	fileBackend.Close()
	os.Remove("test4.log")
}

func BenchmarkFileAsynchronousCallDepth(b *testing.B) {
	log := NewLogger("BenchmarkFileAsynchronousCallDepth")
	fileBackend, err := NewDefaultFileBackend("test4.log")
	if err != nil {
		panic(err)
	}
	log.SetBackend(AddModuleLevel(fileBackend))
	log.ExtraCalldepth = 2
	for i := 0; i < b.N; i++ {
		log.Debug("debug")
	}
	fileBackend.Close()
	os.Remove("test4.log")
}

func BenchmarkFileOnGoroutine(b *testing.B) {
	log := NewLogger("BenchmarkFileOnGoroutine")
	fileBackend, err := NewDefaultFileBackend("test4.log")
	if err != nil {
		panic(err)
	}
	log.SetBackend(AddModuleLevel(fileBackend))
	for i := 0; i < b.N; i++ {
		go log.Debug("debug")
	}
	fileBackend.Close()
	os.Remove("test4.log")
}
