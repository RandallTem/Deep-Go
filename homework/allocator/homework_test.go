package main

import (
	"reflect"
	"testing"
	"unsafe"

	"github.com/stretchr/testify/assert"
)

// go test -v homework_test.go

const oneByte = 1

func Defragment(memory []byte, pointers []unsafe.Pointer) {
	memoryPointer := unsafe.Pointer(&memory[0])
	fragmentedPointersMap := make(map[unsafe.Pointer]bool)
	memoryForOccupationBorder := uintptr(unsafe.Pointer(&memory[len(pointers)-1]))
	for _, pointer := range pointers {
		if uintptr(pointer) <= memoryForOccupationBorder {
			fragmentedPointersMap[pointer] = true
		}
	}
	for i := range pointers {
		for fragmentedPointersMap[memoryPointer] {
			memoryPointer = unsafe.Add(memoryPointer, oneByte)
		}
		if !fragmentedPointersMap[pointers[i]] {
			*(*byte)(memoryPointer) = *(*byte)(pointers[i])
			*(*byte)(pointers[i]) = 0b0
			pointers[i] = memoryPointer
			fragmentedPointersMap[pointers[i]] = true
			memoryPointer = unsafe.Add(memoryPointer, oneByte)
		}
	}
}

func TestDefragmentation(t *testing.T) {
	var fragmentedMemory = []byte{
		0xFF, 0x00, 0x00, 0x00,
		0x00, 0xFF, 0x00, 0x00,
		0x00, 0x00, 0xFF, 0x00,
		0x00, 0x00, 0x00, 0xFF,
	}

	var fragmentedPointers = []unsafe.Pointer{
		unsafe.Pointer(&fragmentedMemory[0]),
		unsafe.Pointer(&fragmentedMemory[5]),
		unsafe.Pointer(&fragmentedMemory[10]),
		unsafe.Pointer(&fragmentedMemory[15]),
	}

	var defragmentedPointers = []unsafe.Pointer{
		unsafe.Pointer(&fragmentedMemory[0]),
		unsafe.Pointer(&fragmentedMemory[1]),
		unsafe.Pointer(&fragmentedMemory[2]),
		unsafe.Pointer(&fragmentedMemory[3]),
	}

	var defragmentedMemory = []byte{
		0xFF, 0xFF, 0xFF, 0xFF,
		0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00,
	}

	Defragment(fragmentedMemory, fragmentedPointers)
	assert.True(t, reflect.DeepEqual(defragmentedMemory, fragmentedMemory))
	assert.True(t, reflect.DeepEqual(defragmentedPointers, fragmentedPointers))
}
