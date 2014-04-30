package fifobench

type ThingVector struct {
	vec []Thing
}

func NewVector() ThingFIFO {
	return &ThingVector{}
}

func (t *ThingVector) Enqueue(thing Thing) {
	t.vec = append(t.vec, thing)
}

func (t *ThingVector) Peek() Thing {
	return t.vec[0]
}

func (t *ThingVector) Dequeue() Thing {
	d := t.vec[0]
	t.vec = t.vec[1:]
	return d
}

func (t *ThingVector) Len() int {
	return len(t.vec)
}

func (t *ThingVector) Empty() bool {
	return t.Len() == 0
}
