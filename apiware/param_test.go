// Copyright 2016 HenryLee. All Rights Reserved.
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

package apiware

import (
	"reflect"
	"testing"
)

func TestParseTags(t *testing.T) {
	var a reflect.StructTag = `param1:"<in:\"path\"> <name:test> <desc:test\\<1,2\\>> <required> <range::4>"`
	var param1 = a.Get("param1")
	var values1 = ParseTags(param1)
	t.Logf("values1:%#v", values1)

	var b reflect.StructTag = `  param2:"   <in:\"path\"> <name : test   > <desc:test\\<1,2\\>> <required:>    <range: :4  >  "   `
	var param2 = b.Get("param2")
	var values2 = ParseTags(param2)
	t.Logf("values2:%#v", values2)

	if values1["in"] != values2["in"] ||
		values1["desc"] != values2["desc"] ||
		values1["required"] != values2["required"] ||
		values1["range"] != values2["range"] ||
		values1["name"] != values2["name"] {

		t.Fail()
	}
}

func TestParseTags2(t *testing.T) {
	var a reflect.StructTag = `param:"<in:query> <name:p> <len: 1:10> <regexp: ^[\\w]*$>"`
	var param = a.Get("param")
	var values = ParseTags(param)
	t.Logf("values:%#v", values)
}
