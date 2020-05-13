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

import (
	"fmt"
	"github.com/aberic/gnomon"
	"github.com/aberic/lilydb/connector"
	"github.com/aberic/lilydb/engine/siam/index"
	"github.com/aberic/lilydb/engine/siam/storage"
	"github.com/aberic/lilydb/engine/siam/utils"
	"reflect"
	"strconv"
	"strings"
	"sync"
	"sync/atomic"
)

const indexAutoID = "lily_do_not_repeat_auto_id"

// NewForm 新建表，会创建默认自增主键
//
// 所属数据库ID
//
// formID 表唯一ID
//
// formName 表名，根据需求可以随时变化
//
// comment 描述
func NewForm(databaseID, formID, formName, comment string) *Form {
	var autoID uint64 = 0
	fm := &Form{
		autoID:     &autoID,
		name:       formName,
		id:         formID,
		comment:    comment,
		indexes:    map[string]*index.Index{},
		formType:   connector.FormTypeSiam,
		databaseID: databaseID,
	}
	fm.NewIndex(indexAutoID, true) // 创建默认自增主键
	return fm
}

// Form 表结构
type Form struct {
	id         string                  // 表唯一ID，不能改变
	name       string                  // 表名，根据需求可以随时变化
	autoID     *uint64                 // 自增id
	comment    string                  // 描述
	formType   connector.FormType      // 表类型 siam
	indexes    map[string]*index.Index // 索引ID集合
	databaseID string                  // 所属数据库ID

	mu sync.RWMutex
}

// AtomicAddAutoID 自增ID
func (f *Form) atomicAddAutoID() uint64 {
	return atomic.AddUint64(f.autoID, 1)
}

// AutoID 返回表当前自增ID值
func (f *Form) AutoID() *uint64 {
	return f.autoID
}

// ID 返回表唯一ID
func (f *Form) ID() string {
	return f.id
}

// Name 返回表名称
func (f *Form) Name() string {
	return f.name
}

// Comment 获取表描述
func (f *Form) Comment() string {
	return f.comment
}

// FormType 获取表类型
func (f *Form) FormType() connector.FormType {
	return f.formType
}

// NewIndex 新建索引
//
// keyStructure 按照规范结构组成的索引字段名称，由对象结构层级字段通过'.'组成，如'i','in.s'
//
// primary 是否主键
func (f *Form) NewIndex(keyStructure string, primary bool) {
	defer f.mu.Unlock()
	f.mu.Lock()
	indexID := f.name2ID4Index(strings.Join([]string{f.name, keyStructure}, "_"))
	f.indexes[indexID] = index.NewIndex(f.databaseID, f.id, indexID, keyStructure, primary)
}

// name2ID4Index 确保表下索引唯一ID不重复
func (f *Form) name2ID4Index(name string) string {
	id := gnomon.HashMD516(name)
	have := true
	for have {
		have = false
		for _, idx := range f.indexes {
			if idx.ID() == id {
				have = true
				id = gnomon.HashMD516(strings.Join([]string{id, gnomon.StringRandSeq(3)}, ""))
				break
			}
		}
	}
	return id
}

// Insert 新增数据
//
// databaseID 数据库唯一ID
//
// value 插入数据对象
//
// 返回 hashKey
func (f *Form) Insert(value interface{}) (uint64, error) {
	//gnomon.Log().Debug("insertDataWithIndexInfo", gnomon.Log().Field("ibs", ibs))
	defer f.mu.Unlock()
	f.mu.Lock()
	// 遍历表索引ID集合，检索并计算当前索引所在文件位置，存储结果
	if err := f.store(value, false); nil != err {
		return 0, err
	}
	return *f.autoID, nil
}

// Update 更新数据，如果存在数据，则更新，如不存在，则插入
//
// databaseID 数据库唯一ID
//
// value 插入数据对象
//
// 返回 hashKey
func (f *Form) Update(value interface{}) (uint64, error) {
	//gnomon.Log().Debug("insertDataWithIndexInfo", gnomon.Log().Field("ibs", ibs))
	defer f.mu.Unlock()
	f.mu.Lock()
	// 遍历表索引ID集合，检索并计算当前索引所在文件位置，存储结果
	if err := f.store(value, true); nil != err {
		return 0, err
	}
	return *f.autoID, nil
}

// Select 根据条件检索
//
// databaseID 数据库唯一ID
//
// selectorBytes 选择器字节数组，自定义转换策略
//
// return count 检索结果总条数
//
// return values 检索结果集合
//
// return err 检索错误信息，如果有
func (f *Form) Select(selectorBytes []byte) (int32, []interface{}, error) {
	var indexes []*index.Index
	for _, idx := range f.indexes {
		indexes = append(indexes, idx)
	}
	selector, err := index.NewSelector(selectorBytes, indexes, f.databaseID, f.id, false)
	if nil != err {
		return 0, nil, err
	}
	count, values := selector.Run()
	return count, values, nil
}

