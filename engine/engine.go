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
	"github.com/aberic/lilydb/config"
	"github.com/aberic/lilydb/connector"
	"github.com/aberic/lilydb/engine/comm"
	"os"
	"path/filepath"
	"sync"
)

var (
	engine *Engine
	once   sync.Once
)

// Obtain 获取表引擎
func Obtain() *Engine {
	once.Do(func() {
		engine = &Engine{
			databases: map[string]*database{},
		}
	})
	return engine
}

// Engine 存储引擎管理器
//
// 全库唯一常住内存对象，并持有所有库的句柄
//
// API 入口
//
// 存储格式 {dataDir}/Data/{dataName}/{formName}/{formName}.dat/idx...
type Engine struct {
	databases map[string]*database
	mu        sync.Mutex
}

// NewDatabase 新建数据库
//
// 新建数据库会同时创建一个名为_default的表，未指定表明的情况下使用put/get等方法会操作该表
//
// databaseID 数据库唯一ID，不能改变
//
// databaseName 数据库名称
//
// comment 数据库描述
func (e *Engine) NewDatabase(databaseID, databaseName, comment string) error {
	if err := mkDataDir(databaseID); nil != err {
		return err
	}
	defer e.mu.Unlock()
	e.mu.Lock()
	e.databases[databaseID] = &database{id: databaseID, name: databaseName, comment: comment, forms: map[string]connector.Form{}}
	return nil
}

// NewForm 新建表，会创建默认自增主键
//
// databaseID 数据库唯一ID，不能改变
//
// formID 表唯一ID，不能改变
//
// formName 表名称
//
// comment 表描述
//
// formType 表类型
func (e *Engine) NewForm(databaseID, formID, formName, comment string, formType connector.FormType) error {
	if db, exist := e.databases[databaseID]; exist {
		db.newForm(formID, formName, comment, formType)
		return nil
	}
	return comm.ErrDataNotFound
}

// Put 新增数据
//
// databaseID 数据库唯一ID
//
// formID 表唯一ID
//
// key 插入的key
//
// value 插入数据对象
//
// 返回 hashKey
func (e *Engine) Put(databaseID, formID, key string, value interface{}) (uint64, error) {
	if db, exist := e.databases[databaseID]; exist {
		return db.put(formID, key, value)
	}
	return 0, comm.ErrDataNotFound
}

// Set 新增或修改数据
//
// databaseID 数据库唯一ID
//
// formID 表唯一ID
//
// key 插入的key
//
// value 插入数据对象
//
// 返回 hashKey
func (e *Engine) Set(databaseID, formID, key string, value interface{}) (uint64, error) {
	if db, exist := e.databases[databaseID]; exist {
		return db.set(formID, key, value)
	}
	return 0, comm.ErrDataNotFound
}

// Get 获取数据
//
// databaseID 数据库唯一ID
//
// formID 表唯一ID
//
// key 指定的key
//
// 返回 获取的数据对象
func (e *Engine) Get(databaseID, formID, key string) (interface{}, error) {
	if db, exist := e.databases[databaseID]; exist {
		return db.get(formID, key)
	}
	return 0, comm.ErrDataNotFound
}

// Del 删除数据
//
// databaseID 数据库唯一ID
//
// formID 表唯一ID
//
// key 指定的key
//
// 返回 删除的数据对象
func (e *Engine) Del(databaseID, formID, key string) (interface{}, error) {
	if db, exist := e.databases[databaseID]; exist {
		return db.del(formID, key)
	}
	return 0, comm.ErrDataNotFound
}

// Insert 新增数据
//
// databaseID 数据库唯一ID
//
// formID 表唯一ID
//
// value 插入数据对象
//
// 返回 hashKey
func (e *Engine) Insert(databaseID, formID string, value interface{}) (uint64, error) {
	if db, exist := e.databases[databaseID]; exist {
		return db.insert(formID, value)
	}
	return 0, comm.ErrDataNotFound
}

// Update 更新数据，如果存在数据，则更新，如不存在，则插入
//
// databaseID 数据库唯一ID
//
// formID 表唯一ID
//
// value 插入数据对象
//
// 返回 hashKey
func (e *Engine) Update(databaseID, formID string, value interface{}) (uint64, error) {
	if db, exist := e.databases[databaseID]; exist {
		return db.update(formID, value)
	}
	return 0, comm.ErrDataNotFound
}

// Select 根据条件检索
//
// databaseID 数据库唯一ID
//
// formID 表唯一ID
//
// selectorBytes 选择器字节数组，自定义转换策略
//
// return count 检索结果总条数
//
// return values 检索结果集合
//
// return err 检索错误信息，如果有
func (e *Engine) Select(databaseID, formID string, selectorBytes []byte) (count int32, values []interface{}, err error) {
	if db, exist := e.databases[databaseID]; exist {
		return db.query(formID, selectorBytes)
	}
	return 0, nil, comm.ErrDataNotFound
}

// Delete 根据条件删除
//
// databaseID 数据库唯一ID
//
// formID 表唯一ID
//
// selectorBytes 选择器字节数组，自定义转换策略
//
// return count 删除结果总条数
//
// return err 删除错误信息，如果有
func (e *Engine) Delete(databaseID, formID string, selectorBytes []byte) (count int32, err error) {
	if db, exist := e.databases[databaseID]; exist {
		return db.delete(formID, selectorBytes)
	}
	return 0, comm.ErrDataNotFound
}

// mkDataDir 创建库存储目录
func mkDataDir(dataName string) (err error) {
	dataPath := filepath.Join(config.Obtain().DataDir, dataName)
	if gnomon.FilePathExists(dataPath) {
		return comm.ErrDatabaseExist
	}
	return os.MkdirAll(dataPath, os.ModePerm)
}
