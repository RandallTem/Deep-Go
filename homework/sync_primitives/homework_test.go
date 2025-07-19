package main

import (
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

type RWMutex struct {
	mutex                   sync.RWMutex
	readingCond             *sync.Cond
	writingCond             *sync.Cond
	numberOfReaders         int32
	haveWaitingWriters      bool
	writerInCriticalSection bool
}

func NewRWMutex() *RWMutex {
	m := &RWMutex{}
	m.readingCond = sync.NewCond(&m.mutex)
	m.writingCond = sync.NewCond(&m.mutex)
	return m
}

func (m *RWMutex) Lock() {
	m.mutex.Lock()
	if m.numberOfReaders > 0 || m.writerInCriticalSection {
		m.haveWaitingWriters = true
		m.writingCond.Wait()
	}
	m.haveWaitingWriters = false
	m.writerInCriticalSection = true
	m.mutex.Unlock()
}

func (m *RWMutex) Unlock() {
	m.mutex.Lock()
	if m.haveWaitingWriters {
		m.writingCond.Signal()
	} else {
		m.readingCond.Broadcast()
	}
	m.mutex.Unlock()
}

func (m *RWMutex) RLock() {
	m.mutex.Lock()
	if m.writerInCriticalSection || m.haveWaitingWriters {
		m.readingCond.Wait()
	}
	m.numberOfReaders++
	m.mutex.Unlock()
}

func (m *RWMutex) RUnlock() {
	m.mutex.Lock()
	m.numberOfReaders--
	if m.numberOfReaders == 0 {
		m.writingCond.Signal()
	}
	m.mutex.Unlock()
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
