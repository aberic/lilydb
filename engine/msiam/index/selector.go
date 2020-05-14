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
	"encoding/json"
	"github.com/aberic/gnomon/log"
	"github.com/aberic/lilydb/engine/siam/utils"
	"reflect"
	"strings"
)

// NewSelector 新建检索选择器
//
// selectorBytes 选择器字节数组，自定义转换策略
//
// indexes 指定表下索引集合
//
// databaseID 数据库唯一ID
//
// formID 表唯一ID
//
// delete 是否删除检索结果
func NewSelector(selectorBytes []byte, indexes []*Index, databaseID, formID string, delete bool) (*Selector, error) {
	selector := &Selector{}
	if err := json.Unmarshal(selectorBytes, selector); nil != err {
		return nil, err
	}
	selector.indexes = indexes
	selector.databaseID = databaseID
	selector.formID = formID
	selector.delete = delete
	return selector, nil
}

// Selector 检索选择器
//
// 查询顺序 scope -> match -> conditions -> skip -> sort -> limit
type Selector struct {
	indexes    []*Index     // indexes 指定表下索引集合
	Conditions []*condition `json:"Conditions"` // Conditions 条件查询
	Skip       uint32       `json:"Skip"`       // Skip 结果集跳过数量
	Sort       *rank        `json:"Sort"`       // Sort 排序方式
	Limit      uint32       `json:"Limit"`      // Limit 结果集顺序数量
	databaseID string       // 数据库唯一ID
	formID     string       // 表唯一ID
	delete     bool         // 是否删除检索结果
}

// Run 执行富查询
//
// return count 检索结果总条数
//
// return values 检索结果集合
//
// return err 检索错误信息，如果有
func (s *Selector) Run() (int32, []interface{}) {
	idx, asc, nc, pcs := s.index() // 根据检索条件获取使用索引对象等信息
	log.Debug("query", log.Field("index", idx.KeyStructure()))
	if s.Limit == 0 { // 默认限制查询1000条数据
		s.Limit = 1000
	}
	if asc { // 是否顺序查询
		return s.leftQueryIndex(idx, nc, pcs)
	}
	return s.rightQueryIndex(idx, nc, pcs)
}

// getIndex 根据检索条件获取使用索引对象
//
// index 已获取索引对象
//
// asc 是否顺序查询
//
// cond 条件对象
func (s *Selector) index() (index *Index, asc bool, nc *nodeCondition, pcs map[string]*paramCondition) {
	// 优先尝试采用条件作为索引，缩小索引范围以提高检索效率
	index, asc, nc, pcs = s.indexCondition()
	if index != nil { // 如果存在条件查询，则优先条件查询
		return
	}
	for _, index = range s.indexes { // 如果存在排序查询，则优先排序查询
		if nil != index && s.Sort != nil && s.Sort.Param == index.KeyStructure() {
			return
		}
	}
	// 取值默认索引来进行查询操作
	index = s.indexes[0]
	return
}

// indexCondition 优先尝试采用条件作为索引，缩小索引范围以提高检索效率
//
// 优先匹配有多个相同Param参数的条件，如果相同数量一样，则按照先后顺序选择最先匹配的
func (s *Selector) indexCondition() (index *Index, asc bool, nc *nodeCondition, pcs map[string]*paramCondition) {
	pcs = make(map[string]*paramCondition)
	ncs := make(map[string]*nodeCondition)
	asc = true
	for _, condition := range s.Conditions { // 遍历检索条件
		for _, idx := range s.indexes { // 遍历检索索引
			if condition.Param == idx.KeyStructure() { // 匹配条件是否存在已有索引，如没有，进入下一轮循环
				if nil != s.Sort && s.Sort.Param == idx.KeyStructure() { // 如果有，则继续判断该索引是否存在排序需求
					index = idx
					asc = s.Sort.ASC
					if ncs[condition.Param] == nil {
						ncs[condition.Param] = &nodeCondition{nss: []*nodeSelector{}}
					}
					s.conditionNode(ncs[condition.Param], condition)
					break
				}
				if index != nil {
					continue
				}
				index = idx
				if ncs[condition.Param] == nil {
					ncs[condition.Param] = &nodeCondition{nss: []*nodeSelector{}}
				}
				s.conditionNode(ncs[condition.Param], condition)
			}
		}

		paramType, paramValue, support := s.formatParam(condition.Value)
		if support {
			pcs[s.pcMapName(condition)] = &paramCondition{paramType: paramType, paramValue: paramValue}
		}
	}
	var nodeCount = 0
	for _, ncNow := range ncs {
		ncNowCount := len(ncNow.nss)
		if ncNowCount > nodeCount {
			nodeCount = ncNowCount
			nc = ncNow
		}
	}
	return
}

