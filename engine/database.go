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

package engine

import (
	"github.com/aberic/gnomon"
	"github.com/aberic/lilydb/connector"
	api "github.com/aberic/lilydb/connector/grpc"
	"github.com/aberic/lilydb/engine/comm"
	"github.com/aberic/lilydb/engine/msiam"
	"github.com/aberic/lilydb/engine/siam"
	"strings"
	"sync"
)

// database 数据库对象
//
// 存储格式 {dataDir}/database/{dataName}/{formName}/{formName}.dat/idx...
type database struct {
	id      string                    // 数据库唯一ID，不能改变
	name    string                    // 数据库名称，根据需求可以随时变化
	comment string                    // 描述
	forms   map[string]connector.Form // 表集合
	mu      sync.Mutex
}

// newForm 新建表，会创建默认自增主键
//
// formName 表名称
//
// comment 表描述
//
// formType 表类型
func (db *database) newForm(formName, comment string, formType api.FormType) error {
	defer db.mu.Unlock()
	db.mu.Lock()
	// 确定库名不重复
	for k := range db.forms {
		if k == formName {
			return comm.ErrFormExist
		}
	}
	// 确保数据库唯一ID不重复
	formID := db.name2ID(formName)
	switch formType {
	default:
		panic("form type error")
	case api.FormType_Siam:
		db.forms[formName] = siam.NewForm(db.id, formID, formName, comment)
	case api.FormType_MSiam:
		db.forms[formName] = msiam.NewForm(db.id, formID, formName, comment)
	}
	return nil
}

// Put 新增数据
//
// key 插入的key
//
// value 插入数据对象
//
// 返回 hashKey
func (db *database) put(formName, key string, value interface{}) (uint64, error) {
	if fm, exist := db.forms[formName]; exist {
		switch fm.FormType() {
		default:
			return 0, comm.ErrFormNotFoundOrSupport
		case api.FormType_MSiam:
			return fm.Put(key, value)
		}
	}
	return 0, comm.ErrFormNotFoundOrSupport
}

// Set 新增或修改数据
//
// key 插入的key
//
// value 插入数据对象
//
// 返回 hashKey
func (db *database) set(formName, key string, value interface{}) (uint64, error) {
	if fm, exist := db.forms[formName]; exist {
		switch fm.FormType() {
		default:
			return 0, comm.ErrFormNotFoundOrSupport
		case api.FormType_MSiam:
			return fm.Set(key, value)
		}
	}
	return 0, comm.ErrFormNotFoundOrSupport
}

// Get 获取数据
//
// key 指定的key
//
// 返回 获取的数据对象
func (db *database) get(formName, key string) (interface{}, error) {
	if fm, exist := db.forms[formName]; exist {
		switch fm.FormType() {
		default:
			return 0, comm.ErrFormNotFoundOrSupport
		case api.FormType_MSiam:
			return fm.Get(key)
		}
	}
	return 0, comm.ErrFormNotFoundOrSupport
}

// Del 删除数据
//
// key 指定的key
//
// 返回 删除的数据对象
func (db *database) del(formName, key string) (interface{}, error) {
	if fm, exist := db.forms[formName]; exist {
		switch fm.FormType() {
		default:
			return 0, comm.ErrFormNotFoundOrSupport
		case api.FormType_MSiam:
			return fm.Del(key)
		}
	}
	return 0, comm.ErrFormNotFoundOrSupport
}

func (db *database) insert(formName string, value interface{}) (uint64, error) {
	if fm, exist := db.forms[formName]; exist && fm.FormType() == api.FormType_Siam {
		return fm.Insert(value)
	}
	return 0, comm.ErrFormNotFoundOrSupport
}

func (db *database) update(formName string, value interface{}) (uint64, error) {
	if fm, exist := db.forms[formName]; exist && fm.FormType() == api.FormType_Siam {
		return fm.Insert(value)
	}
	return 0, comm.ErrFormNotFoundOrSupport
}

func (db *database) query(formName string, selectorBytes []byte) (int32, []interface{}, error) {
	if fm, exist := db.forms[formName]; exist {
		return fm.Select(selectorBytes)
	}
	return 0, nil, comm.ErrFormNotFoundOrSupport
}

func (db *database) delete(formName string, selectorBytes []byte) (int32, error) {
	if fm, exist := db.forms[formName]; exist {
		return fm.Delete(selectorBytes)
	}
	return 0, comm.ErrFormNotFoundOrSupport
}

// name2ID 确保表唯一ID不重复
func (db *database) name2ID(name string) string {
	id := gnomon.HashMD516(name)
	have := true
	for have {
		have = false
		for _, v := range db.forms {
			if v.ID() == id {
				have = true
				id = gnomon.HashMD516(strings.Join([]string{id, gnomon.StringRandSeq(3)}, ""))
				break
			}
		}
	}
	return id
}
