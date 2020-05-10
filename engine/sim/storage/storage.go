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

package storage

import (
	"bufio"
	"errors"
	"github.com/aberic/gnomon"
	"github.com/aberic/gnomon/log"
	"github.com/vmihailenco/msgpack"
	"github/aberic/lilydb/comm"
	"github/aberic/lilydb/config"
	"io"
	"os"
	"strings"
	"sync"
)

var (
	stg         *Storage
	onceStorage sync.Once

	// ErrValueType value type error
	ErrValueType = errors.New("value type error")
	// ErrValueInvalid value is invalid
	ErrValueInvalid = errors.New("value is invalid")
)

func Obtain() *Storage {
	onceStorage.Do(func() {
		if nil == stg {
			stg = &Storage{limitOpenFileChan: make(chan struct{}, config.Obtain().LimitOpenFile)}
		}
	})
	return stg
}

// Storage 文件存储服务
type Storage struct {
	engine            *engine
	limitOpenFileChan chan struct{} // limitOpenFileChan 限制打开文件描述符次数
}

// Take 取出具体内容
//
// filePath 从该路径文件中读取
//
// seekStart value最终存储在文件中的起始位置
//
// seekLast value最终存储在文件中的持续长度
//
// return Read 返回数据读取结果
func (s *Storage) Take(filePath string, seekStart int64, seekLast int) (*Read, error) {
	var (
		file *os.File
		err  error
	)
	defer func() {
		if nil != file {
			<-s.limitOpenFileChan
			_ = file.Close()
		}
	}()
	//log.Debug("read", log.Field("filePath", filePath), log.Field("seekStart", seekStart), log.Field("seekLast", seekLast))
	if file, err = s.openFile(filePath, os.O_RDONLY); err != nil {
		//log.Error("read", log.Err(err))
		return nil, err
	}
	_, err = file.Seek(seekStart, io.SeekStart) //表示文件的起始位置，从第seekStart个字符往后读取
	if err != nil {
		//log.Error("read", log.Err(err))
		return nil, err
	}
	inputReader := bufio.NewReader(file)
	var bytes []byte
	if bytes, err = inputReader.Peek(seekLast); nil != err {
		//log.Error("read", log.Err(err))
		return nil, err
	}
	var value interface{}
	if err = msgpack.Unmarshal(bytes, &value); nil != err {
		//log.Error("read", log.Err(err))
		return nil, err
	}
	switch value.(type) {
	default:
		return nil, ErrValueType
	case map[string]interface{}:
		valueMap := value.(map[string]interface{})
		if valueMap["I"].(bool) {
			return &Read{Key: valueMap["K"].(string), Value: valueMap["V"]}, nil
		}
		return nil, ErrValueInvalid
	}
}

// Store 存储具体内容
//
// databaseID 数据库唯一id
//
// formID 表唯一id
//
// key 唯一key
//
// value 存储具体内容
//
// valid 存储有效性，如无效则表示该记录不可用，即删除
//
// writes 索引即将写入的参考坐标数组
//
// return []*Written 返回完成索引写入后的结果数组
func (s *Storage) Store(databaseID, formID, key string, value interface{}, valid bool, writes []*Write) ([]*Written, error) {
	var (
		formFilePath = comm.PathFormDataFile(databaseID, formID) // path 存储文件路径
		file         *os.File
		seekStart    int64
		seekLast     int
		data         []byte
		err          error
	)
	// 存储数据外包装数据属性
	vd := &valueData{K: key, I: valid, V: value}
	if data, err = msgpack.Marshal(vd); nil != err {
		return nil, err
	}
	fm := s.engine.form(databaseID, formID, formFilePath)
	defer fm.mu.Unlock()
	fm.mu.Lock()
	if nil == fm.file {
		if file, err = s.openFile(formFilePath, os.O_CREATE|os.O_RDWR|os.O_APPEND); nil != err {
			log.Error("storeData", log.Err(err))
			<-s.limitOpenFileChan
			return nil, err
		}
		fm.file = file
	}
	// value最终存储在文件中的起始位置
	if seekStart, err = file.Seek(0, io.SeekEnd); err != nil {
		log.Debug("storeData", log.Err(err))
		return nil, err
	}
	if seekLast, err = file.Write(data); nil != err { // value最终存储在文件中的持续长度
		log.Debug("storeData", log.Err(err))
		return nil, err
	}
	var (
		wg          sync.WaitGroup
		writtenHelp = newWrittenHelp()
	)
	for _, write := range writes {
		wg.Add(1)
		go func(databaseID, formID, formFilePath string, form *form, seekStart int64, seekLast int, write *Write, writtenHelp *WrittenHelp) {
			defer wg.Done()
			written, err := s.storeIndex(databaseID, formID, formFilePath, form, seekStart, seekLast, write)
			if nil == err {
				writtenHelp.mu.Lock()
				writtenHelp.writtenArr = append(writtenHelp.writtenArr, written)
				writtenHelp.mu.Unlock()
			} else {
				writtenHelp.err = err
			}
		}(databaseID, formID, formFilePath, fm, seekStart, seekLast, write, writtenHelp)
	}
	wg.Wait()
	if nil != writtenHelp.err {
		return nil, writtenHelp.err
	}
	return writtenHelp.writtenArr, nil
}