// conditionNode 根据条件匹配节点单元
//
// 该方法可以用更优雅或正确的方式实现，但烧脑，性能影响不大，先这样吧
func (s *Selector) conditionNode(nc *nodeCondition, cond *condition) {
	var (
		hashKey, flexibleKey, nextFlexibleKey, distance uint64
		nextDegree                                      uint16
		ok                                              bool
	)
	if _, hashKey, ok = utils.Type2index(cond.Value); !ok {
		return
	}

	nodeLevel1 := &nodeSelector{level: 1, degreeIndex: 0, cond: cond}
	nc.nss = append(nc.nss, nodeLevel1)
	flexibleKey = hashKey
	distance = levelDistance(nodeLevel1.level)
	nextDegree = uint16(flexibleKey / distance)
	nextFlexibleKey = flexibleKey - uint64(nextDegree)*distance

	nodeLevel2 := &nodeSelector{level: 2, degreeIndex: nextDegree, cond: cond}
	nodeLevel1.nextNode = nodeLevel2
	if nil == nc.nextNode {
		nc.nextNode = &nodeCondition{nss: []*nodeSelector{}}
	}
	nc.nextNode.nss = append(nc.nextNode.nss, nodeLevel2)
	flexibleKey = nextFlexibleKey
	distance = levelDistance(nodeLevel2.level)
	nextDegree = uint16(flexibleKey / distance)
	nextFlexibleKey = flexibleKey - uint64(nextDegree)*distance

	nodeLevel3 := &nodeSelector{level: 3, degreeIndex: nextDegree, cond: cond}
	nodeLevel2.nextNode = nodeLevel3
	if nil == nc.nextNode.nextNode {
		nc.nextNode.nextNode = &nodeCondition{nss: []*nodeSelector{}}
	}
	nc.nextNode.nextNode.nss = append(nc.nextNode.nextNode.nss, nodeLevel3)
	flexibleKey = nextFlexibleKey
	distance = levelDistance(nodeLevel3.level)
	nextDegree = uint16(flexibleKey / distance)
	nextFlexibleKey = flexibleKey - uint64(nextDegree)*distance

	nodeLevel4 := &nodeSelector{level: 4, degreeIndex: nextDegree, cond: cond}
	nodeLevel3.nextNode = nodeLevel4
	if nil == nc.nextNode.nextNode.nextNode {
		nc.nextNode.nextNode.nextNode = &nodeCondition{nss: []*nodeSelector{}}
	}
	nc.nextNode.nextNode.nextNode.nss = append(nc.nextNode.nextNode.nextNode.nss, nodeLevel4)
	flexibleKey = nextFlexibleKey
	distance = levelDistance(nodeLevel4.level)
	nextDegree = uint16(flexibleKey / distance)

	nodeLevel5 := &nodeSelector{level: 5, degreeIndex: nextDegree, cond: cond}
	nodeLevel4.nextNode = nodeLevel5
	nodeLevel3.nextNode = nodeLevel4
	if nil == nc.nextNode.nextNode.nextNode.nextNode {
		nc.nextNode.nextNode.nextNode.nextNode = &nodeCondition{nss: []*nodeSelector{}}
	}
	nc.nextNode.nextNode.nextNode.nextNode.nss = append(nc.nextNode.nextNode.nextNode.nextNode.nss, nodeLevel5)
}

