package main

import (
	"errors"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

// go test -v homework_test.go

const bufferSizeMultiply = 2

var ErrPoolFull = errors.New("pool is full")
var ErrPoolClosed = errors.New("pool is closed")

type WorkerPool struct {
	workersNumber int
	buffer        chan func()
	taskGroup     sync.WaitGroup
	close         chan struct{}
	mutex         sync.RWMutex
}

func NewWorkerPool(workersNumber int) *WorkerPool {
	wp := &WorkerPool{
		workersNumber: workersNumber,
		buffer:        make(chan func(), workersNumber*bufferSizeMultiply),
		taskGroup:     sync.WaitGroup{},
		close:         make(chan struct{}),
	}
	for i := 0; i < wp.workersNumber; i++ {
		go func() {
			for task := range wp.buffer {
				task()
				wp.taskGroup.Done()
			}
		}()
	}
	return wp
}

// Return an error if the pool is full
func (wp *WorkerPool) AddTask(task func()) error {
	wp.mutex.RLock()
	defer wp.mutex.RUnlock()
	select {
	case <-wp.close:
		return ErrPoolClosed
	default:
	}
	select {
	case wp.buffer <- task:
		wp.taskGroup.Add(1)
		return nil
	default:
		return ErrPoolFull
	}
}

// Shutdown all workers and wait for all
// tasks in the pool to complete
func (wp *WorkerPool) Shutdown() {
	wp.mutex.Lock()
	select {
	case <-wp.close:
		return
	default:
		close(wp.close)
	}
	close(wp.buffer)
	wp.mutex.Unlock()
	wp.taskGroup.Wait()
}

func TestWorkerPool(t *testing.T) {
	var counter atomic.Int32
	task := func() {
		time.Sleep(time.Millisecond * 500)
		counter.Add(1)
	}

	pool := NewWorkerPool(2)
	_ = pool.AddTask(task)
	_ = pool.AddTask(task)
	_ = pool.AddTask(task)

	time.Sleep(time.Millisecond * 600)
	assert.Equal(t, int32(2), counter.Load())

	time.Sleep(time.Millisecond * 600)
	assert.Equal(t, int32(3), counter.Load())

	_ = pool.AddTask(task)
	_ = pool.AddTask(task)
	_ = pool.AddTask(task)
	pool.Shutdown() // wait tasks

	assert.Equal(t, int32(6), counter.Load())
}
