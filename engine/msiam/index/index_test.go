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

package index

import (
	"github.com/aberic/gnomon"
	"strconv"
	"testing"
)

func newIndex(databaseID, formID string) *Index {
	return NewIndex(databaseID, formID, "indexID", "indexID", true)
}

func linkFit() *Link {
	link := &Link{md516Key: "md516"}
	link.Fit("key", "md516keyMd516key", 3, 0)
	return link
}

func TestLink_Fit(t *testing.T) {
	link := &Link{md516Key: "md516"}
	link.Fit("key", "md516keyMd516key", 3, 0)
	t.Log(link)
}

func TestLink_Key(t *testing.T) {
	t.Log(linkFit().Key())
}

func TestLink_MD516Key(t *testing.T) {
	t.Log(linkFit().MD516Key())
}

func TestLink_Value(t *testing.T) {
	t.Log(linkFit().Value())
}

func TestLink_Version(t *testing.T) {
	t.Log(linkFit().Version())
}

func TestNewIndex(t *testing.T) {
	t.Log(NewIndex("database", "form", "indexID", "id", true))
}

func TestIndex_ID(t *testing.T) {
	idx := newIndex("database", "form")
	t.Log(idx.ID())
}

func TestIndex_KeyStructure(t *testing.T) {
	idx := newIndex("database", "form")
	t.Log(idx.KeyStructure())
}

func TestIndex_Primary(t *testing.T) {
	idx := newIndex("database", "form")
	t.Log(idx.Primary())
}

func TestIndex_Put(t *testing.T) {
	idx := newIndex("database", "form")
	link, exist, versionGT := idx.Put("key", "md516keyMd516key", 3, 1, 0)
	t.Log(link, exist, versionGT)
}

func TestIndex_Get(t *testing.T) {
	idx := newIndex("database", "form")
	_, _, _ = idx.Put("key", "md516keyMd516key", 3, 1, 0)
	link := idx.Get("md516keyMd516key", 3)
	t.Log(link)
}

func TestIndex_GetFail(t *testing.T) {
	idx := newIndex("database", "form")
	link := idx.Get("md516", 1)
	t.Log(link)
}

var selectorJSONString = `{
		"Conditions":[
			{
				"Param":"Age",
				"Cond":"gt",
				"Value":3
			}
		],
		"Sort":{
			"Param":"Age",
			"Asc":false
		},
		"Skip":1,
		"Limit":3
    }`

var selectorErrorJSONString = `error`

//var selector = &Selector{
//	Conditions: []*condition{
//		{
//			Param: "Age",
//			Cond:  "gt",
//			Value: 0,
//		},
//	},
//	Skip: 1,
//	Sort: &rank{
//		Param: "Age",
//		ASC:   false,
//	},
//	Limit: 3,
//}

func TestNewSelector(t *testing.T) {
	var indexes = []*Index{newIndex("database", "form")}
	selector, err := NewSelector([]byte(selectorJSONString), indexes, "database", "form", false)
	t.Log(selector)
	t.Log(err)
}

func TestNewSelectorFail(t *testing.T) {
	var indexes = []*Index{newIndex("database", "form")}
	selector, err := NewSelector([]byte(selectorErrorJSONString), indexes, "database", "form", false)
	t.Log(selector)
	t.Log(err)
}

func TestSelector_Run(t *testing.T) {
	var (
		idx     = newIndex("database", "form")
		indexes = []*Index{idx}
	)
	type Value struct {
		Name string
		Age  int
	}
	var i uint64
	for i = 0; i < 10; i++ {
		v := &Value{Name: "name", Age: int(i)}
		_, _, _ = idx.Put("key", gnomon.HashMD516(strconv.Itoa(int(i))), i, v, 0)
	}
	selector, err := NewSelector([]byte(selectorJSONString), indexes, "database", "form", false)
	if nil != err {
		t.Error(err)
	}
	count, values := selector.Run()
	t.Log(count, values)
	for _, v := range values {
		t.Log(v)
	}
}