const (
	paramInt64 paramType = iota
	paramUint64
	paramFloat64
	paramString
	paramBool
)

type paramType int

// formatParam 梳理param的类型及值
//
// 类型：int=0;int64=1;uint64=2;string=3;bool=4
func (s *Selector) formatParam(paramValue interface{}) (paramType paramType, value interface{}, support bool) {
	switch paramValue := paramValue.(type) {
	default:
		return -1, nil, false
	case int:
		return paramInt64, int64(paramValue), true
	case int8, int16, int32, int64:
		return paramInt64, paramValue.(int64), true
	case uint8, uint16, uint32, uint, uint64, uintptr:
		return paramUint64, paramValue.(uint64), true
	case float32, float64:
		return paramFloat64, paramValue.(float64), true
	case string:
		return paramString, paramValue, true
	case bool:
		return paramBool, paramValue, true
	}
}

// pcMapName map[string]*paramCondition string
func (s *Selector) pcMapName(cond *condition) string {
	return strings.Join([]string{cond.Param, cond.Cond}, "")
}

// leftQueryIndex 索引顺序检索
//
// index 已获取索引对象
func (s *Selector) leftQueryIndex(index *Index, ns *nodeCondition, pcs map[string]*paramCondition) (int32, []interface{}) {
	var (
		count, nc int32
		nis       []interface{}
		is        = make([]interface{}, 0)
		skipIn    = s.Skip
		limitIn   uint32
	)
	for _, node := range index.node.nodes {
		if nil == ns {
			skipIn, limitIn, nc, nis = s.leftQueryNode(skipIn, limitIn, node, nil, pcs)
		} else {
			skipIn, limitIn, nc, nis = s.leftQueryNode(skipIn, limitIn, node, ns.nextNode, pcs)
		}

		count += nc
		is = append(is, nis...)
		if limitIn >= s.Limit {
			break
		}
	}
	if s.Sort == nil {
		return count, is
	}
	return count, s.shellSort(is)
}

// leftQueryNode 节点顺序检索
//
// skipIn 传入skip表示需要跳过的数量，返回skip表示剩余待跳过的数量
//
// limitIn 传入limit表示已经命中的数量，返回limit表示已经命中的数量
func (s *Selector) leftQueryNode(skipIn, limitIn uint32, node *node, ns *nodeCondition, pcs map[string]*paramCondition) (uint32, uint32, int32, []interface{}) {
	var (
		skip  = skipIn
		limit = limitIn
		count int32
		is    = make([]interface{}, 0)
	)
	if nodes := node.nodes; nil != nodes {
		for _, nd := range nodes {
			var (
				nc  int32
				nis []interface{}
			)
			if ns == nil {
				skip, limit, nc, nis = s.leftQueryNode(skip, limit, nd, nil, pcs)
			} else if s.nodeConditions(nd, ns.nextNode.nss) { // 判断当前条件是否满足，如果满足则继续下一步
				skip, limit, nc, nis = s.leftQueryNode(skip, limit, nd, ns.nextNode, pcs)
			}
			count += nc
			is = append(is, nis...)
			if limit >= s.Limit {
				break
			}
		}
	} else {
		if ns == nil {
			return s.leftQueryLeaf(skip, limit, node, nil, pcs)
		}
		return s.leftQueryLeaf(skip, limit, node, ns, pcs)
	}
	return skip, limit, count, is
}

// leftQueryLeaf 叶子节点顺序检索
//
// skip 传入skip表示需要跳过的数量，返回skip表示剩余待跳过的数量
//
// limit 传入limit表示已经命中的数量，返回limit表示已经命中的数量
func (s *Selector) leftQueryLeaf(skip, limit uint32, leaf *node, ns *nodeCondition, pcs map[string]*paramCondition) (uint32, uint32, int32, []interface{}) {
	var (
		count int32
		is    = make([]interface{}, 0)
	)
	if (nil != ns && s.leafConditions(leaf, ns.nss)) || nil == ns { // 满足等于与不等于条件
		if limit >= s.Limit {
			return skip, limit, 0, is
		}
		for position, link := range leaf.links {
			if nil == pcs || len(pcs) == 0 {
				if skip > 0 {
					skip--
					continue
				}
			}
			if s.isConditionNoIndexLeaf(ns, pcs, link.value) {
				count++
				if skip > 0 {
					skip--
					continue
				}
				limit++
				if s.delete {
					leaf.links = append(leaf.links[:position], leaf.links[position+1:]...)
				}
				is = append(is, link.value)
			}
		}
	}
	return skip, limit, count, is
}

