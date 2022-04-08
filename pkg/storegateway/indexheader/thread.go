// SPDX-License-Identifier: AGPL-3.0-only
// Provenance-includes-location: https://github.com/cpg1111/threadpool-go/blob/master/thread.go
// Provenance-includes-license: MIT
// Provenance-includes-copyright: Christian Grabowski

package indexheader

import (
	"runtime"

	"go.uber.org/atomic"
)

type execResult struct {
	value     interface{}
	err       error
	populated bool
}

type OSThread struct {
	// call is used to submit closures for this thread to execute
	call chan func() (interface{}, error)

	// res is used to send the results of running closures back to callers
	res chan execResult

	// stopping is a request for this thread to stop running
	stopping chan struct{}

	// stopped is confirmation that the thread has stopped running
	stopped chan struct{}

	// localTasks keeps track of if there are any closures currently running
	localTasks *atomic.Uint64
}

func NewOSThread() *OSThread {
	return &OSThread{
		call:       make(chan func() (interface{}, error), 1),
		res:        make(chan execResult, 1),
		stopping:   make(chan struct{}),
		stopped:    make(chan struct{}),
		localTasks: atomic.NewUint64(0),
	}
}

func (o *OSThread) start() {
	runtime.LockOSThread()

	defer func() {
		runtime.UnlockOSThread()
		close(o.stopped)
	}()

	for {
		select {
		case <-o.stopping:
			// Close the channel used to read results so that any callers waiting on
			// results (which can happen if they submitted something before this method
			// has a chance to execute it but after this thread pool has been stopped)
			// so that they immediately get a zero value which they can check for.
			close(o.res)
			return
		case fn := <-o.call:
			o.localTasks.Inc()
			val, err := fn()
			o.localTasks.Dec()
			o.res <- execResult{value: val, err: err, populated: true}
			break
		}
	}
}

func (o *OSThread) Start() {
	go o.start()
}

func (o *OSThread) Stop() {
	close(o.stopping)
}

func (o *OSThread) Join() {
	for range o.stopped {
		tasks := o.localTasks.Load()
		if tasks == 0 {
			break
		}
	}
}

func (o *OSThread) Call(fn func() (interface{}, error)) (interface{}, error) {
	select {
	case <-o.stopping:
		// Make sure the pool (and hence this thread) hasn't been stopped before
		// submitting our function to run. It's still possible for the thread pool
		// to be stopped between this check and our function actually being run by
		// the .start() method. We handle this by closing the o.res channel as part
		// of stopping this thread and testing for a zero value being returned
		// (which is what reading from a closed channel does).
		return nil, ErrPoolStopped
	case o.call <- fn:
		res := <-o.res
		if !res.populated {
			return nil, ErrPoolStopped
		}
		return res.value, res.err
	}
}