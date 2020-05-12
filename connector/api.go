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

// FormType 表引擎类型
type FormType int

const (
	// FormTypeSiam 静态索引方法(static index access method)
	FormTypeSiam FormType = iota
)

// Form 表接口
//
// 提供表基本操作方法
type Form interface {
	AutoID() *uint64    // AutoID 返回表当前自增ID值
	ID() string         // ID 返回表唯一ID
	Name() string       // Name 返回表名称
	Comment() string    // Comment 获取表描述
	FormType() FormType // FormType 获取表类型
	// Insert 新增数据
	//
	// databaseID 数据库唯一ID
	//
	// value 插入数据对象
	//
	// 返回 hashKey
	Insert(value interface{}) (uint64, error)
	// Update 更新数据，如果存在数据，则更新，如不存在，则插入
	//
	// databaseID 数据库唯一ID
	//
	// value 插入数据对象
	//
	// 返回 hashKey
	Update(value interface{}) (uint64, error)
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
	Select(selectorBytes []byte) (count int32, values []interface{}, err error)
	// Delete 根据条件删除
	//
	// databaseID 数据库唯一ID
	//
	// selectorBytes 选择器字节数组，自定义转换策略
	//
	// return count 删除结果总条数
	//
	// return err 删除错误信息，如果有
	Delete(selectorBytes []byte) (count int32, err error)
}

// Index 索引接口
type Index interface {
	// ID 索引唯一ID
	ID() string
	// Primary 是否主键
	Primary() bool
	// KeyStructure 索引字段名称，由对象结构层级字段通过'.'组成，如
	//
	// ref := &ref{
	//		i: 1,
	//		s: "2",
	//		in: refIn{
	//			i: 3,
	//			s: "4",
	//		},
	//	}
	//
	// key可取'i','in.s'
	KeyStructure() string
	// Put 插入数据
	//
	// key 必须string类型
	//
	// hashKey 索引key，可通过hash转换string生成
	//
	// value 存储对象
	//
	// update 本次是否执行更新操作
	Put(key string, hashKey uint64) (link Link, exist bool)
	// Get 获取数据，返回存储对象
	//
	// key 真实key，必须string类型
	//
	// hashKey 索引key，可通过hash转换string生成
	Get(key string, hashKey uint64) Link
	// Recover 重置索引数据
	Recover() error
}

// Link 叶子节点下的链表对象接口
type Link interface {
	// Fit 填充数据
	//
	// 索引最终存储在文件中的起始位置
	//
	// value最终存储在文件中的起始位置
	//
	// value最终存储在文件中的持续长度
	Fit(seekStartIndex int64, seekStart int64, seekLast int)
	MD516Key() string      // 获取md516Key
	SeekStartIndex() int64 // 索引最终存储在文件中的起始位置
	SeekStart() int64      // value最终存储在文件中的起始位置
	SeekLast() int         // value最终存储在文件中的持续长度
}

// Selector 检索选择器
type Selector interface {
	// Run 执行富查询
	//
	// return count 检索结果总条数
	//
	// return values 检索结果集合
	Run() (count int32, values []interface{})
}
