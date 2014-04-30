package fifobench

import (
	"container/list"
)

type ThingList struct {
	list *list.List
}

func NewList() ThingFIFO {
	return &ThingList{list.New()}
}

func (t *ThingList) Enqueue(thing Thing) {
	t.list.PushBack(thing)
}

func (t *ThingList) Peek() Thing {
	return t.list.Front().Value.(Thing)
}

func (t *ThingList) Dequeue() Thing {
	return t.list.Remove(t.list.Front()).(Thing)
}

func (t *ThingList) Len() int {
	return t.list.Len()
}

func (t *ThingList) Empty() bool {
	return t.Len() == 0
}
