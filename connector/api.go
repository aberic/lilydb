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

package connector

import api "github.com/aberic/lilydb/connector/grpc"

// Code 返回码
type Code int

const (
	// Success 返回成功
	Success Code = iota
	// Fail 返回失败
	Fail
)

// Request 请求对象
type Request interface {
}

// Response 返回对象
type Response interface {
	Code() (code Code)
	Error() (error error)
	Data() (value interface{})
}

// Form 表接口
//
// 提供表基本操作方法
type Form interface {
	AutoID() *uint64        // AutoID 返回表当前自增ID值
	ID() string             // ID 返回表唯一ID
	Name() string           // Name 返回表名称
	Comment() string        // Comment 获取表描述
	FormType() api.FormType // FormType 获取表类型
	Indexes() map[string]*api.Index
	// Insert 新增数据
	//
	// value 插入数据对象
	//
	// 返回 hashKey
	Insert(value interface{}) (uint64, error)
	// Update 更新数据，如果存在数据，则更新，如不存在，则插入
	//
	// value 插入数据对象
	//
	// 返回 hashKey
	Update(value interface{}) (uint64, error)
	// Put 新增数据
	//
	// key 插入的key
	//
	// value 插入数据对象
	//
	// 返回 hashKey
	Put(ket string, value interface{}) (uint64, error)
	// Set 新增或修改数据
	//
	// key 插入的key
	//
	// value 插入数据对象
	//
	// 返回 hashKey
	Set(ket string, value interface{}) (uint64, error)
	// Get 获取数据
	//
	// key 指定的key
	//
	// 返回 获取的数据对象
	Get(ket string) (interface{}, error)
	// Del 删除数据
	//
	// key 指定的key
	//
	// 返回 删除的数据对象
	Del(ket string) (interface{}, error)
	// Select 根据条件检索
	//
	// selectorBytes 选择器字节数组，自定义转换策略
	//
	// return count 检索结果总条数
	//
	// return values 检索结果集合
	//
	// return err 检索错误信息，如果有
	Select(selectorBytes []byte) (count int32, values []interface{}, err error)
	// Delete 根据条件删除
	//
	// selectorBytes 选择器字节数组，自定义转换策略
	//
	// return count 删除结果总条数
	//
	// return err 删除错误信息，如果有
	Delete(selectorBytes []byte) (count int32, err error)
}
