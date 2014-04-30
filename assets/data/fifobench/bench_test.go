package fifobench

import (
	"testing"
)

const size = 10

// func Benchmark_Enqueue_32(b *testing.B)      { Enqueue(b, NewVector(), size, 32) }
// func Benchmark_Enqueue_64(b *testing.B)      { Enqueue(b, NewVector(), size, 64) }
// func Benchmark_Enqueue_128(b *testing.B)     { Enqueue(b, NewVector(), size, 128) }
// func Benchmark_Enqueue_256(b *testing.B)     { Enqueue(b, NewVector(), size, 256) }
// func Benchmark_Enqueue_512(b *testing.B)     { Enqueue(b, NewVector(), size, 512) }
// func Benchmark_Enqueue_1024(b *testing.B)    { Enqueue(b, NewVector(), size, 1024) }
// func Benchmark_Enqueue_2048(b *testing.B)    { Enqueue(b, NewVector(), size, 2048) }
// func Benchmark_Enqueue_4096(b *testing.B)    { Enqueue(b, NewVector(), size, 4096) }
// func Benchmark_Enqueue_8192(b *testing.B)    { Enqueue(b, NewVector(), size, 8192) }
// func Benchmark_Enqueue_16384(b *testing.B)   { Enqueue(b, NewVector(), size, 16384) }
// func Benchmark_Enqueue_32768(b *testing.B)   { Enqueue(b, NewVector(), size, 32768) }
// func Benchmark_Enqueue_65536(b *testing.B)   { Enqueue(b, NewVector(), size, 65536) }
// func Benchmark_Enqueue_131072(b *testing.B)  { Enqueue(b, NewVector(), size, 131072) }
// func Benchmark_Enqueue_262144(b *testing.B)  { Enqueue(b, NewVector(), size, 262144) }
// func Benchmark_Enqueue_524288(b *testing.B)  { Enqueue(b, NewVector(), size, 524288) }
// func Benchmark_Enqueue_1048576(b *testing.B) { Enqueue(b, NewVector(), size, 1048576) }
// func Benchmark_Enqueue_2097152(b *testing.B) { Enqueue(b, NewVector(), size, 2097152) }

// func Benchmark_Dequeue_32(b *testing.B)      { Dequeue(b, NewVector(), size, 32) }
// func Benchmark_Dequeue_64(b *testing.B)      { Dequeue(b, NewVector(), size, 64) }
// func Benchmark_Dequeue_128(b *testing.B)     { Dequeue(b, NewVector(), size, 128) }
// func Benchmark_Dequeue_256(b *testing.B)     { Dequeue(b, NewVector(), size, 256) }
// func Benchmark_Dequeue_512(b *testing.B)     { Dequeue(b, NewVector(), size, 512) }
// func Benchmark_Dequeue_1024(b *testing.B)    { Dequeue(b, NewVector(), size, 1024) }
// func Benchmark_Dequeue_2048(b *testing.B)    { Dequeue(b, NewVector(), size, 2048) }
// func Benchmark_Dequeue_4096(b *testing.B)    { Dequeue(b, NewVector(), size, 4096) }
// func Benchmark_Dequeue_8192(b *testing.B)    { Dequeue(b, NewVector(), size, 8192) }
// func Benchmark_Dequeue_16384(b *testing.B)   { Dequeue(b, NewVector(), size, 16384) }
// func Benchmark_Dequeue_32768(b *testing.B)   { Dequeue(b, NewVector(), size, 32768) }
// func Benchmark_Dequeue_65536(b *testing.B)   { Dequeue(b, NewVector(), size, 65536) }
// func Benchmark_Dequeue_131072(b *testing.B)  { Dequeue(b, NewVector(), size, 131072) }
// func Benchmark_Dequeue_262144(b *testing.B)  { Dequeue(b, NewVector(), size, 262144) }
// func Benchmark_Dequeue_524288(b *testing.B)  { Dequeue(b, NewVector(), size, 524288) }
// func Benchmark_Dequeue_1048576(b *testing.B) { Dequeue(b, NewVector(), size, 1048576) }
// func Benchmark_Dequeue_2097152(b *testing.B) { Dequeue(b, NewVector(), size, 2097152) }