// storeIndex 存储索引文件
//
// databaseID 数据库唯一id
//
// formID 表唯一id
//
// form 表级操作结构
//
// seekStart value最终存储在文件中的起始位置
//
// seekLast value最终存储在文件中的持续长度
//
// write 索引即将写入的参考坐标
//
// return Written 返回完成索引写入后的结果
func (s *Storage) storeIndex(databaseID, formID, formFilePath string, form *form, seekStart int64, seekLast int, write *Write) (*Written, error) {
	var (
		file *os.File
		err  error
	)
	idx := s.engine.index(databaseID, formID, write.IndexID, formFilePath, write.FormIndexFilePath)
	defer idx.mu.Unlock()
	idx.mu.Lock()
	if nil == idx.file {
		// 将获取到的索引存储位置传入。如果为0，则表示没有存储过；如果不为0，则覆盖旧的存储记录
		if file, err = s.openFile(write.FormIndexFilePath, os.O_CREATE|os.O_RDWR); nil != err {
			log.Error("storeIndex", log.Err(err))
			<-s.limitOpenFileChan
			return nil, err
		}
		idx.file = file
	}
	md5Key := gnomon.HashMD516(write.Key) // hash(keyStructure) 会发生碰撞，因此这里存储md5结果进行反向验证
	// 写入11位key及16位md5后key
	appendStr := strings.Join([]string{gnomon.StringPrefixSupplementZero(gnomon.ScaleUint64ToDDuoString(write.HashKey), 11), md5Key}, "")
	//log.Debug("storeIndex",
	//	log.Field("appendStr", appendStr),
	//	log.Field("formIndexFilePath", ib.getFormIndexFilePath()),
	//	log.Field("seekStartIndex", ib.getLink().getSeekStartIndex()))
	var seekEnd int64
	//log.Debug("running", log.Field("type", "moldIndex"), log.Field("seekStartIndex", it.link.getSeekStartIndex()))
	if write.SeekStartIndex == -1 {
		if seekEnd, err = file.Seek(0, io.SeekEnd); nil != err {
			log.Error("storeIndex", log.Err(err))
			return nil, err
		}
		//log.Debug("running", log.Field("it.link.seekStartIndex == -1", seekEnd))
	} else {
		if seekEnd, err = file.Seek(write.SeekStartIndex, io.SeekStart); nil != err { // 寻址到原索引起始位置
			log.Error("storeIndex", log.Err(err))
			return nil, err
		}
		//log.Debug("running", log.Field("seekStartIndex", it.link.getSeekStartIndex()), log.Field("it.link.seekStartIndex != -1", seekEnd))
	}
	// 写入11位key及16位md5后key及5位起始seek和4位持续seek
	if _, err = file.WriteString(strings.Join([]string{appendStr,
		gnomon.StringPrefixSupplementZero(gnomon.ScaleInt64ToDDuoString(seekStart), 11),
		gnomon.StringPrefixSupplementZero(gnomon.ScaleIntToDDuoString(seekLast), 4)}, "")); nil != err {
		//log.Error("running", log.Field("seekStartIndex", seekEnd), log.Err(err))
		return nil, err
	}
	//log.Debug("storeIndex", log.Field("ib.getKey()", ib.getKey()), log.Field("md516Key", md516Key), log.Field("seekStartIndex", wf.seekStartIndex))
	//log.Debug("running", log.Field("it.link.seekStartIndex", seekEnd), log.Err(err))
	return &Written{MD516Key: md5Key, SeekStartIndex: seekEnd, SeekStart: seekStart, SeekLast: seekLast}, nil
}

func (s *Storage) openFile(filePath string, flag int) (file *os.File, err error) {
	s.limitOpenFileChan <- struct{}{}
	return os.OpenFile(filePath, flag, 0644)
}
