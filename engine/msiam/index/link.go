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

// Link 叶子节点下的链表对象接口
type Link struct {
	key      string      // 存入key
	md516Key string      // md516后的key
	value    interface{} // 值
	version  int         // 当前索引数据版本号
}

// Fit 填充数据
//
// key 存入key
//
// md516Key md516后的key
//
// value 值
//
// version 当前索引数据版本号
func (l *Link) Fit(key, md516Key string, value interface{}, version int) {
	l.key = key
	l.md516Key = md516Key
	l.value = value
	l.version = version
}

// Key 存入key
func (l *Link) Key() string {
	return l.key
}

// MD516Key md516后的key
func (l *Link) MD516Key() string {
	return l.md516Key
}

// Value 值
func (l *Link) Value() interface{} {
	return l.value
}

// Version 当前索引数据版本号
func (l *Link) Version() int {
	return l.version
}
