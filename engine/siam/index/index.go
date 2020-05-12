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

import (
	"bufio"
	"errors"
	"github.com/aberic/gnomon"
	"github.com/aberic/gnomon/log"
	"github.com/aberic/lilydb/connector"
	"github.com/aberic/lilydb/engine/siam/utils"
	"io"
	"os"
)

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
// hashKey 索引key，可通过hash转换string生成
//
// value 存储对象
//
// update 本次是否执行更新操作
func (i *Index) Put(md516Key string, hashKey uint64) (link connector.Link, exist bool) {
	return i.node.put(md516Key, hashKey, hashKey)
}

// Get 获取数据，返回存储对象
//
// key 真实key，必须string类型
//
// hashKey 索引key，可通过hash转换string生成
func (i *Index) Get(md516Key string, hashKey uint64) connector.Link {
	return i.node.get(md516Key, hashKey, hashKey)
}

// Recover 重置索引数据
func (i *Index) Recover() error {
	indexFilePath := utils.PathFormIndexFile(i.databaseID, i.formID, i.id)
	if gnomon.FilePathExists(indexFilePath) { // 索引文件存在才继续恢复
		var (
			file *os.File
			err  error
		)
		defer func() { _ = file.Close() }()
		if file, err = os.OpenFile(indexFilePath, os.O_RDONLY, 0644); nil != err {
			log.Error("index recover multi read failed", log.Err(err))
			return err
		}
		_, err = file.Seek(0, io.SeekStart) // 文件下标置为文件的起始位置
		if err != nil {
			log.Error("index recover multi read failed", log.Err(err))
			return err
		}
		if err = i.read(file, 0); nil != err && io.EOF != err {
			log.Error("index recover multi read failed", log.Err(err))
			return err
		}
	}
	return nil
}

// todo 回传错误
func (i *Index) read(file *os.File, offset int64) (err error) {
	var (
		inputReader *bufio.Reader
		data        []byte
		peekOnce    = 42000
		haveNext    = true
		position    int64
	)
	if _, err = file.Seek(offset, io.SeekStart); nil != err { //表示文件的起始位置，从第二个字符往后写入。
		return
	}
	inputReader = bufio.NewReaderSize(file, peekOnce)
	data, err = inputReader.Peek(peekOnce)
	if nil != err && io.EOF != err {
		return
	} else if nil != err && io.EOF == err {
		if len(data) == 0 {
			return
		}
		if len(data)%42 != 0 { // 单条索引默认占用长度42个字符
			return errors.New("index lens does't match")
		}
	}
	indexStr := string(data)
	indexStrLen := int64(len(indexStr))
	for haveNext {
		go func(i *Index, position int64) {
			var p0, p1, p2, p3, p4 int64
			// 读取11位key及16位md5后key及5位起始seek和4位持续seek
			p0 = position
			p1 = p0 + 11
			p2 = p1 + 16
			p3 = p2 + 11
			p4 = p3 + 4
			hashKey := gnomon.ScaleDDuoStringToUint64(indexStr[p0:p1])
			md516Key := indexStr[p1:p2]
			seekStart := gnomon.ScaleDDuoStringToInt64(indexStr[p2:p3])     // value最终存储在文件中的起始位置
			seekLast := int(gnomon.ScaleDDuoStringToInt64(indexStr[p3:p4])) // value最终存储在文件中的持续长度
			//log.Debug("read", log.Field("i", i), log.Field("node", i.node))
			link, _ := i.node.put(md516Key, hashKey, hashKey)
			link.Fit(p0, seekStart, seekLast)
			// todo ID自增
			//atomic.AddUint64(i.form.getAutoID(), 1) // ID自增
		}(i, position)
		position += 42 // 单条索引默认占用长度42个字符
		if indexStrLen < position+42 {
			haveNext = false
		}
	}
	if nil == err {
		offset += int64(peekOnce)
		return i.read(file, offset)
	}
	return
}
