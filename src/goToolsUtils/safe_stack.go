package goToolsUtils

import (
	"container/list"
	"sync"
)

var myMap map[int]*list.List = make(map[int]*list.List)
var mutexInstance sync.Mutex;

var mutex sync.Mutex;

var QueueMaxId = 0
func GetQueueById(id int) *list.List {
	mutexInstance.Lock()
	defer func() {
		mutexInstance.Unlock()
	}()
	instance, ok := myMap[id]
	if !ok {
		instance = list.New()
		if id == -1 || id == QueueMaxId {
			myMap[QueueMaxId] = instance
			QueueMaxId += 1
		}else {
			panic("不能随便给id")
		}
	}
	return instance
}

func InputQueueBack(queue *list.List, element interface{}) {
	mutex.Lock()
	defer func() {
		mutex.Unlock()
	}()
	queue.PushBack(element)
}

func InputQueueFront(queue *list.List, element interface{}) {
	mutex.Lock()
	defer func() {
		mutex.Unlock()
	}()
	queue.PushFront(element)
}

func OutputQueueFront(queue *list.List) interface{} {
	mutex.Lock()
	defer func() {
		mutex.Unlock()
	}()
	output := queue.Front()
	if (output == nil) {
		return nil
	}
	queue.Remove(output)
	return output.Value
}

func OutputQueueBack(queue *list.List) interface{} {
	mutex.Lock()
	defer func() {
		mutex.Unlock()
	}()
	output := queue.Back()
	if (output == nil) {
		return nil
	}
	queue.Remove(output)
	return output.Value
}
