{
   "title": "Someone's Right! Building a queue: list or slices?",
   "date": "2014-04-29T22:03:22.386395056-04:00",
   "author": "Antoine Grondin",
   "invisible": false,
   "abstract": "Should you use a slice or a list to build a queue? Your CS textbook says that linked lists are better, is this true?",
   "language": "en"
}

In a previous post, I was troubled about a claim that was contrary to my
instincts; whether [membership testing](/posts/someone_is_wrong_on_the_internet)
was faster using slices or maps, for small $N$.  I thought the individual
proclaiming such things was wrong, but was self-critical enough to consider
being wrong. After testing the claims, my intuitions were confirmed.

## A new challenge to my intuitions

A situation of this kind occured again, this time after reading
[Tv](https://twitter.com/tv) hating on [`container/list`](http://golang.org/pkg/container/list/)
for a second or third time. He's a smart guy so I tend to believe him when he says :

> Let me put it this way: it's actually hard to make a more wasteful data
> structure than `container/list`.  Even for listy operations,
> `container/list` is pretty much the worst possible thing. Here's some
> keywords to read up on: "cacheline", "pipeline stall", "branch prediction".
>
> -- Tv on [#go-nuts](https://botbot.me/freenode/go-nuts/msg/14004767/)

Those look like plausible, convincing words... but I'm hard to convince: this
goes against my instinct and what I've read before. So it's worth verifying.
At least to get a feel for it: see if it's true or not, to see just how similar or
different they are.

Even more plausible and drastic _numbers_ are from [Caleb Space](https://github.com/cespare) who replaced
a `list.List` with a slice in an commit to package [`github.com/bmizerany/perks`](https://github.com/bmizerany/perks),
(computes approximate quantiles on unbounded data streams):

>> ```
>> quantile: Replace container/list with a slice
>> Better data locality trumps asymptotic behavior in this case.
>>
>> benchmark                               old ns/op     new ns/op     delta
>> BenchmarkQuerySmallEpsilon              44491         6782          -84.76%
>> BenchmarkInsertBiasedSmallEpsilon       2641          871           -67.02%
>> BenchmarkQuery                          691           306           -55.72%
>> BenchmarkInsertBiased                   324           177           -45.37%
>> BenchmarkInsertTargetedSmallEpsilon     1016          616           -39.37%
>> BenchmarkInsertTargeted                 294           191           -35.03%
>> ```
>>
>> -- [Cespare](https://github.com/cespare/perks/commit/456f18a8e50eba8f1ea6d8728e8000072e3b322c)

In face of such counter-examples, one can do two things:

* Accept them.
* Accept them after writing benchmarks, because [NIH](https://en.wikipedia.org/wiki/Not_invented_here).

## TL;DR

Slices are much faster than linked list for use as a FIFO.

**caveat**: for the benchmarks below, otherwise YMMV.

## On micro benchmarks

Microbenchmarks like this are pretty useless.  You should definitely not
take the numbers here and walk away thinking this is the New Truth. Conditions
will vary, your usage will be different, etc.  The only use you can make of
these numbers is to have a feel for the difference between the two things tested.

When you ask yourself which of the two is preferable, you will know that:

* X has a nicer API than Y.
* I'm worried that Y might be more performant.
* Now I know that if it's true, the different (will/won't) matter.

## Implementing FIFOs, slices vs `container/list`

The case I wanted to use a queue for was simple:

* Queue things to be processed.
* Dequeue them when they're being processed.
* Re-enqueue things when the processing fails.
* Peek at the next thing to process.
* Check if there's anything to process.
* Have in order visibility over the things in the queue.

I didn't need priorities, deduplication, deletion, etc.  Really, just a queue.

So for testing, I came up with this interface, which fulfils the use case I
would have.

```go
type ThingFIFO interface {
  Enqueue(thing Thing)
  Peek() Thing
  Dequeue() Thing
  Len() int
  Empty() bool
}
```

## Using a slice

A slice-based FIFO (made with a slice) is pretty easy to implement. First
of, the data will look like this, where `Thing` is anything:

```go
type ThingVector struct {
  vec []Thing
}
```

The idea is to `append` to a slice when you enqueue:

```go
func (t *ThingVector) Enqueue(thing Thing) {
  t.vec = append(t.vec, thing)
}
```

... and reslice the slice when you dequeue.

```go
func (t *ThingVector) Dequeue() Thing {
  // Could be done in a single line, but I find that's clearer
  d := t.vec[0]
  t.vec = t.vec[1:]
  return d
}
```
The other methods are trivial, but you can have a peek at the whole source on
my [Github](https://github.com/aybabtme/antoine.im/tree/master/assets/data/fifobench/slice_fifo.go)).

## Using a list (`container/list`)

A queue built using a linked list will also be easy, if perhaps ugly, to
implement.  Go lacking generics, using `container/list` means doing type
assertions over the elements of the list, which is not pretty in the
opinion of some, where "some" includes myself.

Using a `container/list`, it will look like this:

```go
type ThingList struct {
  list *list.List
}

func NewList() ThingFIFO {
  return &ThingList{list.New()}
}
```

To enqueue, put the thing at the end of the queue:

```go
func (t *ThingList) Enqueue(thing Thing) {
  t.list.PushBack(thing)
}
```

...and remove the front element when you dequeue:

```go
func (t *ThingList) Dequeue() Thing {
  return t.list.Remove(t.list.Front()).(Thing)
}
```

Again, the other methods are trivial, but can be found on [Github](https://github.com/aybabtme/antoine.im/tree/master/assets/data/fifobench/list_fifo.go).

## Faceoff

For each type of implementation, we will test how `Enqueue` and `Dequeue` work.
I don't really care about the other operations since they are obviously constant
in time.

We will benchmark the different ways to do **enqueuing** using the following code:

```go
func Enqueue(b *testing.B, fifo ThingFIFO, dataSize, fifoSize int) {
  // Reports memory allocations
  b.ReportAllocs()

  // Create a fifoSize things, each filled with strings of random
  // data of size dataSize
  things := NewThings(dataSize, fifoSize)

  // Start measuring here
  b.ResetTimer()
  for n := 0; n < b.N; n++ {
    // For each measurement, enqueue all the things we've prepared
    for _, thing := range things {
      fifo.Enqueue(thing)
    }
  }
}
```

and **dequeuing** using this code:

```go
func Dequeue(b *testing.B, fifo ThingFIFO, dataSize, fifoSize int) {
  // Reports memory allocations
  b.ReportAllocs()

  // Create a fifoSize things, each filled with strings of random
  // data of size dataSize
  things := NewThings(dataSize, fifoSize)

  // Stop the timer and reset it, because we only want to
  // measure the parts where we dequeue
  b.StopTimer()
  b.ResetTimer()

  for n := 0; n < b.N; n++ {
    // Add all the things to the FIFO
    for _, thing := range things {
      fifo.Enqueue(thing)
    }

    // Then start measuring how much time it takes to
    // dequeue everything
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

```


## Results!

For $dataSize = 10$ and $n>=32$, the results follow.  A positive `delta` means
the slice is faster than the list.

A queue implemented with slices is always faster than a linked list, by
a large margin for enqueuing:

 Enqueue of $n$ things  | slice ns/op  | list ns/op   | delta
------------------------|--------------|--------------|------------
  32                    | 3178         | 8163         | +156.86%
  64                    | 4787         | 14648        | +206.00%
  128                   | 8379         | 36881        | +340.16%
  256                   | 16716        | 78247        | +368.10%
  512                   | 33742        | 145720       | +331.87%
  1024                  | 83273        | 310665       | +273.07%
  2048                  | 151004       | 543738       | +260.08%
  4096                  | 261446       | 936551       | +258.22%
  8192                  | 528281       | 2376402      | +349.84%
  16384                 | 1059136      | 4421926      | +317.50%
  32768                 | 2096680      | 9440943      | +350.28%
  65536                 | 4253232      | 14862841     | +249.45%
  131072                | 10186608     | 39377239     | +286.56%
  262144                | 19282298     | 80871856     | +319.41%
  524288                | 34790846     | 183769573    | +428.21%
  1048576               | 72416166     | 292706216    | +304.20%
  2097152               | 166536357    | 594833520    | +257.18%

...and faster by a non-negligeable margin for dequeuing, which is the case
we would expect to actually favor lists the most:

 Dequeue of $n$ things  | slice ns/op  | list ns/op   | delta
------------------------|--------------|--------------|------------
  32                    | 950          | 1206         | +26.95%
  64                    | 1829         | 2389         | +30.62%
  128                   | 3593         | 4790         | +33.31%
  256                   | 7023         | 9490         | +35.13%
  512                   | 13956        | 18862        | +35.15%
  1024                  | 28074        | 37508        | +33.60%
  2048                  | 55562        | 75552        | +35.98%
  4096                  | 110502       | 154085       | +39.44%
  8192                  | 220145       | 305980       | +38.99%
  16384                 | 441996       | 604432       | +36.75%
  32768                 | 886141       | 1213378      | +36.93%
  65536                 | 1755980      | 2400093      | +36.68%
  131072                | 3511602      | 4837772      | +37.77%
  262144                | 7011034      | 9592371      | +36.82%
  524288                | 13938384     | 19138455     | +37.31%
  1048576               | 29010074     | 38539214     | +32.85%
  2097152               | 57266416     | 80684682     | +40.89%

You can see the full results [here](/assets/data/fifobench/benchcmp.txt).

The results are consistent when $n<32$ for enqueuing, but the dequeuing
benchmarks take too long to converge to meaninful results, so I've gave up
on producing them.

Still, 32 is a pretty decent number and you won't notice any difference,
 using a list or a slice, for 32 or less elements.

## Conclusion

That's it, Tv was right.  The more you know!  If you want to have a look
at the full code, find it [here](https://github.com/aybabtme/antoine.im/tree/master/assets/data/fifobench).

<script type="text/javascript" src="http://cdn.mathjax.org/mathjax/latest/MathJax.js?config=TeX-AMS-MML_HTMLorMML"></script>

<script type="text/javascript">
// Single $ for inline LaTeX
MathJax.Hub.Config({
  tex2jax: {inlineMath: [['$','$'], ['\\(','\\)']]}
});
</script>
