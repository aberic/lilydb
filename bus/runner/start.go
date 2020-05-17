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
	"fmt"
	"github.com/aberic/gnomon/log"
	"github.com/aberic/lilydb/config"
	api "github.com/aberic/lilydb/connector/grpc"
	"google.golang.org/grpc"
	"net"
	"strings"
)

// serverStart 启动服务
func serverStart(configFilepath string) {
	conf := config.InitConfig(configFilepath)
	initLog(conf)
	rpcListener(conf)
}

func initLog(config *config.Config) {
	log.Fit(
		config.LogLevel,
		config.LogDir,
		config.LogFileMaxSize,
		config.LogFileMaxAge,
		config.LogUtc,
		config.Production,
	)
}

// rpcListener 启动rpc监听
func rpcListener(conf *config.Config) {
	var (
		listener net.Listener
		err      error
	)

	fmt.Println(strings.Join([]string{"Listen announces on the local network address with port: ", conf.Port}, ""))
	if listener, err = net.Listen("tcp", strings.Join([]string{":", conf.Port}, "")); nil != err {
		panic(err)
	}
	fmt.Println("creates a gRPC server")
	server := grpc.NewServer()
	fmt.Println("register gRPC listener")
	api.RegisterLilyAPIServer(server, &APIServer{})
	fmt.Println("OFF")
	if err = server.Serve(listener); nil != err {
		panic(err)
	}
}
