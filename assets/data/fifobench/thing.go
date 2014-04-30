package fifobench

import (
	"github.com/dustin/randbo"
)

type ThingFIFO interface {
	Enqueue(thing Thing)
	Peek() Thing
	Dequeue() Thing
	Len() int
	Empty() bool
}

var rand = randbo.New()

type Thing struct {
	Data string
}

func NewThing(n int) Thing {
	return Thing{Data: genString(n)}
}

func NewThings(n, thingCount int) []Thing {
	things := make([]Thing, thingCount)
	for i := range things {
		things[i] = NewThing(n)
	}
	return things
}

func genString(n int) string {
	d := make([]byte, n)
	rand.Read(d)
	return string(d)
}
