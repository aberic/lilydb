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

package main

import (
	"context"
	"encoding/json"
	"github.com/aberic/lilydb/config"
	api "github.com/aberic/lilydb/connector/grpc"
	"github.com/aberic/lilydb/engine"
	"github.com/vmihailenco/msgpack"
	"gopkg.in/yaml.v3"
)

// APIServer APIServer
type APIServer struct {
}

// GetConf 获取数据库引擎对象
func (l *APIServer) GetConf(_ context.Context, _ *api.ReqConf) (*api.RespConf, error) {
	return &api.RespConf{Code: api.Code_Success, Conf: config.Obtain().Conf2RPC()}, nil
}

// ObtainDatabases 获取数据库集合
func (l *APIServer) ObtainDatabases(_ context.Context, _ *api.ReqDatabases) (*api.RespDatabases, error) {
	return &api.RespDatabases{Code: api.Code_Success, Databases: engine.Obtain().Databases()}, nil
}

// ObtainForms 获取数据库表集合
func (l *APIServer) ObtainForms(_ context.Context, req *api.ReqForms) (*api.RespForms, error) {
	return &api.RespForms{Code: api.Code_Success, Forms: engine.Obtain().Forms(req.DatabaseName)}, nil
}

// CreateDatabase 新建数据库
func (l *APIServer) CreateDatabase(_ context.Context, req *api.ReqCreateDatabase) (*api.Resp, error) {
	var err error
	if err = engine.Obtain().NewDatabase(req.Name, req.Comment); nil != err {
		return &api.Resp{Code: api.Code_Fail, ErrMsg: err.Error()}, err
	}
	return &api.Resp{Code: api.Code_Success}, nil
}

// CreateForm 创建表
func (l *APIServer) CreateForm(_ context.Context, req *api.ReqCreateForm) (*api.Resp, error) {
	if err := engine.Obtain().NewForm(req.DatabaseName, req.Name, req.Comment, req.FormType); nil != err {
		return &api.Resp{Code: api.Code_Fail, ErrMsg: err.Error()}, err
	}
	return &api.Resp{Code: api.Code_Success}, nil
}

// CreateKey 新建主键
func (l *APIServer) CreateKey(_ context.Context, _ *api.ReqCreateKey) (*api.Resp, error) {
	//if err := engine.Obtain().CreateKey(req.DatabaseName, req.FormName, req.KeyStructure); nil != err {
	//	return &api.Resp{Code: api.Code_Fail, ErrMsg: err.Error()}, err
	//}
	return &api.Resp{Code: api.Code_Success}, nil
}

// CreateIndex 新建索引
func (l *APIServer) CreateIndex(_ context.Context, _ *api.ReqCreateIndex) (*api.Resp, error) {
	//if err := engine.Obtain().CreateIndex(req.DatabaseName, req.FormName, req.KeyStructure); nil != err {
	//	return &api.Resp{Code: api.Code_Fail, ErrMsg: err.Error()}, err
	//}
	return &api.Resp{Code: api.Code_Success}, nil
}

// Put 新增数据
func (l *APIServer) Put(_ context.Context, req *api.ReqPut) (*api.RespPut, error) {
	var (
		v       interface{}
		hashKey uint64
		err     error
	)
	if err = json.Unmarshal(req.Value, &v); nil == err { // 尝试用json解析
		goto PUT
	}
	if err = yaml.Unmarshal(req.Value, &v); nil == err { // 尝试用yaml解析
		goto PUT
	}
	v = string(req.Value)
PUT:
	if hashKey, err = engine.Obtain().Put(req.DatabaseName, req.FormName, req.Key, v); nil != err {
		return &api.RespPut{Code: api.Code_Fail, ErrMsg: err.Error()}, err
	}
	return &api.RespPut{Code: api.Code_Success, HashKey: hashKey}, nil
}

