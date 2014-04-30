package fifobench

import (
	"testing"
)

func TestVectorIsEmpty(t *testing.T)           { IsEmpty(t, NewVector) }
func TestVectorEnqueues(t *testing.T)          { Enqueues(t, NewVector) }
func TestVectorDequeues(t *testing.T)          { Dequeues(t, NewVector) }
func TestVectorEnqueueAndDequeue(t *testing.T) { EnqueueAndDequeue(t, NewVector) }

func TestListIsEmpty(t *testing.T)           { IsEmpty(t, NewList) }
func TestListEnqueues(t *testing.T)          { Enqueues(t, NewList) }
func TestListDequeues(t *testing.T)          { Dequeues(t, NewList) }
func TestListEnqueueAndDequeue(t *testing.T) { EnqueueAndDequeue(t, NewList) }

func IsEmpty(t *testing.T, fifoMaker func() ThingFIFO) {
	fifo := fifoMaker()

	if !fifo.Empty() {
		t.Fatal("should be empty")
	}

	if fifo.Len() != 0 {
		t.Fatal("should have 0 elements")
	}

}

func Enqueues(t *testing.T, fifoMaker func() ThingFIFO) {
	fifo := fifoMaker()

	for i := 1; i <= 100; i++ {
		thing := NewThing(i)
		fifo.Enqueue(thing)
		if fifo.Empty() {
			t.Fatalf("should not be empty with %d elements", i)
		}

		if fifo.Len() != i {
			t.Fatalf("should have %d elems, got %d", i, fifo.Len())
		}
	}
}

func Dequeues(t *testing.T, fifoMaker func() ThingFIFO) {
	fifo := fifoMaker()

	things := NewThings(10, 100)
	for _, thing := range things {
		fifo.Enqueue(thing)
	}

	if fifo.Empty() {
		t.Fatalf("should not be empty with %d elements", len(things))
	}

	if fifo.Len() != len(things) {
		t.Fatalf("should have %d elems, got %d", len(things), fifo.Len())
	}

	for _, want := range things {
		got := fifo.Dequeue()
		if got != want {
			t.Fatalf("want %q, got %q", want, got)
		}
	}

	if !fifo.Empty() {
		t.Fatal("should be empty")
	}

	if fifo.Len() != 0 {
		t.Fatal("should have 0 elements")
	}
}

func EnqueueAndDequeue(t *testing.T, fifoMaker func() ThingFIFO) {
	fifo := fifoMaker()

	things := NewThings(10, 100)
	for _, junk := range things {
		fifo.Enqueue(junk)
	}

	if fifo.Empty() {
		t.Fatalf("should not be empty with %d elements", len(things))
	}

	if fifo.Len() != len(things) {
		t.Fatalf("should have %d elems, got %d", len(things), fifo.Len())
	}

	for _, want := range things {
		fifo.Enqueue(NewThing(10))
		got := fifo.Dequeue()
		if got != want {
			t.Fatalf("want %q, got %q", want, got)
		}
	}
}
