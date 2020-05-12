/*
 * MIT License
 *
 * Copyright (c) 2020. aberic
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
	"os"
	"sync"
)

// index 索引操作
type index struct {
	path string // path 存储文件路径
	file *os.File
	mu   sync.RWMutex
}

// form 表级操作
type form struct {
	path    string // path 存储文件路径
	file    *os.File
	indexes map[string]*index
	mu      sync.RWMutex
}

// database 表依赖链集合，表ID=表级操作对象
type database struct {
	forms map[string]*form
	mu    sync.RWMutex
}

// engine 数据库依赖链集合，库ID=表依赖链集合
type engine struct {
	databases map[string]*database
	mu        sync.RWMutex
}

func (e *engine) mkDatabase(databaseID string) {
	defer e.mu.Unlock()
	e.mu.Lock()
	if nil != e.databases[databaseID] {
		return
	}
	db := &database{forms: map[string]*form{}}
	e.databases[databaseID] = db
}

func (e *engine) mkFormTry(databaseID, formID, path string) {
	db := e.databases[databaseID]
	defer db.mu.Unlock()
	db.mu.Lock()
	if nil != db.forms[formID] {
		return
	}
	db.forms[formID] = &form{path: path, indexes: map[string]*index{}}
}

func (e *engine) mkIndexTry(databaseID, formID, indexID, path string) {
	fm := e.databases[databaseID].forms[formID]
	defer fm.mu.Unlock()
	fm.mu.Lock()
	if nil != fm.indexes[indexID] {
		return
	}
	fm.indexes[indexID] = &index{path: path}
}

func (e *engine) form(databaseID, formID, path string) *form {
	if db, exist := e.databases[databaseID]; exist {
		if fm, ok := db.forms[formID]; !ok {
			e.mkFormTry(databaseID, formID, path)
		} else {
			return fm
		}
	} else {
		e.mkDatabase(databaseID)
		e.mkFormTry(databaseID, formID, path)
	}
	return e.databases[databaseID].forms[formID]
}

func (e *engine) index(databaseID, formID, indexID, formFilePath, formIndexFilePath string) *index {
	if db, exist := e.databases[databaseID]; exist {
		if fm, ok := db.forms[formID]; !ok {
			e.mkFormTry(databaseID, formID, formFilePath)
			e.mkIndexTry(databaseID, formID, indexID, formIndexFilePath)
		} else {
			if idx, yes := fm.indexes[indexID]; !yes {
				e.mkIndexTry(databaseID, formID, indexID, formIndexFilePath)
			} else {
				return idx
			}
		}
	} else {
		e.mkDatabase(databaseID)
		e.mkFormTry(databaseID, formID, formFilePath)
		e.mkIndexTry(databaseID, formID, indexID, formIndexFilePath)
	}
	return e.databases[databaseID].forms[formID].indexes[indexID]
}

// Handler 存储回调
//
// 索引最终存储在文件中的起始位置
//
// value最终存储在文件中的起始位置
//
// value最终存储在文件中的持续长度
type Handler func(SeekStartIndex int64, SeekStart int64, SeekLast int)

// Write 索引即将写入的参考坐标
type Write struct {
	IndexID           string  // 索引ID
	FormIndexFilePath string  // 索引文件所在路径
	MD516Key          string  // 索引对应字符串key
	HashKey           uint64  // put hash hashKey
	SeekStartIndex    int64   // 上一索引最终存储在文件中的起始位置
	Handler           Handler // 存储回调mu         sync.Mutex
}

// Read 数据读取结果
type Read struct {
	Key   string      // key
	Value interface{} // value
}
