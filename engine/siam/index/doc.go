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

// Package index 数据库存取方法，包含富查询条件选择器
//
// hash array 模型 [00, 01, 02, 03, 04, 05, 06, 07, 08, 09, a, b, c, d, e, f]
//
// b+tree 模型 degree=128;level=4;nodes=[degree^level]/(degree-1)=2113665;
//
// node 内范围控制数量 keyStructure=127
//
// tree 内范围控制数量 treeCount=nodes*keyStructure=268435455
//
// hash array 内范围控制数量 t*16=4294967280
//
// level1间隔 ld1=(treeCount+1)/128=2097152
//
// level2间隔 ld2=(16513*127+1)/128=16384
//
// level3间隔 ld3=(129*127+1)/128=128
//
// level4间隔 ld3=(1*127+1)/128=1
//
//////////////////////////////////////////////////////////////////////////////////////////
//
// 表索引树总数量为 nodes = 18446744073709551616
//
// nodes = max uint64 即 1<<64 = 18446744073709551616
//
// b+tree 模型 degree=65536;level=4;nodes=degree^level=18446744073709551616;
//
// node 内范围控制数量 nodeKeyStructure=65536
//
// leaf 内范围控制数量 leafKeyStructure=65537
//
// tree从上至下level1间隔 ld1=degree^(level-1)=65536^3=281474976710656
//
// tree从上至下level2间隔 ld1=degree^(level-1)=65536^2=4294967296
//
// tree从上至下level3间隔 ld1=degree^(level-1)=65536^1=65536
//
// tree从上至下level4间隔 ld1=degree^(level-1)=65536^0=1
//
// 存储格式 {dataDir}/database/{dataName}/{formName}/{formName}.dat/idx...
//
// 索引格式
package index
