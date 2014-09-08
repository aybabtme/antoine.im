{
   "title": "Debugging trick: debugging broken formats in binary data",
   "date": "2014-09-08T02:01:44.99914873-04:00",
   "author": "Antoine Grondin",
   "invisible": false,
   "abstract": "Helpful trick to help debugging code that manipulates binary data.",
   "language": "en"
}

I was asked to write a package that would implement a binary format. That binary format
was originally implemented in C as a shell tool, and my task was to implement the
same algorithm in Go, this time as a package that could be imported.

As always, there were edge case bugs to be found in my initial draft.
Luckily, testing simple binary formats is easy, especially when you have an existing
implementation to test against.

As I wrote the tests and fixed them, I came to feel rather confident about my
implementation. Still, herring on the side of safety, I decided to add some
fuzzy tests to stress the implementation and try to draw out some edge cases.

Quickly, I hit a bug that was rather hard to dig out.

## TL;DR

When using random data in tests, try to give a discernable pattern to the data
so that you can visually identify breaks in the pattern.

## The bug

For some reasons, my encoding would break on data generally over 1GB,
but typically not on data of less than that. To add to the weirdness, those bugs
would not always show up. Sometimes a 2GB buffer would encode and decode
without problems, then the next time, a 2GB buffer would fail, panicing
with _out of bounds_ errors.

Usually, I'd associate intermittent failures to bad concurrency, but in this
case there was no concurrency involved.

## Initial exploration

When an out of bounds error occur, my first reflex is always to obtain the value
that is sent to the code. The stacktrace reported a panic coming from a slice
allocation, during the decoding phase of the test:

```go
data := make([]byte, size)
```

So I put a logging statement before that allocation:

```go
log.Printf("size=%d", size)
data := make([]byte, size)
```

The value printed was larger than a signed 64 bits integer, so it couldn't be
allocated. This resulted in the panic. Tracing down where `size` was computed,
I noticed that it was an offset computed from the product of an index and a
block size:

```go
offset := blockIndex * blockSize
```

Now the question was to find which of those were wrong. The answer was easy
to guess: `blockSize` was a value computed once and then unchanged. Since
the decoding was working for a while and then failing, it couldn't be this
value, otherwise the decoding would have failed at the first block.

On the other side, `blockIndex` was a value extracted from a stream of data:

```go
err := binary.Read(src, byteOrder, &blockIndex)
```

The problem was thus likely that the value in the encoded buffer was wrong.
The question would now be to find what was that value, and why it was properly
read most of the time, but not always.

So I verified my assumption by printing the value of `blockIndex`:

```go
err := binary.Read(src, byteOrder, &blockIndex)
log.Printf("blockIndex=%d", blockIndex)
```

Which proved that the `blockIndex` value was taking normal values until the
very last iteration, where it would take a randomly large value.

### Looking at the data

Because the data at play is rather large, it's not practical to print it to
screen. Knowing that the data was structured in blocks of `blockSize`, I decided
to print a hexdump of the surrounding of the guilty value:

```go
block := src.Bytes()[:blockSize]
fmt.Println(hex.Dump(block))
err := binary.Read(src, byteOrder, &blockIndex)
```

The result was not very helpful, printing gribberish. This is because
the test was failing only during the fuzz tests (throwing randomly shaped
data to try and draw out corner cases).

So I decided to dump the input data, the encoded data and the decoded data
to three files, for manual examination with a hex editor.

The first thing that I noticed was that the input was indeed random; which
was expected. Then I noticed that the encoded data was looking pretty normal,
with the right type of headers and blocks, with the right values everywhere.

Finally, the decoded output was mostly like the input data, aside for the
obviously missing parts after the panic.

### Finding a needle in a haystack

What I needed to do was to find the provenance of the value I was
incorrectly reading. This value was a little endian `uint64` in a 2GB
file of random data. Luckily, changing my log statement to print the
value in its hexadecimal form:

```go
log.Printf("blockIndex=%x", blockIndex)
```

I could search for the hex value and see where it came from.

Unfortunately, being 2GB of random data, the hex value was found
in several places. It would thus be hard to find a clue about the
type of error at play here.

I say this because it's usually easy to guess the type of error we're dealing
with:

<table>
<tr>
    <th>Symptoms</th>
    <th>Probable cause</th>
    <th>Why</th>
</tr>
<tr>
    <td>Random failures correlated with use of concurrent execution.</td>
    <td>Race condition.</td>
    <td>Concurrency can produce a myriad of interleaving of executions, some of them resulting in observable failures, some of the time.</td>
</tr>
<tr>
    <td>Data off by a factor</td>
    <td>Off by one error.</td>
    <td>A consistent value that is off by one leads to data consistently off.</td>
</tr>

<tr>
    <td>Data off by a pattern.</td>
    <td>Off by one error.</td>
    <td>A value is created function of another one, and this other value is off by one.</td>
</tr>
<tr>
    <td>Data randomly off without clear pattern, intermittendly, without concurrency.</td>
    <td>#dunnolol</td>
    <td>You're off for a bad day.</td>
</tr>
</table>

But as I said, at this point I still couldn't say what kind of bug I was likely
fighting. Worst, the concurrency one would have been a good candidate, given
the intermittent errors, but was impossible since there was no concurrency in play.

## Adjusting the data to help debugging.

To be able to find a pattern in the data, while still being able of triggering
the bug, I would need to create a fuzz test that would have patterns that
could help be distinguish data from encoding headers.

Headers in this case were beginning with an incrementing integer, representing
the `blockIndex`. It would be easy to spot them from the rest of the data
if only I could create random data that is unique per block, but still
random.

The solution was rather simple:

* Pick a random integer less than `uint64 - fileSize/size(uint64)`.
* Add to that integer, the biggest value of `blockIndex` that I expected.
* Write a stream of integers counting up from this random starting point.

```go
import (
    mrand "math/rand"
    crand "crypto/rand"
)

bigI, _ := crand.Int(crand.Reader, big.NewInt(math.MaxUint64 - fileSize/size(uint64)))
ui64 := bigI.Uint64()
for i := 0; i < len(data)-width; i += width {
    binary.PutUvarint(data[i:i+width], ui64)
    ui64++
}
```


This way, I could know that all the bytes counting up starting from a large
value were part of the data, while the rest were headers.

## Hah ah!

Quickly, I identified a series of tests where the `blockIndex` was
wrongly read from. Since the wrong value of `blockIndex` was now
unique over the data, it was easy to find the region leading to
the erroneous read.

It turns out that the decoding iteration __before__
the one crashing was reading less data from the buffer than its
expected `blockSize`. This would mean that the next read
would decode a `blockIndex` from somewhere in the previous block,
leading to garbage.

Adding log output for the number of bytes read from the buffer:

```
blockIndex=344
bytesRead=4096                // normal read
blockIndex=345
bytesRead=4080                // short read
blockIndex=98177827817892442  // garbage index
bytesRead=4096                // normal read
panic: !!!
```

## In the end

In the end, I realized that I was not reading enough data from the buffer
during a specific call to `Read`. My encoded stream was wrapped in a
`bufio.Reader`, which I wrongly expected to `Read()` fully.

This wrong expectation was correct 99% of the time, except sometimes where
it would read less than all the bytes.

This would not explain why the error was intermittent: the short read
should have occured consistently from runs to run. Thinking more about it,
I'm still not sure why, but I assume it might have to do with garbage collection
running in non-deterministic ways, at random times, given the fuzzy tests.
