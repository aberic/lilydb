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

import (
	"github.com/aberic/lilydb/engine/comm"
	"reflect"
	"testing"
)

func TestHash(t *testing.T) {
	t.Log(comm.Hash("test"))
}

func TestPathFormFile(t *testing.T) {
	t.Log(PathFormFile("databaseID", "formID"))
}

func TestPathFormIndexFile(t *testing.T) {
	t.Log(PathFormIndexFile("databaseID", "formID", "indexID"))
}

func TestType2index(t *testing.T) {
	t.Log(Type2index(1))
	t.Log(Type2index(int64(-64)))
	t.Log(Type2index(uint64(9223372036854775808)))
	t.Log(Type2index(uint64(64)))
	t.Log(Type2index(6.4))
	t.Log(Type2index("test"))
	t.Log(Type2index(true))
	t.Log(Type2index(false))
	t.Log(Type2index(struct{}{}))
}

func TestValueType2index(t *testing.T) {
	reflectValue := reflect.ValueOf(1)
	t.Log(ValueType2index(&reflectValue))
	reflectValue = reflect.ValueOf(uint64(9223372036854775808))
	t.Log(ValueType2index(&reflectValue))
	reflectValue = reflect.ValueOf(-1)
	t.Log(ValueType2index(&reflectValue))
	reflectValue = reflect.ValueOf(1.5)
	t.Log(ValueType2index(&reflectValue))
	reflectValue = reflect.ValueOf("test")
	t.Log(ValueType2index(&reflectValue))
	reflectValue = reflect.ValueOf(true)
	t.Log(ValueType2index(&reflectValue))
	reflectValue = reflect.ValueOf(false)
	t.Log(ValueType2index(&reflectValue))
	reflectValue = reflect.ValueOf(struct{}{})
	t.Log(ValueType2index(&reflectValue))
}

func TestValue2hashKey(t *testing.T) {
	reflectValue := reflect.ValueOf(1)
	t.Log(Value2hashKey(&reflectValue))
	reflectValue = reflect.ValueOf(uint64(9223372036854775808))
	t.Log(Value2hashKey(&reflectValue))
	reflectValue = reflect.ValueOf(-1)
	t.Log(Value2hashKey(&reflectValue))
	reflectValue = reflect.ValueOf(1.5)
	t.Log(Value2hashKey(&reflectValue))
	reflectValue = reflect.ValueOf("test")
	t.Log(Value2hashKey(&reflectValue))
	reflectValue = reflect.ValueOf(true)
	t.Log(Value2hashKey(&reflectValue))
	reflectValue = reflect.ValueOf(false)
	t.Log(Value2hashKey(&reflectValue))
	reflectValue = reflect.ValueOf(struct{}{})
	t.Log(Value2hashKey(&reflectValue))
}
