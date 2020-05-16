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
	"sync"
	"time"
)

type jobGroup struct {
	pool            *Pool  // 连接池对象
	stabilizeJobs   []*job // 固定工作组
	instabilityJobs []*job // 动态工作组
	once            sync.Once
	mu              sync.Mutex
}

func newJobGroup(pool *Pool) *jobGroup {
	stabilizeJobs := make([]*job, pool.options.minIdle)
	for i := 0; i < pool.options.minIdle; i++ {
		stabilizeJobs[i] = newStabilizeJob(pool)
	}
	return &jobGroup{
		pool:            pool,
		stabilizeJobs:   stabilizeJobs,
		instabilityJobs: []*job{},
	}
}

func (jg *jobGroup) run() {
	// 首先根据连接池的minIdler数量限定,开启固定数量的Worker,
	// 每一个Worker用一个Goroutine承载
	jg.once.Do(func() {
		for _, job := range jg.stabilizeJobs {
			go job.run()
		}
	})
}

func (jg *jobGroup) work(intent Intent, handler Handler) bool {
	task := task{intent: intent, handler: handler}
	if jg.pool.infinite { // 如果该池是否为无限容量
		job := newInstabilityJobs(jg.pool)
		jg.mu.Lock()
		jg.instabilityJobs = append(jg.instabilityJobs, job)
		jg.mu.Unlock()
		go job.run()
	} else if jg.pool.Active() >= jg.pool.options.maxActive { // 如果当前活跃连接数大于限制连接数
		return false
	}
	jg.pool.jobsChannel <- task
	return true
}

func (jg *jobGroup) check() {
	for index, job := range jg.instabilityJobs {
		if job.lock() {
			jg.mu.Lock()
			jg.instabilityJobs = append(jg.instabilityJobs[0:index], jg.instabilityJobs[index+1:]...)
			jg.mu.Unlock()
			break
		}
	}
	time.Sleep(jg.pool.options.expiryDuration)
	jg.check()
}
