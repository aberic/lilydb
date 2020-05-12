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

package config

import (
	"errors"
	"github.com/aberic/gnomon"
	"github.com/aberic/gnomon/log"
	api "github.com/aberic/lilydb/connector/grpc"
	"gopkg.in/yaml.v3"
	"io/ioutil"
	"path/filepath"
	"sync"
)

var (
	// version 版本号
	version      = "1.0"
	confInstance *Config
	onceConfig   sync.Once
)

// YamlConfig lily启动配置文件根项目
type YamlConfig struct {
	Config *Config `yaml:"conf"`
}

// Config lily启动配置文件子项目
type Config struct {
	Port                     string `yaml:"Port"`                     // Port 开放端口，便于其它应用访问
	RootDir                  string `yaml:"RootDir"`                  // RootDir Lily服务默认存储路径
	DataDir                  string `yaml:"DataDir"`                  // DataDir Lily服务数据默认存储路径
	LimitOpenFile            int32  `yaml:"LimitOpenFile"`            // LimitOpenFile 限制打开文件描述符次数
	TLS                      bool   `yaml:"TLS"`                      // TLS 是否开启 TLS
	TLSServerKeyFile         string `yaml:"TLSServerKeyFile"`         // TLSServerKeyFile lily服务私钥
	TLSServerCertFile        string `yaml:"TLSServerCertFile"`        // TLSServerCertFile lily服务数字证书
	Limit                    bool   `yaml:"Limit"`                    // Limit 是否启用服务限流策略
	LimitMillisecond         int32  `yaml:"LimitMillisecond"`         // LimitMillisecond 请求限定的时间段（毫秒）
	LimitCount               int32  `yaml:"LimitCount"`               // LimitCount 请求限定的时间段内允许的请求次数
	LimitIntervalMicrosecond int32  `yaml:"LimitIntervalMicrosecond"` // LimitIntervalMillisecond 请求允许的最小间隔时间（微秒），0表示不限
	LogDir                   string `yaml:"LogDir"`                   // LogDir Lily服务默认日志存储路径
	LogLevel                 string `yaml:"LogLevel"`                 // LogLevel 日志级别(debug/info/warn/Error/panic/fatal)
	LogFileMaxSize           int    `yaml:"LogFileMaxSize"`           // LogFileMaxSize 每个日志文件保存的最大尺寸 单位：M
	LogFileMaxAge            int    `yaml:"LogFileMaxAge"`            // LogFileMaxAge 文件最多保存多少天
	LogUtc                   bool   `yaml:"LogUtc"`                   // LogUtc CST & UTC 时间
	Production               bool   `yaml:"Production"`               // Production 是否生产环境，在生产环境下控制台不会输出任何日志
	LilyLockFilePath         string `yaml:"lily_lock_file_path"`      // LilyLockFilePath Lily当前进程地址存储文件地址
	LilyBootstrapFilePath    string `yaml:"lily_bootstrap_file_path"` // LilyBootstrapFilePath Lily重启引导文件地址
}

// InitConfig 根据文件地址获取Config对象
func InitConfig(filePath string) *Config {
	onceConfig.Do(func() {
		confInstance = &Config{}
		if gnomon.StringIsNotEmpty(filePath) {
			if err := confInstance.yaml2Config(filePath); nil != err {
				panic(err)
			}
		}
		if _, err := confInstance.scanDefault(); nil != err {
			panic(err)
		}
		log.Fit(
			confInstance.LogLevel,
			confInstance.LogDir,
			confInstance.LogFileMaxSize,
			confInstance.LogFileMaxAge,
			confInstance.LogUtc,
			confInstance.Production,
		)
	})
	return confInstance
}

// Obtain 根据文件地址获取Config对象
func Obtain() *Config {
	onceConfig.Do(func() {
		confInstance = &Config{}
		if _, err := confInstance.scanDefault(); nil != err {
			log.Panic("ObtainConfig", log.Err(err))
		}
	})
	return confInstance
}

