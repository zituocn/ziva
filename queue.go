package ziva

type TodoQueue interface {
	Add(tas *Task)

	AddTasks(list []*Task)

	Pop() *Task

	Clear() bool

	Size() int

	IsEmpty() bool

	Print()
}

type CreateQueue func() TodoQueue
