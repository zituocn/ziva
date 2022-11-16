package ziva

import (
	"fmt"
	"sync"
)

type MemQueue struct {
	mu   *sync.Mutex
	list []*Task
}

func NewMemQueue() TodoQueue {
	return &MemQueue{
		list: make([]*Task, 0),
		mu:   &sync.Mutex{},
	}
}

func (m *MemQueue) Add(task *Task) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.list = append(m.list, task)
}

func (m *MemQueue) AddTasks(list []*Task) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.list = append(m.list, list...)
}

func (m *MemQueue) Pop() *Task {
	m.mu.Lock()
	defer m.mu.Unlock()

	if m.IsEmpty() {
		return nil
	}
	first := m.list[0]
	m.list = m.list[1:]
	return first
}

func (m *MemQueue) Clear() bool {
	if m.IsEmpty() {
		return false
	}
	for i := 0; i < m.Size(); i++ {
		m.list[i].Url = ""
	}
	m.list = nil
	return true
}

func (m *MemQueue) Size() int {
	return len(m.list)
}

func (m *MemQueue) IsEmpty() bool {
	if len(m.list) == 0 {
		return true
	}
	return false
}

func (m *MemQueue) Print() {
	fmt.Println(m.list)
}
