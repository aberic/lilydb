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
	api "github.com/aberic/lilydb/connector/grpc"
	"github.com/aberic/lilydb/engine/comm"
	"os"
	"path/filepath"
	"strings"
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

// Databases 获取数据库集合
func (e *Engine) Databases() []*api.Database {
	var respDBs []*api.Database
	for _, db := range e.databaseArr() {
		respDBs = append(respDBs, &api.Database{ID: db.id, Name: db.name, Comment: db.comment, Forms: e.formatForms(db)})
	}
	return respDBs
}

// databaseArr 获取数据库集合
func (e *Engine) databaseArr() []*database {
	var dbs []*database
	for _, db := range e.databases {
		dbs = append(dbs, db)
	}
	return dbs
}

func (e *Engine) formatForms(db *database) map[string]*api.Form {
	var fms = make(map[string]*api.Form)
	for _, form := range db.forms {
		fms[form.Name()] = &api.Form{
			ID:       form.ID(),
			Name:     form.Name(),
			Comment:  form.Comment(),
			FormType: form.FormType(),
			Indexes:  form.Indexes(),
		}
	}
	return fms
}

// Forms 根据数据库名获取表集合
func (e *Engine) Forms(databaseName string) []*api.Form {
	return e.databases[databaseName].formArr()
}

// NewDatabase 新建数据库
//
// 新建数据库会同时创建一个名为_default的表，未指定表明的情况下使用put/get等方法会操作该表
//
// databaseName 数据库名称
//
// comment 数据库描述
func (e *Engine) NewDatabase(databaseName, comment string) error {
	defer e.mu.Unlock()
	e.mu.Lock()
	// 确定库名不重复
	for k := range e.databases {
		if k == databaseName {
			return comm.ErrDatabaseExist
		}
	}
	// 确保数据库唯一ID不重复
	databaseID := e.name2ID(databaseName)
	if err := e.mkDataDir(databaseID); nil != err {
		return err
	}
	e.databases[databaseName] = &database{id: databaseID, name: databaseName, comment: comment, forms: map[string]connector.Form{}}
	return nil
}

// NewForm 新建表，会创建默认自增主键
//
// databaseName 数据库名
//
// formName 表名称
//
// comment 表描述
//
// formType 表类型
func (e *Engine) NewForm(databaseName, formName, comment string, formType api.FormType) error {
	if db, exist := e.databases[databaseName]; exist {
		return db.newForm(formName, comment, formType)
	}
	return comm.ErrDataNotFound
}

// Put 新增数据
//
// databaseID 数据库名
//
// formName 表名
//
// key 插入的key
//
// value 插入数据对象
//
// 返回 hashKey
func (e *Engine) Put(databaseName, formName, key string, value interface{}) (uint64, error) {
	if db, exist := e.databases[databaseName]; exist {
		return db.put(formName, key, value)
	}
	return 0, comm.ErrDataNotFound
}

// Set 新增或修改数据
//
// databaseID 数据库名
//
// formName 表名
//
// key 插入的key
//
// value 插入数据对象
//
// 返回 hashKey
func (e *Engine) Set(databaseName, formName, key string, value interface{}) (uint64, error) {
	if db, exist := e.databases[databaseName]; exist {
		return db.set(formName, key, value)
	}
	return 0, comm.ErrDataNotFound
}

// Get 获取数据
//
// databaseID 数据库名
//
// formName 表名
//
// key 指定的key
//
// 返回 获取的数据对象
func (e *Engine) Get(databaseName, formName, key string) (interface{}, error) {
	if db, exist := e.databases[databaseName]; exist {
		return db.get(formName, key)
	}
	return 0, comm.ErrDataNotFound
}

// Del 删除数据
//
// databaseID 数据库名
//
// formName 表名
//
// key 指定的key
//
// 返回 删除的数据对象
func (e *Engine) Del(databaseName, formName, key string) (interface{}, error) {
	if db, exist := e.databases[databaseName]; exist {
		return db.del(formName, key)
	}
	return 0, comm.ErrDataNotFound
}

// Insert 新增数据
//
// databaseID 数据库名
//
// formName 表名
//
// value 插入数据对象
//
// 返回 hashKey
func (e *Engine) Insert(databaseName, formName string, value interface{}) (uint64, error) {
	if db, exist := e.databases[databaseName]; exist {
		return db.insert(formName, value)
	}
	return 0, comm.ErrDataNotFound
}

// Update 更新数据，如果存在数据，则更新，如不存在，则插入
//
// databaseID 数据库名
//
// formName 表名
//
// value 插入数据对象
//
// 返回 hashKey
func (e *Engine) Update(databaseName, formName string, value interface{}) (uint64, error) {
	if db, exist := e.databases[databaseName]; exist {
		return db.update(formName, value)
	}
	return 0, comm.ErrDataNotFound
}

// Select 根据条件检索
//
// databaseID 数据库名
//
// formName 表名
//
// selectorBytes 选择器字节数组，自定义转换策略
//
// return count 检索结果总条数
//
// return values 检索结果集合
//
// return err 检索错误信息，如果有
func (e *Engine) Select(databaseName, formName string, selectorBytes []byte) (count int32, values []interface{}, err error) {
	if db, exist := e.databases[databaseName]; exist {
		return db.query(formName, selectorBytes)
	}
	return 0, nil, comm.ErrDataNotFound
}

// Delete 根据条件删除
//
// databaseID 数据库名
//
// formName 表名
//
// selectorBytes 选择器字节数组，自定义转换策略
//
// return count 删除结果总条数
//
// return err 删除错误信息，如果有
func (e *Engine) Delete(databaseName, formName string, selectorBytes []byte) (count int32, err error) {
	if db, exist := e.databases[databaseName]; exist {
		return db.delete(formName, selectorBytes)
	}
	return 0, comm.ErrDataNotFound
}

// name2ID 确保数据库唯一ID不重复
func (e *Engine) name2ID(name string) string {
	id := gnomon.HashMD516(name)
	have := true
	for have {
		have = false
		for _, v := range e.databases {
			if v.id == id {
				have = true
				id = gnomon.HashMD516(strings.Join([]string{id, gnomon.StringRandSeq(3)}, ""))
				break
			}
		}
	}
	return id
}

// mkDataDir 创建库存储目录
func (e *Engine) mkDataDir(dataName string) (err error) {
	dataPath := filepath.Join(config.Obtain().DataDir, dataName)
	if gnomon.FilePathExists(dataPath) {
		return comm.ErrDatabaseExist
	}
	return os.MkdirAll(dataPath, os.ModePerm)
}