// Delete 根据条件删除
//
// databaseID 数据库唯一ID
//
// selectorBytes 选择器字节数组，自定义转换策略
//
// return count 删除结果总条数
//
// return err 删除错误信息，如果有
func (f *Form) Delete(selectorBytes []byte) (int32, error) {
	var indexes []*index.Index
	for _, idx := range f.indexes {
		indexes = append(indexes, idx)
	}
	selector, err := index.NewSelector(selectorBytes, indexes, f.databaseID, f.id, true)
	if nil != err {
		return 0, err
	}
	count, _ := selector.Run()
	return count, nil
}

// rangeIndexes 遍历表索引ID集合，检索并计算所有索引返回对象集合
func (f *Form) store(value interface{}, update bool) error {
	var (
		wg     sync.WaitGroup
		writes []*storage.Write
		wMu    sync.Mutex
		err    error
	)
	// 遍历表索引ID集合，检索并计算当前索引所在文件位置
	for _, idx := range f.indexes {
		wg.Add(1)
		go func(index *index.Index) {
			defer wg.Done()
			var (
				key     string
				hashKey uint64
			)
			//gnomon.Log().Debug("rangeIndexes", gnomon.Log().Field("index.id", index.getID()), gnomon.Log().Field("index.keyStructure", index.getKeyStructure()))
			if index.KeyStructure() == indexAutoID {
				hashKey = f.atomicAddAutoID() // ID自增
				key = strconv.FormatUint(hashKey, 10)
			} else {
				key, hashKey, err = f.getCustomIndex(index, value)
				if nil != err {
					return
				}
			}
			md516Key := gnomon.HashMD516(key)
			link, exist, _ := f.indexes[index.ID()].Put(md516Key, hashKey, 0)
			if !update && exist { // 如果当前是插入操作，且已存在对应key的值
				err = fmt.Errorf("the same key %s already exist", index.KeyStructure())
				return
			}
			defer wMu.Unlock()
			wMu.Lock()
			writes = append(writes, &storage.Write{
				IndexID:           index.ID(),
				FormIndexFilePath: utils.PathFormIndexFile(f.databaseID, f.id, index.ID()),
				MD516Key:          md516Key,
				HashKey:           hashKey,
				SeekStartIndex:    link.SeekStartIndex(),
				Handler: func(SeekStartIndex int64, SeekStart int64, SeekLast int) {
					link.Fit(SeekStartIndex, SeekStart, SeekLast, 0)
				},
			})

		}(idx)
	}
	wg.Wait()
	if nil != err {
		return err
	}
	return storage.Obtain().Store(f.databaseID, f.id, value, writes)
}

// getCustomIndex 获取自定义索引预插入返回对象
func (f *Form) getCustomIndex(idx *index.Index, value interface{}) (key string, hashKey uint64, err error) {
	reflectValue := reflect.ValueOf(value) // 反射对象，通过reflectObj获取存储在里面的值，还可以去改变值
	params := strings.Split(idx.KeyStructure(), ".")
	switch kind := reflectValue.Kind(); kind {
	default:
		err = fmt.Errorf("index %s with type is invalid", idx.KeyStructure())
		return
	case reflect.Map:
		var (
			item      interface{}
			paramsLen = len(params)
			position  int
			itemMap   = value.(map[string]interface{})
		)
		for _, param := range params {
			position++
			item = itemMap[param]
			if position == paramsLen { // 表示没有后续参数
				break
			}
			switch item := item.(type) {
			default:
				err = fmt.Errorf("index %s with map is invalid", idx.KeyStructure())
				return
			case map[string]interface{}:
				itemMap = item
				continue
			}
		}
		if keyNew, hashKeyNew, valid := utils.Type2index(item); valid {
			return keyNew, hashKeyNew, nil
		}
		err = fmt.Errorf("index %s with map is invalid", idx.KeyStructure())
		return
	case reflect.Ptr:
		checkValue := reflectValue
		for _, param := range params {
			checkNewValue := checkValue.Elem().FieldByName(param)
			if checkNewValue.IsValid() { // 子字段有效
				checkValue = checkNewValue
				continue
			}
			err = fmt.Errorf("index %s with ptr is invalid", idx.KeyStructure())
			return
		}
		if keyNew, hashKeyNew, valid := utils.ValueType2index(&checkValue); valid {
			return keyNew, hashKeyNew, nil
		}
		err = fmt.Errorf("index %s with ptr is invalid", idx.KeyStructure())
		return
	}
}
