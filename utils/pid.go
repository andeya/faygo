package utils

import (
	"fmt"
	"os"
	"path/filepath"
)

func WritePid(pidFilename string) {
	abs, err := filepath.Abs(pidFilename)
	if err != nil {
		panic(err)
	}
	dir := filepath.Dir(abs)
	os.MkdirAll(dir, 0777)
	pid := os.Getpid()
	f, err := os.OpenFile(abs, os.O_CREATE|os.O_RDWR|os.O_TRUNC, 0666)
	if err != nil {
		panic(err)
	}
	defer f.Close()
	f.WriteString(fmt.Sprintf("%d\n", pid))
}