// Version 取得db版本号
func (c *Config) Version() string {
	return version
}

// scanDefault 扫描填充默认值
func (c *Config) scanDefault() (*Config, error) {
	if gnomon.StringIsEmpty(c.Port) {
		c.Port = "19877"
	}
	if gnomon.StringIsEmpty(c.RootDir) {
		c.RootDir = "lilyDB"
	}
	if gnomon.StringIsEmpty(c.DataDir) {
		c.DataDir = filepath.Join(c.RootDir, "data")
	}
	if gnomon.StringIsEmpty(c.LogDir) {
		c.LogDir = filepath.Join(c.RootDir, "log")
	}
	if gnomon.StringIsEmpty(c.LogLevel) {
		c.LogLevel = "debug"
	}
	if c.LogFileMaxSize < 1 {
		c.LogFileMaxSize = 64
	}
	if c.LogFileMaxAge < 1 {
		c.LogFileMaxAge = 3
	}
	if c.LimitOpenFile < 1000 {
		c.LimitOpenFile = 10000
	}
	if c.TLS {
		if gnomon.StringIsEmpty(c.TLSServerKeyFile) || gnomon.StringIsEmpty(c.TLSServerCertFile) {
			return nil, errors.New("tls server key file or cert file is nil")
		}
	}
	if c.Limit {
		if c.LimitCount < 0 || c.LimitMillisecond < 0 {
			return nil, errors.New("limit count or millisecond can not be zero")
		}
	}
	c.LilyLockFilePath = filepath.Join(c.RootDir, "lily.lock")
	c.LilyBootstrapFilePath = filepath.Join(c.DataDir, "lily.sync")
	return c, nil
}

// yaml2Config YML转配置对象
func (c *Config) yaml2Config(filePath string) error {
	data, err := ioutil.ReadFile(filePath)
	if nil != err {
		return err
	}
	ymlConfig := YamlConfig{}
	err = yaml.Unmarshal([]byte(data), &ymlConfig)
	if err != nil {
		return err
	}
	confInstance = ymlConfig.Config
	return nil
}

// Conf2RPC 转rpc对象
func (c *Config) Conf2RPC() *api.Config {
	return &api.Config{
		Port:                     c.Port,
		RootDir:                  c.RootDir,
		DataDir:                  c.DataDir,
		LogDir:                   c.LogDir,
		LimitOpenFile:            c.LimitOpenFile,
		TLS:                      c.TLS,
		TLSServerKeyFile:         c.TLSServerKeyFile,
		TLSServerCertFile:        c.TLSServerCertFile,
		Limit:                    c.Limit,
		LimitMillisecond:         c.LimitMillisecond,
		LimitCount:               c.LimitCount,
		LimitIntervalMicrosecond: c.LimitIntervalMicrosecond,
		LilyLockFilePath:         c.LilyLockFilePath,
		LilyBootstrapFilePath:    c.LilyBootstrapFilePath,
	}
}

// RPC2Config rpc转对象
func (c *Config) RPC2Config(conf *api.Config) {
	c.Port = conf.Port
	c.RootDir = conf.RootDir
	c.DataDir = conf.DataDir
	c.LogDir = conf.LogDir
	c.LimitOpenFile = conf.LimitOpenFile
	c.TLS = conf.TLS
	c.TLSServerKeyFile = conf.TLSServerKeyFile
	c.TLSServerCertFile = conf.TLSServerCertFile
	c.Limit = conf.Limit
	c.LimitMillisecond = conf.LimitMillisecond
	c.LimitCount = conf.LimitCount
	c.LimitIntervalMicrosecond = conf.LimitIntervalMicrosecond
	c.LilyLockFilePath = conf.LilyLockFilePath
	c.LilyBootstrapFilePath = conf.LilyBootstrapFilePath
}
