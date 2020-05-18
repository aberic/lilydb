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

package pool

import "time"

// Option represents the optional function.
type Option func(opts *Options)

func loadOptions(options ...Option) *Options {
	opts := new(Options)
	for _, option := range options {
		option(opts)
	}
	return opts
}

// Options 包含参数配置
type Options struct {
	// 连接池中最小空闲连接数，当连接数少于此值时，连接池会创建连接来补充到该值的数量
	minIdle int
	// 连接池支持的最大连接数，这里取值为20，表示同时最多有20个数据库连接。一般把maxActive设置成可能的并发量就行了设 0 为没有限制。
	maxActive int
	// 连接池中连接用完时,新的请求等待时间，这里取值-1，表示无限等待，直到超时为止，也可取值9000，表示9秒后超时。超过时间会出错误信息
	maxWait time.Duration
	// 单个连接过期时间
	expiryDuration time.Duration
}

// WithMaxActive 连接池支持的最大连接数
func WithMaxActive(active int) Option {
	return func(opts *Options) {
		opts.maxActive = active
	}
}

// WithMinIdle 连接池中最小空闲连接数
func WithMinIdle(idle int) Option {
	return func(opts *Options) {
		opts.minIdle = idle
	}
}

// WithMaxWait 连接池中连接用完时,新的请求等待时间
func WithMaxWait(wait time.Duration) Option {
	return func(opts *Options) {
		opts.maxWait = wait
	}
}

// WithExpiryDuration 单个连接过期时间
func WithExpiryDuration(expiryDuration time.Duration) Option {
	return func(opts *Options) {
		opts.expiryDuration = expiryDuration
	}
}
