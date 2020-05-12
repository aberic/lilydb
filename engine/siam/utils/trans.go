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
	"github.com/aberic/gnomon"
	"reflect"
	"strconv"
)

// Type2index 通过存储内容获取索引信息
func Type2index(value interface{}) (key string, hashKey uint64, support bool) {
	support = true
	switch value := value.(type) {
	default:
		return "", 0, false
	case int:
		i64, err := strconv.ParseInt(strconv.Itoa(value), 10, 64)
		if nil != err {
			return "", 0, false
		}
		key = strconv.FormatInt(i64, 10)
		hashKey = uint64(i64 + 9223372036854775807 + 1)
	case int8, int16, int32, int64:
		i64 := value.(int64)
		key = strconv.FormatInt(i64, 10)
		hashKey = uint64(i64 + 9223372036854775807 + 1)
	case uint8, uint16, uint32, uint, uint64, uintptr:
		ui64 := value.(uint64)
		key = strconv.FormatUint(ui64, 10)
		if ui64 > 9223372036854775807 { // 9223372036854775808 = 1 << 63
			return "", 0, false
		}
		hashKey = ui64 + 9223372036854775807 + 1
	case float32, float64:
		i64 := gnomon.ScaleFloat64toInt64(value.(float64), 4)
		key = strconv.FormatInt(i64, 10)
		hashKey = uint64(i64 + 9223372036854775807 + 1)
	case string:
		hashKey = Hash(value)
	case bool:
		if value {
			key = "true"
			hashKey = 1
		} else {
			key = "false"
			hashKey = 2
		}
	}
	return
}

// ValueType2index 通过存储内容获取索引信息
func ValueType2index(value *reflect.Value) (key string, hashKey uint64, support bool) {
	support = true
	switch value.Kind() {
	default:
		return "", 0, false
	case reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int, reflect.Int64:
		i64 := value.Int()
		key = strconv.FormatInt(i64, 10)
		hashKey = uint64(i64 + 9223372036854775807 + 1)
	case reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint, reflect.Uint64, reflect.Uintptr:
		ui64 := value.Uint()
		key = strconv.FormatUint(ui64, 10)
		if ui64 > 9223372036854775807 { // 9223372036854775808 = 1 << 63
			return "", 0, false
		}
		hashKey = ui64 + 9223372036854775807 + 1
	case reflect.Float32, reflect.Float64:
		i64 := gnomon.ScaleFloat64toInt64(value.Float(), 4)
		key = strconv.FormatInt(i64, 10)
		hashKey = uint64(i64 + 9223372036854775807 + 1)
	case reflect.String:
		key = value.String()
		hashKey = Hash(key)
	case reflect.Bool:
		if value.Bool() {
			key = value.String()
			hashKey = 1
		} else {
			key = value.String()
			hashKey = 2
		}
	}
	return
}

// Value2hashKey 通过存储内容获取索引信息
func Value2hashKey(value *reflect.Value) (hashKey uint64, support bool) {
	support = true
	switch value.Kind() {
	default:
		return 0, false
	case reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int, reflect.Int64:
		i64 := value.Int()
		hashKey = uint64(i64 + 9223372036854775807 + 1)
	case reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint, reflect.Uint64, reflect.Uintptr:
		ui64 := value.Uint()
		if ui64 > 9223372036854775808 { // 9223372036854775808 = 1 << 63
			return 0, false
		}
		hashKey = ui64 + 9223372036854775807 + 1
	case reflect.Float32, reflect.Float64:
		hashKey = uint64(gnomon.ScaleFloat64toInt64(value.Float(), 4) + 9223372036854775807 + 1)
	case reflect.String:
		hashKey = Hash(value.String())
	case reflect.Bool:
		if value.Bool() {
			hashKey = 1
		} else {
			hashKey = 2
		}
	}
	return
}
