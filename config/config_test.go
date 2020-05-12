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

package config

import "testing"

func TestObtain(t *testing.T) {
	t.Log(Obtain())
}

func TestInitConfig(t *testing.T) {
	t.Log(InitConfig("example/conf.yaml"))
}

func TestInitConfigFileEmptyFail(t *testing.T) {
	defer func() {
		if r := recover(); r != nil {
			t.Log("recover...:", r)
		}
	}()
	t.Log(InitConfig("example/conf2.yaml"))
}

func TestInitConfigFileWrongContentFail(t *testing.T) {
	defer func() {
		if r := recover(); r != nil {
			t.Log("recover...:", r)
		}
	}()
	t.Log(InitConfig("example/conf3.yaml"))
}

func TestInitConfigTLSFail(t *testing.T) {
	defer func() {
		if r := recover(); r != nil {
			t.Log("recover...:", r)
		}
	}()
	t.Log(InitConfig("example/conf4.yaml"))
}

func TestInitConfigLimitFail(t *testing.T) {
	defer func() {
		if r := recover(); r != nil {
			t.Log("recover...:", r)
		}
	}()
	t.Log(InitConfig("example/conf5.yaml"))
}

func TestInitConfigFileNotFoundFail(t *testing.T) {
	defer func() {
		if r := recover(); r != nil {
			t.Log("recover...:", r)
		}
	}()
	t.Log(InitConfig("example/conf6.yaml"))
}

func TestConfig_Version(t *testing.T) {
	t.Log(Obtain().Version())
}

func TestConfig_Conf2RPC(t *testing.T) {
	t.Log(Obtain().Conf2RPC())
}

func TestConfig_RPC2Config(t *testing.T) {
	Obtain().RPC2Config(Obtain().Conf2RPC())
}
