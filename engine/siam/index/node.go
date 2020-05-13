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
	"sort"
	"strings"
	"sync"
)

// node 手提袋
//
// 这里面能存放很多个包装盒
//
// box 包装盒集合
//
// 存储格式 {dataDir}/database/{dataName}/{formName}.form...
type node struct {
	level       uint8  // 当前节点所在树层级
	degreeIndex uint16 // 当前节点所在集合中的索引下标，该坐标不一定在数组中的正确位置，但一定是逻辑正确的
	preNode     *node  // node 所属 trolley
	nodes       []*node
	links       []*link
	mu          sync.RWMutex
}

// put 预插入数据
//
// md516Key md516Key，必须string类型
//
// hashKey 索引key，可通过hash转换string生成
//
// flexibleKey 下一级最左最小树所对应真实key
//
// version 当前索引数据版本号
//
// return exist 返回是否存在
func (n *node) put(md516Key string, hashKey, flexibleKey uint64, version int) (link *link, exist, versionGT bool) {
	var (
		nextDegree      uint16 // 下一节点所在当前节点下度的坐标
		nextFlexibleKey uint64 // 下一级最左最小树所对应真实key
		distance        uint64 // 指定Level层级节点内各个子节点之前的差
		nd              *node
	)
	if n.level < 5 {
		distance = levelDistance(n.level)
		//gnomon.Log().Debug("put", gnomon.Log().Field("key", key), gnomon.Log().Field("distance", distance))
		nextDegree = uint16(flexibleKey / distance)
		nextFlexibleKey = flexibleKey - uint64(nextDegree)*distance
		if n.level == 4 {
			nd = n.createOrTakeLeaf(nextDegree) // 创建或获取下一个叶子节点
		} else {
			nd = n.createOrTakeNode(nextDegree) // 创建或获取下一个子节点
		}
	} else {
		return n.link(md516Key, version)
	}
	return nd.put(md516Key, hashKey, nextFlexibleKey, version)
}

// get 获取数据，返回存储数据的可检索对象
//
// md516Key md516Key，必须string类型
//
// hashKey 索引key，可通过hash转换string生成
//
// flexibleKey 下一级最左最小树所对应真实key
func (n *node) get(md516Key string, hashKey, flexibleKey uint64) *link {
	var (
		nextDegree      uint16 // 下一节点所在当前节点下度的坐标
		nextFlexibleKey uint64 // 下一级最左最小树所对应真实key
		distance        uint64 // 指定Level层级节点内各个子节点之前的差
	)
	if n.level < 5 {
		distance = levelDistance(n.level)
		nextDegree = uint16(flexibleKey / distance)
		nextFlexibleKey = flexibleKey - uint64(nextDegree)*distance
	} else {
		//gnomon.Log().Debug("box-get", gnomon.Log().Field("key", key))
		if realIndex, exist := n.existLink(md516Key); exist {
			return n.links[realIndex]
		}
		return nil
	}
	if realIndex, err := n.existNode(nextDegree); nil == err {
		return n.nodes[realIndex].get(md516Key, hashKey, nextFlexibleKey)
	}
	return nil
}

func (n *node) existNode(index uint16) (realIndex int, err error) {
	return n.binaryMatchData(index)
}

func (n *node) createOrTakeNode(index uint16) *node {
	var (
		realIndex int
		err       error
	)
	defer n.mu.Unlock()
	n.mu.Lock()
	if realIndex, err = n.existNode(index); nil != err {
		level := n.level + 1
		newNode := &node{
			level:       level,
			degreeIndex: index,
			preNode:     n,
			nodes:       []*node{},
		}
		return n.appendNodal(newNode)
	}
	return n.nodes[realIndex]
}

func (n *node) createOrTakeLeaf(index uint16) *node {
	var (
		realIndex int
		err       error
	)
	defer n.mu.Unlock()
	n.mu.Lock()
	if realIndex, err = n.binaryMatchData(index); nil != err {
		level := n.level + 1
		leaf := &node{
			level:       level,
			degreeIndex: index,
			preNode:     n,
			links:       []*link{},
		}
		return n.appendNodal(leaf)
	}
	return n.nodes[realIndex]
}

// link 获取link叶子，如果存在，则判断版本号，如果不存在，则新建一个空并返回
//
// md516Key 索引md516Key
//
// version 当前索引数据版本号
func (n *node) link(md516Key string, version int) (lk *link, exist, versionGT bool) {
	if pos, exist := n.existLink(md516Key); exist {
		lk = n.links[pos]
		if version > lk.version {
			return n.links[pos], true, true
		}
		return n.links[pos], true, false
	}
	defer n.mu.Unlock()
	n.mu.Lock()
	lk = &link{md516Key: md516Key, version: version}
	n.links = append(n.links, lk)
	return lk, false, true
}

func (n *node) existLink(md516Key string) (int, bool) {
	defer n.mu.RUnlock()
	n.mu.RLock()
	for index, link := range n.links {
		//gnomon.Log().Debug("existLink", gnomon.Log().Field("link.md516Key", link.getMD516Key()), gnomon.Log().Field("md516", gnomon.CryptoHash().MD516(key)))
		if strings.EqualFold(link.MD516Key(), md516Key) {
			return index, true
		}
	}
	return 0, false
}

func (n *node) appendNodal(node *node) *node {
	nodesLen := len(n.nodes)
	if nodesLen == 0 {
		n.nodes = append(n.nodes, node)
		return node
	}
	n.nodes = append(n.nodes, node)
	sort.Slice(n.nodes, func(i, j int) bool {
		return n.nodes[i].degreeIndex < n.nodes[j].degreeIndex
	})
	//for i := nodesLen - 2; i >= 0; i-- {
	//	if n.nodes[i].degreeIndex < index {
	//		n.nodes[i+1] = node
	//		break
	//	} else if n.nodes[i].degreeIndex > index {
	//		n.nodes[i+1] = n.nodes[i]
	//		n.nodes[i] = node
	//	} else {
	//		return n.nodes[i]
	//	}
	//}
	return node
}

// binaryMatchData 'Nodal'内子节点数组二分查找基本方法
//
// matchIndex 要查找的值
//
// data 'binaryMatcher'接口支持的获取‘Nodal’接口的内置方法对象
//
// realIndex 返回查找到的真实的元素下标，该下标是对应数组内的下标，并非树中节点数组原型的下标
//
// 如果没找到，则返回err
func (n *node) binaryMatchData(matchIndex uint16) (realIndex int, err error) {
	var (
		leftIndex   int
		middleIndex int
		rightIndex  int
	)
	leftIndex = 0
	nodes := n.nodes
	rightIndex = len(nodes) - 1
	for leftIndex <= rightIndex {
		middleIndex = (leftIndex + rightIndex) / 2
		// 如果要找的数比midVal大
		if nodes[middleIndex].degreeIndex > matchIndex {
			// 在arr数组的左边找
			rightIndex = middleIndex - 1
		} else if nodes[middleIndex].degreeIndex < matchIndex {
			// 在arr数组的右边找
			leftIndex = middleIndex + 1
		} else if nodes[middleIndex].degreeIndex == matchIndex {
			return middleIndex, nil
		}
	}
	return 0, errors.New("index is nil")
}
