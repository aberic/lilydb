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

package index

// NewIndex 新建索引
//
// databaseID 数据库唯一ID
//
// formID 表唯一ID
//
// id 索引唯一ID
//
// keyStructure 按照规范结构组成的索引字段名称，由对象结构层级字段通过'.'组成，如'i','in.s'
//
// primary 是否主键
func NewIndex(databaseID, formID, id, keyStructure string, primary bool) *Index {
	return &Index{
		id:           id,
		primary:      primary,
		keyStructure: keyStructure,
		node:         &node{level: 1, degreeIndex: 0, preNode: nil, nodes: []*node{}},
		databaseID:   databaseID,
		formID:       formID,
	}
}

// Index Siam索引
//
// 5位key及16位md5后key及5位起始seek和4位持续seek
type Index struct {
	id           string // id 索引唯一ID
	primary      bool   // 是否主键
	keyStructure string // keyStructure 按照规范结构组成的索引字段名称，由对象结构层级字段通过'.'组成，如'i','in.s'
	node         *node  // 节点
	databaseID   string // 所属数据库ID
	formID       string // 所属表ID
}

// ID 索引唯一ID
func (i *Index) ID() string {
	return i.id
}

// Primary 是否主键
func (i *Index) Primary() bool {
	return i.primary
}

// KeyStructure 索引字段名称，由对象结构层级字段通过'.'组成，如
func (i *Index) KeyStructure() string {
	return i.keyStructure
}

// Put 插入数据
//
// key 必须string类型
//
// md516Key md516Key，必须string类型
//
// hashKey 索引key，可通过hash转换string生成，下一级最左最小树所对应真实key
//
// value 存储对象
//
// version 当前索引数据版本号
func (i *Index) Put(key, md516Key string, hashKey uint64, value interface{}, version int) (link *Link, exist, versionGT bool) {
	return i.node.put(key, md516Key, hashKey, value, version)
}

// Get 获取数据，返回存储对象
//
// key 真实key，必须string类型
//
// hashKey 索引key，可通过hash转换string生成
func (i *Index) Get(md516Key string, hashKey uint64) *Link {
	return i.node.get(md516Key, hashKey, hashKey)
}

// Del 删除数据
//
// key 真实key，必须string类型
//
// hashKey 索引key，可通过hash转换string生成
func (i *Index) Del(md516Key string, hashKey uint64) (interface{}, error) {
	return i.node.del(md516Key, hashKey, hashKey)
}
