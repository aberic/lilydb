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
	Run(engine *engine.Engine, handler Handler)
}

// IntentNewForm 新建表意图
type IntentNewForm struct {
	databaseID string
	formID     string
	formName   string
	comment    string
	formType   connector.FormType
}

// Run 执行
func (i *IntentNewForm) Run(engine *engine.Engine, handler Handler) {
	err := engine.NewForm(i.databaseID, i.formID, i.formName, i.comment, i.formType)
	if nil != err {
		handler(resultFail(err))
	}
	handler(resultSuccess(nil))
}
