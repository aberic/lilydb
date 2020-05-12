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
	"errors"
	"github.com/aberic/gnomon"
	"github.com/aberic/lilydb/config"
	"os"
	"path/filepath"
)

var (
	// ErrDatabaseExist 自定义error信息
	ErrDatabaseExist = errors.New("database already exist")
	// ErrDataNotFound 自定义error信息
	ErrDataNotFound = errors.New("database not found")
	// ErrFormNotFound 自定义error信息
	ErrFormNotFound = errors.New("form not found")
	//// ErrFormExist 自定义error信息
	//ErrFormExist = errors.New("form already exist")
	//// ErrKeyExist 自定义error信息
	//ErrKeyExist = errors.New("key already exist")
	//// ErrIndexExist 自定义error信息
	//ErrIndexExist = errors.New("index already exist")
	//// ErrKeyIsNil 自定义error信息
	//ErrKeyIsNil = errors.New("put keyStructure can not be nil")
)

// mkDataDir 创建库存储目录
func mkDataDir(dataName string) (err error) {
	dataPath := filepath.Join(config.Obtain().DataDir, dataName)
	if gnomon.FilePathExists(dataPath) {
		return ErrDatabaseExist
	}
	return os.MkdirAll(dataPath, os.ModePerm)
}