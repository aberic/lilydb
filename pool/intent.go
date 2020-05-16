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

package pool

import (
	"github.com/aberic/lilydb/connector"
	"github.com/aberic/lilydb/engine"
)

// task 任务对象
type task struct {
	intent  Intent  // 意图
	handler Handler // 执行后的回调
}

// Handler 执行后的回调
type Handler func(connector.Response)

// Intent 意图接口
type Intent interface {
	run(engine *engine.Engine, handler Handler) // 执行
}

// IntentNewDatabase 新建库意图
type IntentNewDatabase struct {
	databaseID   string
	databaseName string
	comment      string
}

func (i *IntentNewDatabase) run(engine *engine.Engine, handler Handler) {
	err := engine.NewDatabase(i.databaseID, i.databaseName, i.comment)
	if nil != err {
		handler(connector.ResultFail(err))
	}
	handler(connector.ResultSuccess(nil))
}

// IntentNewForm 新建表意图
type IntentNewForm struct {
	databaseID string
	formID     string
	formName   string
	comment    string
	formType   connector.FormType
}

func (i *IntentNewForm) run(engine *engine.Engine, handler Handler) {
	err := engine.NewForm(i.databaseID, i.formID, i.formName, i.comment, i.formType)
	if nil != err {
		handler(connector.ResultFail(err))
	}
	handler(connector.ResultSuccess(nil))
}

// IntentPut 新建新增数据意图
type IntentPut struct {
	databaseID string
	formID     string
	key        string
	value      interface{}
}

func (i *IntentPut) run(engine *engine.Engine, handler Handler) {
	hashKey, err := engine.Put(i.databaseID, i.formID, i.key, i.value)
	if nil != err {
		handler(connector.ResultFail(err))
	}
	handler(connector.ResultSuccess(hashKey))
}

// IntentSet 新建新增或修改数据意图
type IntentSet struct {
	databaseID string
	formID     string
	key        string
	value      interface{}
}

func (i *IntentSet) run(engine *engine.Engine, handler Handler) {
	hashKey, err := engine.Set(i.databaseID, i.formID, i.key, i.value)
	if nil != err {
		handler(connector.ResultFail(err))
	}
	handler(connector.ResultSuccess(hashKey))
}

// IntentGet 新建获取数据意图
type IntentGet struct {
	databaseID string
	formID     string
	key        string
}

func (i *IntentGet) run(engine *engine.Engine, handler Handler) {
	value, err := engine.Get(i.databaseID, i.formID, i.key)
	if nil != err {
		handler(connector.ResultFail(err))
	}
	handler(connector.ResultSuccess(value))
}

// IntentDel 新建删除数据意图
type IntentDel struct {
	databaseID string
	formID     string
	key        string
}

func (i *IntentDel) run(engine *engine.Engine, handler Handler) {
	value, err := engine.Del(i.databaseID, i.formID, i.key)
	if nil != err {
		handler(connector.ResultFail(err))
	}
	handler(connector.ResultSuccess(value))
}

// IntentInsert 新建新增数据意图
type IntentInsert struct {
	databaseID string
	formID     string
	value      interface{}
}

func (i *IntentInsert) run(engine *engine.Engine, handler Handler) {
	hashKey, err := engine.Insert(i.databaseID, i.formID, i.value)
	if nil != err {
		handler(connector.ResultFail(err))
	}
	handler(connector.ResultSuccess(hashKey))
}

// IntentUpdate 新建更新数据意图
type IntentUpdate struct {
	databaseID string
	formID     string
	value      interface{}
}

func (i *IntentUpdate) run(engine *engine.Engine, handler Handler) {
	hashKey, err := engine.Update(i.databaseID, i.formID, i.value)
	if nil != err {
		handler(connector.ResultFail(err))
	}
	handler(connector.ResultSuccess(hashKey))
}

// IntentSelect 新建根据条件检索意图
type IntentSelect struct {
	databaseID    string
	formID        string
	selectorBytes []byte
}

type data struct {
	Count  int32
	Values []interface{}
}

func (i *IntentSelect) run(engine *engine.Engine, handler Handler) {
	count, values, err := engine.Select(i.databaseID, i.formID, i.selectorBytes)
	if nil != err {
		handler(connector.ResultFail(err))
	}
	handler(connector.ResultSuccess(&data{Count: count, Values: values}))
}

// IntentDelete 新建根据条件删除意图
type IntentDelete struct {
	databaseID    string
	formID        string
	selectorBytes []byte
}

func (i *IntentDelete) run(engine *engine.Engine, handler Handler) {
	count, err := engine.Delete(i.databaseID, i.formID, i.selectorBytes)
	if nil != err {
		handler(connector.ResultFail(err))
	}
	handler(connector.ResultSuccess(count))
}
