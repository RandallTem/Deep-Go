package main

import (
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

type RWMutex struct {
	mutex               sync.Mutex
	readCondition       *sync.Cond
	writeCondition      *sync.Cond
	numOfReaders        int
	canWrite            bool
	numOfWaitingWriters int
}

func NewRWMutex() *RWMutex {
	m := &RWMutex{}
	m.readCondition = sync.NewCond(&m.mutex)
	m.writeCondition = sync.NewCond(&m.mutex)
	return m
}

func (m *RWMutex) Lock() {
	m.mutex.Lock()
	defer func() {
		m.mutex.Unlock()
		m.numOfWaitingWriters--
		m.canWrite = true
	}()
	m.numOfWaitingWriters++
	for m.numOfReaders > 0 || m.canWrite {
		m.writeCondition.Wait()
	}
}

func (m *RWMutex) Unlock() {
	m.mutex.Lock()
	defer m.mutex.Unlock()
	m.canWrite = false
	if m.numOfWaitingWriters == 0 {
		m.writeCondition.Signal()
	} else {
		m.readCondition.Broadcast()
	}
}

func (m *RWMutex) RLock() {
	m.mutex.Lock()
	defer m.mutex.Unlock()
	for m.canWrite || m.numOfWaitingWriters > 0 {
		m.readCondition.Wait()
	}
	m.numOfReaders++
}

func (m *RWMutex) RUnlock() {
	m.mutex.Lock()
	defer m.mutex.Unlock()
	m.numOfReaders--
	if m.numOfReaders == 0 {
		m.writeCondition.Signal()
	}
}

func TestRWMutexWithWriter(t *testing.T) {
	mutex := NewRWMutex()
	mutex.Lock() // writer

	var mutualExlusionWithWriter atomic.Bool
	mutualExlusionWithWriter.Store(true)
	var mutualExlusionWithReader atomic.Bool
	mutualExlusionWithReader.Store(true)

	go func() {
		mutex.Lock() // another writer
		mutualExlusionWithWriter.Store(false)
	}()

	go func() {
		mutex.RLock() // another reader
		mutualExlusionWithReader.Store(false)
	}()

	time.Sleep(time.Second)
	assert.True(t, mutualExlusionWithWriter.Load())
	assert.True(t, mutualExlusionWithReader.Load())
}

func TestRWMutexWithReaders(t *testing.T) {
	mutex := NewRWMutex()
	mutex.RLock() // reader

	var mutualExlusionWithWriter atomic.Bool
	mutualExlusionWithWriter.Store(true)

	go func() {
		mutex.Lock() // another writer
		mutualExlusionWithWriter.Store(false)
	}()

	time.Sleep(time.Second)
	assert.True(t, mutualExlusionWithWriter.Load())
}

func TestRWMutexMultipleReaders(t *testing.T) {
	mutex := NewRWMutex()
	mutex.RLock() // reader

	var readersCount atomic.Int32
	readersCount.Add(1)

	go func() {
		mutex.RLock() // another reader
		readersCount.Add(1)
	}()

	go func() {
		mutex.RLock() // another reader
		readersCount.Add(1)
	}()

	time.Sleep(time.Second)
	assert.Equal(t, int32(3), readersCount.Load())
}

func TestRWMutexWithWriterPriority(t *testing.T) {
	mutex := NewRWMutex()
	mutex.RLock() // reader

	var mutualExlusionWithWriter atomic.Bool
	mutualExlusionWithWriter.Store(true)
	var readersCount atomic.Int32
	readersCount.Add(1)

	go func() {
		mutex.Lock() // another writer is waiting for reader
		mutualExlusionWithWriter.Store(false)
	}()

	time.Sleep(time.Second)

	go func() {
		mutex.RLock() // another reader is waiting for a higher priority writer
		readersCount.Add(1)
	}()

	go func() {
		mutex.RLock() // another reader is waiting for a higher priority writer
		readersCount.Add(1)
	}()

	time.Sleep(time.Second)

	assert.True(t, mutualExlusionWithWriter.Load())
	assert.Equal(t, int32(1), readersCount.Load())
}