// rightQueryIndex 索引倒序检索
//
// index 已获取索引对象
func (s *Selector) rightQueryIndex(index *Index, ns *nodeCondition, pcs map[string]*paramCondition) (int32, []interface{}) {
	var (
		count, nc int32
		nis       []interface{}
		is        = make([]interface{}, 0)
		skipIn    = s.Skip
		limitIn   uint32
	)
	lenNode := len(index.node.nodes)
	for i := lenNode - 1; i >= 0; i-- {
		if nil == ns {
			skipIn, limitIn, nc, nis = s.rightQueryNode(skipIn, limitIn, index.node.nodes[i], nil, pcs)
		} else {
			skipIn, limitIn, nc, nis = s.rightQueryNode(skipIn, limitIn, index.node.nodes[i], ns.nextNode, pcs)
		}
		count += nc
		is = append(is, nis...)
		if limitIn >= s.Limit {
			break
		}
	}
	return count, is
}

// rightQueryNode 节点倒序检索
//
// skip 传入skip表示需要跳过的数量，返回skip表示剩余待跳过的数量
//
// limit 传入limit表示已经命中的数量，返回limit表示已经命中的数量
func (s *Selector) rightQueryNode(skipIn, limitIn uint32, node *node, ns *nodeCondition, pcs map[string]*paramCondition) (uint32, uint32, int32, []interface{}) {
	var (
		skip  = skipIn
		limit = limitIn
		count int32
		is    = make([]interface{}, 0)
	)
	if nodes := node.nodes; nil != nodes {
		lenNode := len(nodes)
		for i := lenNode - 1; i >= 0; i-- {
			var (
				nc  int32
				nis []interface{}
			)
			if ns == nil {
				skip, limit, nc, nis = s.rightQueryNode(skip, limit, nodes[i], nil, pcs)
			} else if s.nodeConditions(nodes[i], ns.nextNode.nss) { // 判断当前条件是否满足，如果满足则继续下一步
				skip, limit, nc, nis = s.rightQueryNode(skip, limit, nodes[i], ns.nextNode, pcs)
			}
			count += nc
			is = append(is, nis...)
		}
	} else {
		if ns == nil {
			return s.rightQueryLeaf(skip, limit, node, nil, pcs)
		}
		return s.rightQueryLeaf(skip, limit, node, ns, pcs)
	}
	return skip, limit, count, is
}

// rightQueryLeaf 叶子节点倒序检索
//
// skip 传入skip表示需要跳过的数量，返回skip表示剩余待跳过的数量
//
// limit 传入limit表示已经命中的数量，返回limit表示已经命中的数量
func (s *Selector) rightQueryLeaf(skip, limit uint32, leaf *node, ns *nodeCondition, pcs map[string]*paramCondition) (uint32, uint32, int32, []interface{}) {
	var (
		count int32
		is    = make([]interface{}, 0)
	)
	lenLink := len(leaf.links)
	if (nil != ns && s.leafConditions(leaf, ns.nss)) || nil == ns { // 满足等于与不等于条件
		if limit >= s.Limit {
			return skip, limit, 0, is
		}
		for i := lenLink - 1; i >= 0; i-- {
			if nil == pcs || len(pcs) == 0 {
				if skip > 0 {
					skip--
					continue
				}
			}
			link := leaf.links[i]
			if s.isConditionNoIndexLeaf(ns, pcs, link.value) {
				count++
				if skip > 0 {
					skip--
					continue
				}
				limit++
				if s.delete {
					leaf.links = append(leaf.links[:i], leaf.links[i+1:]...)
				}
				is = append(is, link.value)
			}
		}
	}
	return skip, limit, count, is
}

