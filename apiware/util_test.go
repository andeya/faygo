package apiware

import (
	"testing"
)

func TestInterfaceToSnake(t *testing.T) {
	type SampleModel struct{}
	name := interfaceToSnake(&SampleModel{})
	if name != "sample_model" {
		t.Fatal("wrong table name", name)
	}
	name = interfaceToSnake(SampleModel{})
	if name != "sample_model" {
		t.Fatal("wrong table name", name)
	}
}

func TestSnakeToUpperCamel(t *testing.T) {
	if s := snakeToUpperCamel("table_name"); s != "TableName" {
		t.Fatal("wrong string", s)
	}
}
