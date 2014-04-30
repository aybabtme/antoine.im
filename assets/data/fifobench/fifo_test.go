package fifobench

import (
	"github.com/dustin/go-humanize"
	"runtime"
	"testing"
)

func TestListIsEmpty(t *testing.T)           { IsEmpty(t, NewList) }
func TestListEnqueues(t *testing.T)          { Enqueues(t, NewList) }
func TestListDequeues(t *testing.T)          { Dequeues(t, NewList) }
func TestListEnqueueAndDequeue(t *testing.T) { EnqueueAndDequeue(t, NewList) }
func TestListMemoryIsBounded(t *testing.T)   { MemoryIsBounded(t, NewList) }

func TestVectorIsEmpty(t *testing.T)           { IsEmpty(t, NewVector) }
func TestVectorEnqueues(t *testing.T)          { Enqueues(t, NewVector) }
func TestVectorDequeues(t *testing.T)          { Dequeues(t, NewVector) }
func TestVectorEnqueueAndDequeue(t *testing.T) { EnqueueAndDequeue(t, NewVector) }
func TestVectorMemoryIsBounded(t *testing.T)   { MemoryIsBounded(t, NewVector) }

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

func MemoryIsBounded(t *testing.T, fifoMaker func() ThingFIFO) {
	runtime.GC()

	fifo := fifoMaker()
	mem := runtime.MemStats{}

	size := int(1e6)

	things := NewThings(10, size)
	for _, junk := range things {
		fifo.Enqueue(junk)
	}

	runtime.ReadMemStats(&mem)
	start := mem.TotalAlloc

	for i := 0; i < size; i++ {
		fifo.Enqueue(fifo.Dequeue())
	}

	runtime.ReadMemStats(&mem)
	after := mem.TotalAlloc

	t.Logf("used ~%s", humanize.Bytes(after-start))

	atmost := 2 * start
	if after > atmost {
		t.Errorf("memory unbounded, want at most %s got %s",
			humanize.Bytes(atmost),
			humanize.Bytes(after))
	}
}