// nodeConditions 判断当前条件集合是否满足
func (s *Selector) nodeConditions(node *node, nss []*nodeSelector) bool {
	for _, ns := range nss {
		if !s.isConditionNode(node, ns) {
			return false
		}
	}
	return true
}

// isConditionNode 判断当前条件是否满足
func (s *Selector) isConditionNode(node *node, ns *nodeSelector) bool {
	if ns != nil {
		for _, cond := range s.Conditions {
			if cond != ns.cond {
				continue
			}
			switch cond.Cond {
			case "gt":
				return s.conditionGT(node, ns)
			case "lt":
				return s.conditionLT(node, ns)
			}
		}
	}
	return true
}

// conditionGT 条件大于判断
func (s *Selector) conditionGT(node *node, ns *nodeSelector) bool {
	switch ns.level {
	default:
		return ns.degreeIndex < node.degreeIndex
	case 1, 2, 3, 4:
		return ns.degreeIndex <= node.degreeIndex
	}
}

// conditionLT 条件小于判断
func (s *Selector) conditionLT(node *node, ns *nodeSelector) bool {
	switch ns.level {
	default:
		return ns.degreeIndex > node.degreeIndex
	case 1, 2, 3, 4:
		return ns.degreeIndex >= node.degreeIndex
	}
}

// leafConditions 判断当前条件集合是否满足
func (s *Selector) leafConditions(node *node, nss []*nodeSelector) bool {
	for _, ns := range nss {
		if !s.isConditionLeaf(node, ns) {
			return false
		}
	}
	return true
}

// conditionLeaf 判断当前条件是否满足
func (s *Selector) isConditionLeaf(node *node, ns *nodeSelector) bool {
	if ns != nil {
		for _, cond := range s.Conditions {
			if cond != ns.cond {
				continue
			}
			switch cond.Cond {
			case "eq":
				return ns.level == 5 && ns.degreeIndex == node.degreeIndex
			case "dif":
				return ns.level == 5 && ns.degreeIndex != node.degreeIndex
			}
		}
	}
	return true
}

// isConditionNoIndexLeaf 判断当前条件是否满足
func (s *Selector) isConditionNoIndexLeaf(ns *nodeCondition, pcs map[string]*paramCondition, value interface{}) bool {
	for _, cond := range s.Conditions {
		if nil != ns && cond.Param == ns.nss[0].cond.Param {
			continue
		}
		pc := pcs[s.pcMapName(cond)]
		if nil == pc {
			continue
		}
		if !s.conditionValue(cond.Cond, strings.Split(cond.Param, "."), pc.paramType, pc.paramValue, value) {
			return false
		}
	}
	return true
}

// conditionValue 判断当前条件是否满足
func (s *Selector) conditionValue(cond string, params []string, paramType paramType, paramValue, objValue interface{}) bool {
	var value interface{}
	if value = s.valueFromParams(params, objValue); nil == value {
		return false
	}
	switch value := value.(type) {
	default:
		return false
	case int, int8, int16, int32, int64:
		return conditionValueInt64(cond, paramType, paramValue, value)
	case uint8, uint16, uint32, uint, uint64, uintptr:
		return conditionValueUint64(cond, paramType, paramValue, value)
	case float32, float64:
		return conditionValueFloat64(cond, paramType, paramValue, value)
	case string:
		return conditionValueString(cond, paramType, paramValue, value)
	case bool:
		if paramType != paramBool {
			return false
		}
		return conditionValueBool(cond, paramValue.(bool), value)
	}
}

