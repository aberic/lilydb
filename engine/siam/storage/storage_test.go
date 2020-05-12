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

package storage

import (
	"github.com/aberic/gnomon"
	"github.com/aberic/lilydb/engine/siam/utils"
	"strconv"
	"testing"
)

func TestObtain(t *testing.T) {
	s := Obtain()
	type Value struct {
		Name string
		Age  int
	}
	for i := 0; i < 10; i++ {
		v := &Value{Name: "name", Age: i}
		if err := s.Store("databaseID", "formID", v, []*Write{
			{
				IndexID:           "indexID",
				FormIndexFilePath: utils.PathFormIndexFile("databaseID", "formID", "indexID"),
				MD516Key:          gnomon.HashMD516(strconv.Itoa(i)),
				HashKey:           1234,
				SeekStartIndex:    -1,
				Handler: func(SeekStartIndex int64, SeekStart int64, SeekLast int) {
					t.Log(SeekStartIndex, SeekStart, SeekLast)
				},
			},
		}); nil != err {
			t.Error(err)
		}
	}
}

func TestStorage_Store(t *testing.T) {
	err := Obtain().Store("databaseID", "formID", "value", []*Write{
		{
			IndexID:           "indexID",
			FormIndexFilePath: utils.PathFormIndexFile("databaseID", "formID", "indexID"),
			MD516Key:          "key",
			HashKey:           1234,
			SeekStartIndex:    -1,
			Handler: func(SeekStartIndex int64, SeekStart int64, SeekLast int) {
				t.Log(SeekStartIndex, SeekStart, SeekLast)
			},
		},
		{
			IndexID:           "index2ID",
			FormIndexFilePath: utils.PathFormIndexFile("databaseID", "formID", "index2ID"),
			MD516Key:          "key",
			HashKey:           1234,
			SeekStartIndex:    -1,
			Handler: func(SeekStartIndex int64, SeekStart int64, SeekLast int) {
				t.Log(SeekStartIndex, SeekStart, SeekLast)
			},
		},
	})
	t.Log(err)
}

func TestStorage_StoreFailIndexFilePath(t *testing.T) {
	err := Obtain().Store("databaseID", "formID", "value", []*Write{
		{
			IndexID:           "indexID",
			FormIndexFilePath: "/usr/tmp",
			MD516Key:          "key",
			HashKey:           1234,
			SeekStartIndex:    -1,
			Handler: func(SeekStartIndex int64, SeekStart int64, SeekLast int) {
				t.Log(SeekStartIndex, SeekStart, SeekLast)
			},
		},
	})
	t.Skip(err)
}

func TestStorage_Take(t *testing.T) {
	r, err := Obtain().Take(utils.PathFormFile("databaseID", "formID"), 0, 42)
	t.Log(err, r)
}

func TestStorage_TakeFail(t *testing.T) {
	r, err := Obtain().Take(utils.PathFormFile("databaseID", "formID"), 0, 1)
	if nil != err {
		t.Skip(err)
	} else {
		t.Log(r)
	}
}
