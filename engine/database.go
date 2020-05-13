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
	"github.com/aberic/lilydb/connector"
	"github.com/aberic/lilydb/engine/comm"
	"github.com/aberic/lilydb/engine/siam"
)

// database 数据库对象
//
// 存储格式 {dataDir}/database/{dataName}/{formName}/{formName}.dat/idx...
type database struct {
	id      string                    // 数据库唯一ID，不能改变
	name    string                    // 数据库名称，根据需求可以随时变化
	comment string                    // 描述
	forms   map[string]connector.Form // 表集合
}

// newForm 新建表，会创建默认自增主键
//
// formID 表唯一ID，不能改变
//
// formName 表名称
//
// comment 表描述
//
// formType 表类型
func (db *database) newForm(formID, formName, comment string, formType connector.FormType) {
	switch formType {
	default:
		panic("form type error")
	case connector.FormTypeSiam:
		db.forms[formID] = siam.NewForm(db.id, formID, formName, comment)
	}
}

func (db *database) insert(formID string, value interface{}) (uint64, error) {
	if fm, exist := db.forms[formID]; exist && fm.FormType() == connector.FormTypeSiam {
		return fm.Insert(value)
	}
	return 0, comm.ErrFormNotFoundOrSupport
}

func (db *database) update(formID string, value interface{}) (uint64, error) {
	if fm, exist := db.forms[formID]; exist && fm.FormType() == connector.FormTypeSiam {
		return fm.Insert(value)
	}
	return 0, comm.ErrFormNotFoundOrSupport
}

func (db *database) query(formID string, selectorBytes []byte) (int32, []interface{}, error) {
	if fm, exist := db.forms[formID]; exist {
		return fm.Select(selectorBytes)
	}
	return 0, nil, comm.ErrFormNotFoundOrSupport
}

func (db *database) delete(formID string, selectorBytes []byte) (int32, error) {
	if fm, exist := db.forms[formID]; exist {
		return fm.Delete(selectorBytes)
	}
	return 0, comm.ErrFormNotFoundOrSupport
}
