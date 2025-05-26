package main

import (
	"container/heap"
	"testing"

	"github.com/stretchr/testify/assert"
)

type Task struct {
	Identifier int
	Priority   int
}

type Queue []*Task

func (queue Queue) Len() int {
	return len(queue)
}

func (queue Queue) Less(i, j int) bool {
	return queue[i].Priority > queue[j].Priority
}

func (queue Queue) Swap(i, j int) {
	queue[i], queue[j] = queue[j], queue[i]
}

func (queue *Queue) Push(x interface{}) {
	task := x.(*Task)
	*queue = append(*queue, task)
}

func (queue *Queue) Pop() interface{} {
	oldQueue := *queue
	queueLength := len(oldQueue)
	task := oldQueue[queueLength-1]
	*queue = oldQueue[0 : queueLength-1]

	return task
}

func (queue *Queue) UpdatePriority(taskIdentifier int, newPriority int) {
	for i, task := range *queue {
		if task.Identifier == taskIdentifier {
			task.Priority = newPriority
			heap.Fix(queue, i)
			return
		}
	}
}

type Scheduler struct {
	heap *Queue
}

func NewScheduler() Scheduler {
	queue := &Queue{}
	heap.Init(queue)

	return Scheduler{
		heap: queue,
	}
}

func (s *Scheduler) AddTask(task Task) {
	heap.Push(s.heap, &task)
}

func (s *Scheduler) ChangeTaskPriority(taskID int, newPriority int) {
	s.heap.UpdatePriority(taskID, newPriority)
}

func (s *Scheduler) GetTask() Task {
	return *heap.Pop(s.heap).(*Task)
}

func TestTrace(t *testing.T) {
	task1 := Task{Identifier: 1, Priority: 10}
	task2 := Task{Identifier: 2, Priority: 20}
	task3 := Task{Identifier: 3, Priority: 30}
	task4 := Task{Identifier: 4, Priority: 40}
	task5 := Task{Identifier: 5, Priority: 50}

	scheduler := NewScheduler()
	scheduler.AddTask(task1)
	scheduler.AddTask(task2)
	scheduler.AddTask(task3)
	scheduler.AddTask(task4)
	scheduler.AddTask(task5)

	task := scheduler.GetTask()
	assert.Equal(t, task5, task)

	task = scheduler.GetTask()
	assert.Equal(t, task4, task)

	scheduler.ChangeTaskPriority(1, 100)
	task1 = Task{Identifier: 1, Priority: 100}

	task = scheduler.GetTask()
	assert.Equal(t, task1, task)

	task = scheduler.GetTask()
	assert.Equal(t, task3, task)
}
