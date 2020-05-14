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
func conditionValueInt64(cond string, param, value int64) bool {
	switch cond {
	default:
		return false
	case "gt":
		return value > param
	case "lt":
		return value < param
	case "eq":
		return value == param
	case "dif":
		return value != param
	}
}

// conditionValueUint64 判断当前条件是否满足
func conditionValueUint64(cond string, param, value uint64) bool {
	switch cond {
	default:
		return false
	case "gt":
		return value > param
	case "lt":
		return value < param
	case "eq":
		return value == param
	case "dif":
		return value != param
	}
}

// conditionValueFloat64 判断当前条件是否满足
func conditionValueFloat64(cond string, param, value float64) bool {
	switch cond {
	default:
		return false
	case "gt":
		return value > param
	case "lt":
		return value < param
	case "eq":
		return value == param
	case "dif":
		return value != param
	}
}

// conditionValueString 判断当前条件是否满足
func conditionValueString(cond string, param, value string) bool {
	switch cond {
	default:
		return false
	case "gt":
		return value > param
	case "lt":
		return value < param
	case "eq":
		return value == param
	case "dif":
		return value != param
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
