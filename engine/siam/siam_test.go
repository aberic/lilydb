/*
 * MIT License
 *
 * Copyright (c) 2020 aberic
 *
 * Permission is hereby granted, free of charge, to any person obtaining a copy
 * of this software and associated documentation files (the "Software"), to deal
 * in the Software without restriction, including without limitation the rights
 * to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
 * copies of the Software, and to permit persons to whom the Software is
 * furnished to do so, subject to the following conditions:
 *
 * The above copyright notice and this permission notice shall be included in all
 * copies or substantial portions of the Software.
 *
 * THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
 * IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
 * FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
 * AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
 * LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
 * OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
 * SOFTWARE.
 */

package siam

import "testing"

func form() *Form {
	return NewForm("databaseID", "formID", "formName", "comment")
}

func TestNewForm(t *testing.T) {
	t.Log(NewForm("databaseID", "formID", "formName", "comment"))
}

func TestForm_AutoID(t *testing.T) {
	t.Log(form().AutoID())
}

func TestForm_ID(t *testing.T) {
	t.Log(form().ID())
}

func TestForm_Name(t *testing.T) {
	t.Log(form().Name())
}

func TestForm_Comment(t *testing.T) {
	t.Log(form().Comment())
}

func TestForm_FormType(t *testing.T) {
	t.Log(form().FormType())
}

func TestForm_NewIndex(t *testing.T) {
	form().NewIndex("id", true)
}

type Value struct {
	Name string
	Age  int
}

var selectorJSONString = `{
		"Conditions":[
			{
				"Param":"Age",
				"Cond":"gt",
				"Value":0
			}
		],
		"Sort":{
			"Param":"Age",
			"Asc":false
		},
		"Skip":0,
		"Limit":10
    }`

func TestForm_Insert(t *testing.T) {
	for i := 0; i < 10; i++ {
		t.Log(form().Insert(&Value{Name: "name", Age: i}))
	}
}

func TestForm_Select(t *testing.T) {
	fm := form()
	for _, idx := range fm.indexes {
		if err := idx.Recover(); nil != err {
			t.Skip(err) // todo
		}
	}
	t.Log(fm.Select([]byte(selectorJSONString)))
}
