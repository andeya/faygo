package apiware

import (
	"reflect"
	"strings"
	"testing"
	"time"
)

func TestParsetags(t *testing.T) {
	m := ParseTags(`<in:path> <required> <desc:banana>`)
	if x, ok := m[KEY_REQUIRED]; !ok {
		t.Fatal("wrong value", ok, x)
	}
	if x, ok := m[KEY_DESC]; !ok || x != "banana" {
		t.Fatal("wrong value", x)
	}
}

func TestFieldIsZero(t *testing.T) {
	if !isZero(reflect.ValueOf(0)) {
		t.Fatal("should be zero")
	}
	if !isZero(reflect.ValueOf("")) {
		t.Fatal("should be zero")
	}
	if !isZero(reflect.ValueOf(false)) {
		t.Fatal("should be zero")
	}
	if isZero(reflect.ValueOf(true)) {
		t.Fatal("should not be zero")
	}
	if isZero(reflect.ValueOf(-1)) {
		t.Fatal("should not be zero")
	}
	if isZero(reflect.ValueOf(1)) {
		t.Fatal("should not be zero")
	}
	if isZero(reflect.ValueOf("asdf")) {
		t.Fatal("should not be zero")
	}
}

func TestFieldvalidate(t *testing.T) {
	type Schema struct {
		A string  `param:"<in:path> <len: 3:6> <name:p> <err:This is a custom error!>"`
		B float32 `param:"<in:query> <range: 10:20>"`
		C string  `param:"<in:query> <len: :4> <nonzero>"`
		D string  `param:"<in:query> <regexp: ^[a-zA-Z0-9_.+-]+@[a-zA-Z0-9-]+\\.[a-zA-Z0-9-.]+$>"`
	}
	m, _ := NewParamsAPI(&Schema{B: 9.999999}, nil, nil)
	a := m.params[0]
	if x := len(a.tags); x != 5 {
		t.Fatal("wrong len", x, a.tags)
	}
	if x, ok := a.tags[KEY_LEN]; !ok || x != "3:6" {
		t.Fatal("wrong value", x, ok)
	}
	if err := a.validate(a.rawValue); err == nil || err.Error() != "This is a custom error!" {
		t.Fatal("should not validate")
	}
	if err := a.validate(reflect.ValueOf("abc")); err != nil {
		t.Fatal("should validate", err)
	}
	if err := a.validate(reflect.ValueOf("abcdefg")); err == nil || err.Error() != "This is a custom error!" {
		t.Fatal("should not validate")
	}

	b := m.params[1]
	if x := len(b.tags); x != 2 {
		t.Fatal("wrong len", x)
	}
	if err := b.validate(b.rawValue); err == nil || !strings.Contains(err.Error(), "small") {
		t.Fatal("should not validate")
	}
	if err := b.validate(reflect.ValueOf(10)); err != nil {
		t.Fatal("should validate", err)
	}
	if err := b.validate(reflect.ValueOf(21)); err == nil || !strings.Contains(err.Error(), "big") {
		t.Fatal("should not validate")
	}

	c := m.params[2]
	if x := len(c.tags); x != 3 {
		t.Fatal("wrong len", x)
	}
	if err := c.validate(c.rawValue); err == nil || !strings.Contains(err.Error(), "not set") {
		t.Fatal("should not validate")
	}
	if err := c.validate(reflect.ValueOf("a")); err != nil {
		t.Fatal("should validate", err)
	}
	if err := c.validate(reflect.ValueOf("abcde")); err == nil || !strings.Contains(err.Error(), "long") {
		t.Fatal("should not validate")
	}

	d := m.params[3]
	if x := len(d.tags); x != 2 {
		t.Fatal("wrong len", x)
	}
	if err := d.validate(reflect.ValueOf("gggg@gmail.com")); err != nil {
		t.Fatal("should validate", err)
	}
	if err := d.validate(reflect.ValueOf("www.google.com")); err == nil || !strings.Contains(err.Error(), "not match") {
		t.Fatal("should not validate", err)
	}
}

func TestFieldOmit(t *testing.T) {
	type schema struct {
		A string `param:"-"`
		B string
	}
	m, _ := NewParamsAPI(&schema{}, nil, nil)
	if x := len(m.params); x != 0 {
		t.Fatal("wrong len", x)
	}
}

func TestInterfaceNewParamsAPIWithEmbedded(t *testing.T) {
	type third struct {
		Num int64 `param:"<in:query>"`
	}
	type embed struct {
		Name  string `param:"<in:query>"`
		Value string `param:"<in:query>"`
		third
	}
	type table struct {
		ColPrimary int64 `param:"<in:query>"`
		embed
	}
	table1 := &table{
		6, embed{"Mrs. A", "infinite", third{Num: 12345}},
	}
	m, err := NewParamsAPI(table1, nil, nil)
	if err != nil {
		t.Fatal("error not nil", err)
	}
	f := m.params[1]
	if x, ok := toString(f.rawValue); !ok || x != "Mrs. A" {
		t.Fatal("wrong value from embedded struct")
	}
	f = m.params[3]
	if x, _ := f.Raw().(int64); x != 12345 {
		t.Fatal("wrong value from third struct")
	}
}

type indexedTable struct {
	ColIsRequired string `param:"<in:query> <required>"`
	ColVarChar    string `param:"<in:query> <desc:banana>"`
	ColTime       time.Time
}

func TestInterfaceNewParamsAPI(t *testing.T) {
	now := time.Now()
	table1 := &indexedTable{
		ColVarChar: "orange",
		ColTime:    now,
	}
	m, err := NewParamsAPI(table1, nil, nil)
	if err != nil {
		t.Fatal("error not nil", err)
	}
	if x := len(m.params); x != 2 {
		t.Fatal("wrong value", x)
	}
	f := m.params[0]
	if !f.IsRequired() {
		t.Fatal("wrong value")
	}
	f = m.params[1]
	if x, ok := toString(f.rawValue); !ok || x != "orange" {
		t.Fatal("wrong value", x)
	}
	if isZero(f.rawValue) {
		t.Fatal("wrong value")
	}
	if f.Description() != "banana" {
		t.Fatal("should value", f.Description())
	}
	if f.IsRequired() {
		t.Fatal("wrong value")
	}
}

func makeWhitespaceVisible(s string) string {
	s = strings.Replace(s, "\t", "\\t", -1)
	s = strings.Replace(s, "\r\n", "\\r\\n", -1)
	s = strings.Replace(s, "\r", "\\r", -1)
	s = strings.Replace(s, "\n", "\\n", -1)
	return s
}
func isZero(v reflect.Value) bool {
	return v.Interface() == reflect.Zero(v.Type()).Interface()
}

func toString(v reflect.Value) (string, bool) {
	s, ok := v.Interface().(string)
	return s, ok
}
