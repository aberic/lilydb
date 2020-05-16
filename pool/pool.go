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

import (
	"github.com/aberic/lilydb/engine"
	"sync/atomic"
)

// Pool 连接池
type Pool struct {
	engine      *engine.Engine // 存储引擎管理器
	jobGroup    *jobGroup      // 工作组
	active      int32          // 正在活跃的连接数量
	blockingNum int            // blockingNum 已Submit的数量
	infinite    bool           // 表示该池是否为无限容量
	options     *Options       // 参数配置
	//协程池内部的任务就绪队列
	jobsChannel chan task
}

// NewPool 新建连接池
func NewPool(engine *engine.Engine, options ...Option) (*Pool, error) {
	opts := loadOptions(options...)
	if expiry := opts.expiryDuration; expiry < 0 {
		return nil, ErrInvalidPoolExpiry
	} else if expiry == 0 {
		opts.expiryDuration = DefaultCleanIntervalTime
	}

	p := &Pool{
		engine:      engine,
		options:     opts,
		jobsChannel: make(chan task),
	}
	if p.options.minIdle <= 0 {
		p.infinite = true
	}
	p.jobGroup = newJobGroup(p)
	return p, nil
}

func (p *Pool) run() {
	p.jobGroup.run()
}

// Submit 提交任务
func (p *Pool) Submit(intent Intent, handler Handler) {
	if !p.jobGroup.work(intent, handler) {
		handler(resultFail(ErrPoolOverload))
	}
}

// Active returns the number of the currently active goroutines.
func (p *Pool) Active() int {
	return int(atomic.LoadInt32(&p.active))
}

// incActive increases the number of the currently active goroutines.
func (p *Pool) incActive() {
	atomic.AddInt32(&p.active, 1)
}

// decActive decreases the number of the currently active goroutines.
func (p *Pool) decActive() {
	atomic.AddInt32(&p.active, -1)
}
