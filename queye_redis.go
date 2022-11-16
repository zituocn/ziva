package ziva

import "github.com/go-redis/redis"

type RedisQueue struct {
	key string

	rdb *redis.Client
}

func NewRedisQueue(key string) TodoQueue {
	return &RedisQueue{
		key: key,
		rdb: nil,
	}
}

func (r RedisQueue) Add(tas *Task) {
	//TODO implement me
	panic("implement me")
}

func (r RedisQueue) AddTasks(list []*Task) {
	//TODO implement me
	panic("implement me")
}

func (r RedisQueue) Pop() *Task {
	//TODO implement me
	panic("implement me")
}

func (r RedisQueue) Clear() bool {
	//TODO implement me
	panic("implement me")
}

func (r RedisQueue) Size() int {
	//TODO implement me
	panic("implement me")
}

func (r RedisQueue) IsEmpty() bool {
	//TODO implement me
	panic("implement me")
}

func (r RedisQueue) Print() {
	//TODO implement me
	panic("implement me")
}
