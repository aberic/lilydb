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

package utils

const (
	// LenHashKey 11位hashKey
	LenHashKey = 11
	// LenMD5Key 16位md5Key
	LenMD5Key = 16
	// LenSeekStart 11位起始seek
	LenSeekStart = 11
	// LenSeekLast 4位持续seek
	LenSeekLast = 4
	// LenVersion 4位版本号
	LenVersion = 4
	// LenIndex 单条索引长度 = 46
	LenIndex = 46
	// LenIndex64 单条索引长度 = 46
	LenIndex64 int64 = 46
	// LenPeekOnce 单次恢复索引长度 = 46000
	LenPeekOnce = 46000
	// LenPeekOnce64 单次恢复索引长度 = 46000
	LenPeekOnce64 int64 = 46000
)
