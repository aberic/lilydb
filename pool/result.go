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
)

func resultSuccess(value interface{}) *Result {
	return &Result{code: connector.Success, value: value}
}

func resultFail(error error) *Result {
	return &Result{code: connector.Fail, error: error}
}

// Result 返回对象
type Result struct {
	code  connector.Code
	error error
	value interface{}
}

// Code 返回码
func (r *Result) Code() (code connector.Code) {
	return r.code
}

// Error 错误信息
func (r *Result) Error() (error error) {
	return r.error
}

// Data 数据
func (r *Result) Data() (value interface{}) {
	return r.value
}