func Benchmark_Dequeue_32(b *testing.B)      { Dequeue(b, NewList(), size, 32) }
func Benchmark_Dequeue_64(b *testing.B)      { Dequeue(b, NewList(), size, 64) }
func Benchmark_Dequeue_128(b *testing.B)     { Dequeue(b, NewList(), size, 128) }
func Benchmark_Dequeue_256(b *testing.B)     { Dequeue(b, NewList(), size, 256) }
func Benchmark_Dequeue_512(b *testing.B)     { Dequeue(b, NewList(), size, 512) }
func Benchmark_Dequeue_1024(b *testing.B)    { Dequeue(b, NewList(), size, 1024) }
func Benchmark_Dequeue_2048(b *testing.B)    { Dequeue(b, NewList(), size, 2048) }
func Benchmark_Dequeue_4096(b *testing.B)    { Dequeue(b, NewList(), size, 4096) }
func Benchmark_Dequeue_8192(b *testing.B)    { Dequeue(b, NewList(), size, 8192) }
func Benchmark_Dequeue_16384(b *testing.B)   { Dequeue(b, NewList(), size, 16384) }
func Benchmark_Dequeue_32768(b *testing.B)   { Dequeue(b, NewList(), size, 32768) }
func Benchmark_Dequeue_65536(b *testing.B)   { Dequeue(b, NewList(), size, 65536) }
func Benchmark_Dequeue_131072(b *testing.B)  { Dequeue(b, NewList(), size, 131072) }
func Benchmark_Dequeue_262144(b *testing.B)  { Dequeue(b, NewList(), size, 262144) }
func Benchmark_Dequeue_524288(b *testing.B)  { Dequeue(b, NewList(), size, 524288) }
func Benchmark_Dequeue_1048576(b *testing.B) { Dequeue(b, NewList(), size, 1048576) }
func Benchmark_Dequeue_2097152(b *testing.B) { Dequeue(b, NewList(), size, 2097152) }

func Benchmark_Enqueue_32(b *testing.B)      { Enqueue(b, NewList(), size, 32) }
func Benchmark_Enqueue_64(b *testing.B)      { Enqueue(b, NewList(), size, 64) }
func Benchmark_Enqueue_128(b *testing.B)     { Enqueue(b, NewList(), size, 128) }
func Benchmark_Enqueue_256(b *testing.B)     { Enqueue(b, NewList(), size, 256) }
func Benchmark_Enqueue_512(b *testing.B)     { Enqueue(b, NewList(), size, 512) }
func Benchmark_Enqueue_1024(b *testing.B)    { Enqueue(b, NewList(), size, 1024) }
func Benchmark_Enqueue_2048(b *testing.B)    { Enqueue(b, NewList(), size, 2048) }
func Benchmark_Enqueue_4096(b *testing.B)    { Enqueue(b, NewList(), size, 4096) }
func Benchmark_Enqueue_8192(b *testing.B)    { Enqueue(b, NewList(), size, 8192) }
func Benchmark_Enqueue_16384(b *testing.B)   { Enqueue(b, NewList(), size, 16384) }
func Benchmark_Enqueue_32768(b *testing.B)   { Enqueue(b, NewList(), size, 32768) }
func Benchmark_Enqueue_65536(b *testing.B)   { Enqueue(b, NewList(), size, 65536) }
func Benchmark_Enqueue_131072(b *testing.B)  { Enqueue(b, NewList(), size, 131072) }
func Benchmark_Enqueue_262144(b *testing.B)  { Enqueue(b, NewList(), size, 262144) }
func Benchmark_Enqueue_524288(b *testing.B)  { Enqueue(b, NewList(), size, 524288) }
func Benchmark_Enqueue_1048576(b *testing.B) { Enqueue(b, NewList(), size, 1048576) }
func Benchmark_Enqueue_2097152(b *testing.B) { Enqueue(b, NewList(), size, 2097152) }

func Enqueue(b *testing.B, fifo ThingFIFO, dataSize, fifoSize int) {
	b.ReportAllocs()

	things := NewThings(dataSize, fifoSize)

	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		for _, thing := range things {
			fifo.Enqueue(thing)
		}
	}
}

func Dequeue(b *testing.B, fifo ThingFIFO, dataSize, fifoSize int) {
	b.ReportAllocs()

	things := NewThings(dataSize, fifoSize)

	b.StopTimer()
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		for _, thing := range things {
			fifo.Enqueue(thing)
		}

		b.StartTimer()
		for _, thing := range things {
			dq := fifo.Dequeue()
			if dq != thing {
				b.FailNow()
			}
		}
		b.StopTimer()
	}
}
