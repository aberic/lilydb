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

import "testing"

func newIndex() *Index {
	return NewIndex("indexID", "lily_do_not_repeat_auto_id", true)
}

func linkFit() *link {
	link := &link{md516Key: "md516"}
	link.Fit(1, 2, 3)
	return link
}

func TestLink_Fit(t *testing.T) {
	link := &link{md516Key: "md516"}
	link.Fit(1, 2, 3)
	t.Log(link)
}

func TestLink_MD516Key(t *testing.T) {
	t.Log(linkFit().MD516Key())
}

func TestLink_SeekStartIndex(t *testing.T) {
	t.Log(linkFit().SeekStartIndex())
}

func TestLink_SeekStart(t *testing.T) {
	t.Log(linkFit().SeekStart())
}

func TestLink_SeekLast(t *testing.T) {
	t.Log(linkFit().SeekLast())
}

func TestNewIndex(t *testing.T) {
	t.Log(NewIndex("indexID", "id", true))
}

func TestIndex_ID(t *testing.T) {
	idx := newIndex()
	t.Log(idx.ID())
}

func TestIndex_KeyStructure(t *testing.T) {
	idx := newIndex()
	t.Log(idx.KeyStructure())
}

func TestIndex_Primary(t *testing.T) {
	idx := newIndex()
	t.Log(idx.Primary())
}

func TestIndex_Put(t *testing.T) {
	idx := newIndex()
	link, exist := idx.Put("md516", 1)
	t.Log(link, exist)
}

func TestIndex_Get(t *testing.T) {
	idx := newIndex()
	_, _ = idx.Put("md516", 1)
	link := idx.Get("md516", 1)
	t.Log(link)
}

func TestIndex_GetFail(t *testing.T) {
	idx := newIndex()
	_, _ = idx.Put("md516", 1)
	link := idx.Get("md516", 2)
	t.Log(link)
}

func TestIndex_Recover(t *testing.T) {
	idx := newIndex()
	t.Skip(idx.Recover("database", "form"))
}

func TestIndex_RecoverLenMatchFail(t *testing.T) {
	idx := newIndex()
	t.Skip(idx.Recover("database", "form1"))
}

func TestIndex_RecoverDataEmptyFail(t *testing.T) {
	idx := newIndex()
	t.Skip(idx.Recover("database", "form2"))
}

func TestIndex_RecoverNoFile(t *testing.T) {
	idx := newIndex()
	t.Skip(idx.Recover("database1", "form"))
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

func TestNewSelector(t *testing.T) {
	var indexes = []*Index{newIndex()}
	selector, err := NewSelector([]byte(selectorJSONString), indexes, "database", "form", false)
	t.Log(selector)
	t.Log(err)
}

func TestNewSelectorFail(t *testing.T) {
	var indexes = []*Index{newIndex()}
	selector, err := NewSelector([]byte(selectorErrorJSONString), indexes, "database", "form", false)
	t.Log(selector)
	t.Log(err)
}

func TestSelector_Run(t *testing.T) {
	var indexes = []*Index{newIndex()}
	selector, err := NewSelector([]byte(selectorJSONString), indexes, "database", "form", false)
	if nil != err {
		t.Error(err)
	}
	count, values := selector.Run()
	t.Log(count, values)
}
