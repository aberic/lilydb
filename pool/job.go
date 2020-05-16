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
	"github.com/aberic/gnomon"
	"runtime"
	"sync"
	"time"
)

type job struct {
	pool          *Pool     // 连接池对象
	stabilization bool      // 是否稳定工作对象
	alive         bool      // 是否存活
	createTime    time.Time // 创建时间
	lastActive    time.Time // 最后活跃时间
	mu            sync.Locker
}

func newStabilizeJob(pool *Pool) *job {
	return &job{pool: pool, stabilization: true, alive: true}
}

func newInstabilityJobs(pool *Pool) *job {
	if runtime.NumCPU() == 1 {
		return &job{pool: pool, stabilization: false, alive: true, createTime: time.Now(), lastActive: time.Now(), mu: &sync.Mutex{}}
	}
	return &job{pool: pool, stabilization: false, alive: true, createTime: time.Now(), lastActive: time.Now(), mu: gnomon.NewSpinLock()}
}

func (j *job) run() {
	j.pool.incActive()
	// 不断的从JobsChannel内部任务队列中拿任务
	for j.alive {
		if task, ok := <-j.pool.jobsChannel; ok {
			j.pool.incIdle()
			j.mu.Lock()
			// 如果拿到任务,则执行task任务
			task.intent.run(j.pool.engine, task.handler)
			j.lastActive = time.Now()
			j.mu.Unlock()
			j.pool.decIdle()
		}
	}
}

func (j *job) lockCheckRelease() bool {
	defer j.mu.Unlock()
	j.mu.Lock()
	if j.lastActive.Add(j.pool.options.expiryDuration).Before(time.Now()) {
		j.release()
		return true
	}
	return false
}

func (j *job) release() {
	j.pool.decActive()
	j.alive = false
}
