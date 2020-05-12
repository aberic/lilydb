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

type link struct {
	md516Key       string
	seekStartIndex int64 // 索引最终存储在文件中的起始位置
	seekStart      int64 // value最终存储在文件中的起始位置
	seekLast       int   // value最终存储在文件中的持续长度
	version        int   // 当前索引数据版本号
}

// Fit 填充数据
//
// seekStartIndex 索引最终存储在文件中的起始位置
//
// seekStart value最终存储在文件中的起始位置
//
// seekLast value最终存储在文件中的持续长度
//
// version 当前索引数据版本号
func (l *link) Fit(seekStartIndex int64, seekStart int64, seekLast, version int) {
	l.seekStartIndex = seekStartIndex
	l.seekStart = seekStart
	l.seekLast = seekLast
}

func (l *link) MD516Key() string {
	return l.md516Key
}

func (l *link) SeekStartIndex() int64 {
	return l.seekStartIndex
}

func (l *link) SeekStart() int64 {
	return l.seekStart
}

func (l *link) SeekLast() int {
	return l.seekLast
}

func (l *link) Version() int {
	return l.version
}
