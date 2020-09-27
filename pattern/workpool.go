package pattern

import (
	"errors"
	"fmt"
	"runtime"
	"sync"
	"sync/atomic"
)

var (
	ErrCapacity = errors.New("routine pool's capacity is full")
)

type Task interface {
	Exec(workRoutine int)
}

type taskWrapper struct {
	task   Task
	result chan error
}

type WorkPool struct {
	shutdownQueueChannel chan string   // Channel used to shut down the queue routine.
	shutdownWorkChannel  chan struct{} // Channel used to shut down the task routines.
	noTaskChannel        chan struct{}
	shutdownWaitGroup    sync.WaitGroup   // The WaitGroup for shutting down existing routines.
	queueChannel         chan taskWrapper // Channel used to sync access to the queue.
	workChannel          chan Task        // Channel used to process task.
	queuedWork           int32            // The number of task items queued.
	activeRoutines       int32            // The number of routines active.
	queueCapacity        int32            // The max number of items we can store in the queue.
}

func New(numberOfRoutines int, queueCapacity int32) *WorkPool {
	workPool := WorkPool{
		shutdownQueueChannel: make(chan string),
		shutdownWorkChannel:  make(chan struct{}),
		noTaskChannel:        make(chan struct{}),
		queueChannel:         make(chan taskWrapper),
		workChannel:          make(chan Task, queueCapacity),
		queuedWork:           0,
		activeRoutines:       0,
		queueCapacity:        queueCapacity,
	}

	workPool.shutdownWaitGroup.Add(numberOfRoutines)

	// 执行任务的routine
	for i := 0; i < numberOfRoutines; i++ {
		go workPool.workRoutine(i)
	}

	// 提交任务的routine
	go workPool.queueRoutine()

	return &workPool
}

// Shutdown will release resources and shutdown all processing.
func (pool *WorkPool) Shutdown() (err error) {
	defer catchPanic(&err, "Shutdown")

	// 用来跳出提交任务队列的goroutine
	pool.shutdownQueueChannel <- "Down"
	<-pool.shutdownQueueChannel
	close(pool.shutdownQueueChannel)

	// 关闭提交任务队列
	close(pool.queueChannel)

	// 用来跳出执行任务队列的goroutine
	close(pool.shutdownWorkChannel)
	pool.shutdownWaitGroup.Wait()
	// 关闭真正工作队列
	close(pool.workChannel)

	return err
}

func (pool *WorkPool) PostWork(work Task) (err error) {
	defer catchPanic(&err, "PostWork")
	poolWork := taskWrapper{work, make(chan error)}
	defer close(poolWork.result)
	pool.queueChannel <- poolWork
	err = <-poolWork.result
	return err
}

func (pool *WorkPool) NoTaskChannel() chan struct{} {
	return pool.noTaskChannel
}

func (pool *WorkPool) QueuedWork() int32 {
	return atomic.AddInt32(&pool.queuedWork, 0)
}

func (pool *WorkPool) ActiveRoutines() int32 {
	return atomic.AddInt32(&pool.activeRoutines, 0)
}

func catchPanic(err *error, functionName string) {
	if r := recover(); r != nil {
		buf := make([]byte, 10000)
		runtime.Stack(buf, false)

		//logs.Debugf(functionName, "PANIC Defered [%v] : Stack Trace : %v", r, string(buf))
		if err != nil {
			*err = fmt.Errorf("%v", r)
		}
	}
}

func (pool *WorkPool) safelyDoWork(workRoutine int, poolWorker Task) {
	defer catchPanic(nil, "SafelyDoWork")
	defer func() {
		active := atomic.AddInt32(&pool.activeRoutines, -1)
		if active == 0 {
			pool.noTaskChannel <- struct{}{}
		}
	}()

	atomic.AddInt32(&pool.queuedWork, -1)
	atomic.AddInt32(&pool.activeRoutines, 1)
	poolWorker.Exec(workRoutine)
}

func (pool *WorkPool) workRoutine(workRoutine int) {
	for {
		select {
		case <-pool.shutdownWorkChannel:
			pool.shutdownWaitGroup.Done()
			return

		case poolWorker := <-pool.workChannel:
			pool.safelyDoWork(workRoutine, poolWorker)
			break
		}
	}
}

func (pool *WorkPool) queueRoutine() {
	for {
		select {
		case <-pool.shutdownQueueChannel:
			pool.shutdownQueueChannel <- "Down"
			return

		case queueItem := <-pool.queueChannel:
			// if atomic.AddInt32(&pool.queuedWork, 0) == pool.queueCapacity {
			// 	queueItem.result <- ErrCapacity
			// 	continue
			// }
			pool.workChannel <- queueItem.task
			atomic.AddInt32(&pool.queuedWork, 1)
			queueItem.result <- nil
			break
		}
	}
}
