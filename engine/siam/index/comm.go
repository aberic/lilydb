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
	"errors"
	"github.com/aberic/gnomon"
	"strconv"
)

const (
	// level1Distance level1间隔 65536^3 = 281474976710656 | 测试 4^3 = 64
	level1Distance uint64 = 281474976710656
	// level2Distance level2间隔 65536^2 = 4294967296 | 测试 4^2 = 16
	level2Distance uint64 = 4294967296
	// level3Distance level3间隔 65536^1 = 65536 | 测试 4^1 = 4
	level3Distance uint64 = 65536
	// level4Distance level4间隔 65536^0 = 1 | 测试 4^0 = 1
	level4Distance uint64 = 1
)

var (
	// ErrIndexFileNotFound 自定义error信息
	ErrIndexFileNotFound = errors.New("index file not found")
)

// levelDistance 根据节点所在层级获取当前节点内部子节点之间的差
func levelDistance(level uint8) uint64 {
	switch level {
	case 1:
		return level1Distance
	case 2:
		return level2Distance
	case 3:
		return level3Distance
	case 4:
		return level4Distance
	}
	return 0
}

// conditionValueInt64 判断当前条件是否满足
func conditionValueInt64(cond string, paramType paramType, paramValue interface{}, value int64) bool {
	var (
		paramInt64Value int64
	)
	if paramType != paramInt64 {
		var success bool
		paramInt64Value, success = param2Int64(paramType, paramValue)
		if !success {
			return false
		}
	} else {
		paramInt64Value = paramValue.(int64)
	}
	switch cond {
	default:
		return false
	case "gt":
		return value > paramInt64Value
	case "lt":
		return value < paramInt64Value
	case "eq":
		return value == paramInt64Value
	case "dif":
		return value != paramInt64Value
	}
}

// conditionValueUint64 判断当前条件是否满足
func conditionValueUint64(cond string, paramType paramType, paramValue interface{}, value uint64) bool {
	var paramUint64Value uint64
	if paramType != paramUint64 {
		var success bool
		paramUint64Value, success = param2UInt64(paramType, paramValue)
		if !success {
			return false
		}
	} else {
		paramUint64Value = paramValue.(uint64)
	}
	switch cond {
	default:
		return false
	case "gt":
		return value > paramUint64Value
	case "lt":
		return value < paramUint64Value
	case "eq":
		return value == paramUint64Value
	case "dif":
		return value != paramUint64Value
	}
}

// conditionValueFloat64 判断当前条件是否满足
func conditionValueFloat64(cond string, paramType paramType, paramValue interface{}, value float64) bool {
	var paramFloat64Value float64
	if paramType != paramFloat64 {
		var success bool
		paramFloat64Value, success = param2Float64(paramType, paramValue)
		if !success {
			return false
		}
	} else {
		paramFloat64Value = paramValue.(float64)
	}
	switch cond {
	default:
		return false
	case "gt":
		return value > paramFloat64Value
	case "lt":
		return value < paramFloat64Value
	case "eq":
		return value == paramFloat64Value
	case "dif":
		return value != paramFloat64Value
	}
}

// conditionValueString 判断当前条件是否满足
func conditionValueString(cond string, paramType paramType, paramValue interface{}, value string) bool {
	var paramStringValue string
	if paramType != paramString {
		var success bool
		paramStringValue, success = param2String(paramType, paramValue)
		if !success {
			return false
		}
	} else {
		paramStringValue = paramValue.(string)
	}
	switch cond {
	default:
		return false
	case "gt":
		return value > paramStringValue
	case "lt":
		return value < paramStringValue
	case "eq":
		return value == paramStringValue
	case "dif":
		return value != paramStringValue
	}
}

// conditionValueBool 判断当前条件是否满足
func conditionValueBool(cond string, param, value bool) bool {
	switch cond {
	default:
		return false
	case "eq":
		return value == param
	case "dif":
		return value != param
	}
}

func param2Int64(paramType paramType, paramValue interface{}) (paramInt int64, trans bool) {
	switch paramType {
	default:
		return 0, false
	case paramInt64:
		return paramValue.(int64), true
	case paramUint64:
		uint64Str := strconv.FormatUint(paramValue.(uint64), 10) // uint64转string
		intNum, err := strconv.Atoi(uint64Str)                   // string转int
		if nil != err {
			return 0, false
		}
		return int64(intNum), true
	case paramFloat64:
		return gnomon.ScaleFloat64toInt64(paramValue.(float64), 0), true
	case paramString:
		intNum, err := strconv.Atoi(paramValue.(string)) // string转int
		if nil != err {
			return 0, false
		}
		return int64(intNum), true
	}
}

func param2UInt64(paramType paramType, paramValue interface{}) (paramInt uint64, trans bool) {
	switch paramType {
	default:
		return 0, false
	case paramInt64:
		int64Str := strconv.FormatInt(paramValue.(int64), 10)
		intNum, err := strconv.Atoi(int64Str)
		if nil != err {
			return 0, false
		}
		return uint64(intNum), true
	case paramUint64:
		return paramValue.(uint64), true
	case paramFloat64:
		return param2UInt64(paramInt64, gnomon.ScaleFloat64toInt64(paramValue.(float64), 0))
	case paramString:
		intNum, err := strconv.Atoi(paramValue.(string)) // string转int
		if nil != err {
			return 0, false
		}
		return uint64(intNum), true
	}
}

func param2Float64(paramType paramType, paramValue interface{}) (paramInt float64, trans bool) {
	switch paramType {
	default:
		return 0, false
	case paramInt64:
		return gnomon.ScaleInt64toFloat64(paramValue.(int64), 0), true
	case paramUint64:
		uint64Str := strconv.FormatUint(paramValue.(uint64), 10) // uint64转string
		return param2Float64(paramString, uint64Str)
	case paramFloat64:
		return paramValue.(float64), true
	case paramString:
		intNum, err := strconv.Atoi(paramValue.(string)) // string转int
		if nil != err {
			return 0, false
		}
		return param2Float64(paramInt64, int64(intNum))
	}
}

func param2String(paramType paramType, paramValue interface{}) (paramInt string, trans bool) {
	switch paramType {
	default:
		return "", false
	case paramInt64:
		return strconv.FormatInt(paramValue.(int64), 10), true
	case paramUint64:
		return strconv.FormatUint(paramValue.(uint64), 10), true
	case paramFloat64:
		return strconv.FormatFloat(paramValue.(float64), 'E', -1, 64), true
	case paramString:
		return paramValue.(string), true
	}
}
