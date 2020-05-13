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
	"context"
	"errors"
	"github.com/aberic/gnomon"
	"github.com/aberic/gnomon/log"
	"github.com/aberic/lilydb/connector"
	"github.com/aberic/lilydb/engine/siam/utils"
	"io"
	"os"
	"sync"
	"sync/atomic"
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
func (i *Index) Put(md516Key string, hashKey uint64, version int) (link connector.Link, exist, versionGT bool) {
	return i.node.put(md516Key, hashKey, hashKey, version)
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
func (i *Index) Recover() (autoID *uint64, err error) {
	indexFilePath := utils.PathFormIndexFile(i.databaseID, i.formID, i.id)
	if gnomon.FilePathExists(indexFilePath) { // 索引文件存在才继续恢复
		var (
			indexContentSize int64  // 索引文件长度
			pieceCount       int64  // 将索引分块后的数量
			id               uint64 = 0
		)
		autoID = &id
		// 获取索引文件长度
		if indexContentSize, err = i.indexFileSize(indexFilePath); nil != err {
			return
		}
		// 获取将索引分块后的数量，20/18=1
		pieceCount = indexContentSize / utils.LenPeekOnce64
		var (
			wg          sync.WaitGroup
			offset      int64
			ctx, cancel = context.WithCancel(context.Background())
		)
		defer cancel()                                   // 确保所有路径都取消了上下文，以避免上下文泄漏
		for offset = pieceCount; offset >= 0; offset-- { // 将索引进行分块恢复，倒序恢复能尽量避免重复赋值，提升效率
			wg.Add(1)
			go func(ctx context.Context, indexFilePath string, offset int64) {
				defer wg.Done()
				if err := i.read(ctx, autoID, indexFilePath, offset*utils.LenPeekOnce64); nil != err {
					cancel()
				}
			}(ctx, indexFilePath, offset)
		}
		wg.Wait()
		return
	}
	return nil, ErrIndexFileNotFound
}

func (i *Index) read(ctx context.Context, autoID *uint64, indexFilePath string, offset int64) (err error) {
	var (
		data    []byte
		success = make(chan struct{})
	)
	if data, err = i.readAppointData(indexFilePath, offset); nil != err {
		if io.EOF == err {
			if len(data) == 0 {
				return
			}
			if len(data)%utils.LenIndex != 0 { // 单条索引默认占用长度
				return errors.New("index lens does't match")
			}
		} else {
			return
		}
	}
	go func(data []byte) {
		var (
			wg          sync.WaitGroup
			haveNext          = true
			indexStr          = string(data)
			indexStrLen       = int64(len(indexStr))
			position    int64 = 0
		)
		for haveNext {
			wg.Add(1)
			go func(position int64, indexStr string) {
				defer wg.Done()
				// 恢复索引中link数据，同时对路径上的node进行恢复
				i.recoverLink(autoID, position, indexStr)
			}(position, indexStr)
			position += utils.LenIndex64 // 单条索引默认占用长度
			if indexStrLen < position+utils.LenIndex64 {
				haveNext = false
			}
		}
		wg.Wait()
		success <- struct{}{}
	}(data)
	for {
		select {
		case <-ctx.Done(): // 上下文停止，直接返回
			return
		case <-success: // 处理完成，返回
			return
		}
	}
}

// recoverLink 恢复索引中link数据，同时对路径上的node进行恢复
func (i *Index) recoverLink(autoID *uint64, position int64, indexStr string) {
	var p0, p1, p2, p3, p4, p5 int64
	// 读取 11位hashKey + 16位md5Key + 11位起始seek + 4位持续seek + 4位版本号 = 46
	p0 = position
	p1 = p0 + utils.LenHashKey
	p2 = p1 + utils.LenMD5Key
	p3 = p2 + utils.LenSeekStart
	p4 = p3 + utils.LenSeekLast
	p5 = p4 + utils.LenVersion
	hashKey := gnomon.ScaleDDuoStringToUint64(indexStr[p0:p1])
	md516Key := indexStr[p1:p2]
	seekStart := gnomon.ScaleDDuoStringToInt64(indexStr[p2:p3])     // value最终存储在文件中的起始位置
	seekLast := int(gnomon.ScaleDDuoStringToInt64(indexStr[p3:p4])) // value最终存储在文件中的持续长度
	version := int(gnomon.ScaleDDuoStringToInt64(indexStr[p4:p5]))
	//log.Debug("read", log.Field("i", i), log.Field("node", i.node))
	link, _, versionGT := i.node.put(md516Key, hashKey, hashKey, version)
	if versionGT {
		link.Fit(p0, seekStart, seekLast, version)
		atomic.AddUint64(autoID, 1) // ID自增
	}
}

// readAppointData 读取指定文件中指定起始位置的指定长度的内容
func (i *Index) readAppointData(indexFilePath string, offset int64) (data []byte, err error) {
	var (
		file        *os.File
		inputReader *bufio.Reader
	)
	if file, err = os.OpenFile(indexFilePath, os.O_RDONLY, 0644); nil != err {
		log.Error("index recover multi read failed", log.Err(err))
		return
	}

	// 将file下标移动到指定位置作为读取数据的起始位置
	if _, err = file.Seek(offset, io.SeekStart); nil != err { //表示文件的起始位置，从第二个字符往后写入。
		return
	}
	inputReader = bufio.NewReaderSize(file, utils.LenPeekOnce) // 将 rd 封装成一个拥有 size 大小缓存的 bufio.Reader 对象
	return inputReader.Peek(utils.LenPeekOnce)                 // 返回缓存的一个切片，该切片引用缓存中前 n 字节数据
}

// indexFileSize 获取索引文件长度
func (i *Index) indexFileSize(indexFilePath string) (bufLen int64, err error) {
	var file *os.File
	defer func() { _ = file.Close() }()
	if file, err = os.OpenFile(indexFilePath, os.O_RDONLY, 0644); nil == err {
		//文件指针指向文件末尾 获取文件大小保存于bufLen
		return file.Seek(0, io.SeekEnd)
	}
	log.Error("index recover multi read failed", log.Err(err))
	return
}
