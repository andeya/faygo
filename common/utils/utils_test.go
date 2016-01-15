package utils

import (
	"testing"
)

// 遍历并返回指定类型范围的文件名列表
// 默认返回所有文件
func TestWalkFiles(t *testing.T) {
	var path, suffixes = "./", []string{}
	t.Logf("%#v", WalkFiles(path, suffixes...))
}

// 遍历并返回目录列表
func TestWalkDir(t *testing.T) {
	var path = "./"
	t.Logf("%#v", WalkDir(path))
}

func TestRelPath(t *testing.T) {
	var targpath = "E:\\__HENRY__\\Program\\Go\\src\\github.com\\henrylee2cn\\morechat\\utils\\pool\\p"
	t.Log(RelPath(targpath))
}
func TestSnakeString(t *testing.T) {
	t.Log(SnakeString("/MaxMin"))
	t.Log(CamelString("_/MaxMin"))
}