// getValueFromParams 根据索引描述获取当前value
func (s *Selector) valueFromParams(params []string, value interface{}) interface{} {
	reflectObj := reflect.ValueOf(value) // 反射对象，通过reflectObj获取存储在里面的值，还可以去改变值
	valueType := reflectObj.Kind()
	if valueType == reflect.Ptr {
		if reflectObj.IsNil() {
			log.Debug("valueFromParams", log.Field("kind", valueType), log.Field("support", false))
			return nil
		}
		reflectObj = reflectObj.Elem()
		valueType = reflectObj.Kind()
	} else {
		reflectObj = reflectObj.Elem()
	}
	switch valueType {
	default:
		log.Debug("valueFromParams", log.Field("kind", reflectObj.Kind()), log.Field("support", false))
		return nil
	case reflect.Map:
		var valueResult interface{}
		interMap := value.(map[string]interface{})
		lenParams := len(params)
		for i, param := range params {
			if i == lenParams-1 {
				valueResult = interMap[param]
				break
			}
			interMap = interMap[param].(map[string]interface{})
		}
		return valueResult
	case reflect.Struct:
		var valueResult interface{}
		lenParams := len(params)
		for i, param := range params {
			if i == lenParams-1 {
				valueResult = reflectObj.FieldByName(param).Interface()
				break
			}
			reflectObj = reflectObj.FieldByName(param)
		}
		return valueResult
	}
}

// shellSort 希尔排序
func (s *Selector) shellSort(is []interface{}) []interface{} {
	log.Debug("shellSort 希尔排序", log.Field("s.Sort", s.Sort))
	if s.Sort.ASC {
		log.Debug("shellAsc 希尔顺序排序")
		return s.shellAsc(is)
	}
	log.Debug("shellDesc 希尔倒序排序")
	return s.shellDesc(is)
}

// shellAsc 希尔顺序排序
func (s *Selector) shellAsc(is []interface{}) []interface{} {
	length := len(is)
	gap := length / 2
	params := strings.Split(s.Sort.Param, ".")
	for gap > 0 {
		for i := gap; i < length; i++ {
			tempI := is[i]
			temp := s.hashKeyFromValue(params, is[i])
			preIndex := i - gap
			for preIndex >= 0 && s.hashKeyFromValue(params, is[preIndex]) > temp {
				is[preIndex+gap] = is[preIndex]
				preIndex -= gap
			}
			is[preIndex+gap] = tempI
		}
		gap /= 2
	}
	return is
}

// shellDesc 希尔倒序排序
func (s *Selector) shellDesc(is []interface{}) []interface{} {
	length := len(is)
	gap := length / 2
	params := strings.Split(s.Sort.Param, ".")
	for gap > 0 {
		for i := gap; i < length; i++ {
			tempI := is[i]
			temp := s.hashKeyFromValue(params, is[i])
			preIndex := i - gap
			for preIndex >= 0 && s.hashKeyFromValue(params, is[preIndex]) < temp {
				is[preIndex+gap] = is[preIndex]
				preIndex -= gap
			}
			is[preIndex+gap] = tempI
		}
		gap /= 2
	}
	return is
}

// hashKeyFromValue 通过Param获取该参数所属hashKey
func (s *Selector) hashKeyFromValue(params []string, value interface{}) uint64 {
	hashKey, support := s.getInterValue(params, value)
	if !support {
		return 0
	}
	return hashKey
}

// getInterValue 根据索引描述和当前检索到的value对象获取当前value对象所在索引的hashKey
func (s *Selector) getInterValue(params []string, value interface{}) (hashKey uint64, support bool) {
	reflectObj := reflect.ValueOf(value) // 反射对象，通过reflectObj获取存储在里面的值，还可以去改变值
	if reflectObj.Kind() == reflect.Map {
		interMap := value.(map[string]interface{})
		lenParams := len(params)
		var valueResult interface{}
		for i, param := range params {
			if i == lenParams-1 {
				valueResult = interMap[param]
				break
			}
			interMap = interMap[param].(map[string]interface{})
		}
		//log.Debug("getInterValue", log.Field("valueResult", valueResult))
		checkValue := reflect.ValueOf(valueResult)
		return utils.Value2hashKey(&checkValue)
	}
	log.Debug("getInterValue", log.Field("kind", reflectObj.Kind()), log.Field("support", false))
	return 0, false
}