// Set 新增数据
func (l *APIServer) Set(_ context.Context, req *api.ReqSet) (*api.RespSet, error) {
	var (
		v       interface{}
		hashKey uint64
		err     error
	)
	if err = json.Unmarshal(req.Value, &v); nil == err { // 尝试用json解析
		goto PUT
	}
	if err = yaml.Unmarshal(req.Value, &v); nil == err { // 尝试用yaml解析
		goto PUT
	}
	v = string(req.Value)
PUT:
	if hashKey, err = engine.Obtain().Set(req.DatabaseName, req.FormName, req.Key, v); nil != err {
		return &api.RespSet{Code: api.Code_Fail, ErrMsg: err.Error()}, err
	}
	return &api.RespSet{Code: api.Code_Success, HashKey: hashKey}, nil
}

// Get 获取数据
func (l *APIServer) Get(_ context.Context, req *api.ReqGet) (*api.RespGet, error) {
	var (
		v    interface{}
		data []byte
		err  error
	)
	if v, err = engine.Obtain().Get(req.DatabaseName, req.FormName, req.Key); nil != err {
		return &api.RespGet{Code: api.Code_Fail, ErrMsg: err.Error()}, err
	}
	if data, err = msgpack.Marshal(v); nil != err {
		return &api.RespGet{Code: api.Code_Fail, ErrMsg: err.Error()}, err
	}
	return &api.RespGet{Code: api.Code_Success, Value: data}, nil
}

// Remove 删除数据
func (l *APIServer) Remove(_ context.Context, req *api.ReqRemove) (*api.Resp, error) {
	if _, err := engine.Obtain().Del(req.DatabaseName, req.FormName, req.Key); nil != err {
		return &api.Resp{Code: api.Code_Fail, ErrMsg: err.Error()}, err
	}
	return &api.Resp{Code: api.Code_Success}, nil
}

// Insert 新增数据
func (l *APIServer) Insert(_ context.Context, _ *api.ReqInsert) (*api.RespInsert, error) {
	//	var (
	//		v       interface{}
	//		hashKey uint64
	//		err     error
	//	)
	//	if err = json.Unmarshal(req.Value, &v); nil == err { // 尝试用json解析
	//		goto PUT
	//	}
	//	if err = yaml.Unmarshal(req.Value, &v); nil == err { // 尝试用yaml解析
	//		goto PUT
	//	}
	//	v = string(req.Value)
	//PUT:
	//	if hashKey, err = engine.Obtain().PutD(req.Key, v); nil != err {
	//		return &api.RespInsert{Code: api.Code_Fail, ErrMsg: err.Error()}, err
	//	}
	return &api.RespInsert{Code: api.Code_Success, HashKey: 0}, nil
}

// Select 获取数据
func (l *APIServer) Select(_ context.Context, _ *api.ReqSelect) (*api.RespSelect, error) {
	//var (
	//	count int32
	//	v     interface{}
	//	s     = &Selector{}
	//	data  []byte
	//	err   error
	//)
	//if err = s.formatAPI(req.Selector); nil != err {
	//	return nil, err
	//}
	//if count, v, err = engine.Obtain().Select(req.DatabaseName, req.FormName, s); nil != err {
	//	return &api.RespSelect{Code: api.Code_Fail, ErrMsg: err.Error()}, err
	//}
	//if data, err = msgpack.Marshal(v); nil != err {
	//	return &api.RespSelect{Code: api.Code_Fail, ErrMsg: err.Error()}, err
	//}
	return &api.RespSelect{Code: api.Code_Success, Count: 0, Value: nil}, nil
}

// Delete 删除数据
func (l *APIServer) Delete(_ context.Context, _ *api.ReqDelete) (*api.RespDelete, error) {
	//var (
	//	s     = &Selector{}
	//	count int32
	//	err   error
	//)
	//if err = s.formatAPI(req.Selector); nil != err {
	//	return nil, err
	//}
	//if count, err = engine.Obtain().Delete(req.DatabaseName, req.FormName, s); nil != err {
	//	return &api.RespDelete{Code: api.Code_Fail, ErrMsg: err.Error()}, err
	//}
	return &api.RespDelete{Code: api.Code_Success, Count: 0}, nil
}
