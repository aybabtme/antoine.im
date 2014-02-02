{
    "title":"Go: Non blocking producers and consumers",
    "author":"Antoine Grondin",
    "date":"2014-02-01T15:09:10.000Z",
    "invisible": false,
    "abstract":"Go offers many facilities to make concurrency easy. I demonstrate two simple idioms that are convenient.",
    "language":"en"
}
Go offers the `select` keyword as a kind of `switch` for channels. This
construct is very convenient when dealing with concurrency.  One very
common idiom is to `select` on many channels, waiting to receive on one
 of them:

```
select {
  case u := <- uChan:
    // use u
  case v := <- vChan:
    // use v
}
```

`select` will pick the path that is ready to receive data, either `u` or
`v`.  If none of the two paths are ready to receive, `select` will sleep
until such a situation occur.

Alike `switch` statements, `select` also has a default case.

```
select {
  case u := <- uChan:
    // use u
  case v := <- vChan:
    // use v
  default:
    // do something while uChan and vChan are empty
}
```

The usefulness of `default` might not be obvious to you.  On a first glance,
you might wonder what situation would lead to a good use of it.

In this post, I will illustrate two cases in which this `default` construct is
convenient.  At the end of the post, I will refresh the reader's mind on
Go channels.

## Slow producers

The `default` statement is convenient when you deal with a slow producer.
Say you want to perform an action every second, but you also want to see if
some data is ready on a channel.  This is a pretty obvious case of using
`select` with `default`:

```
// A non-buffered channel
nobodyTalking := make(chan struct{})

// Start a producer that's quite slow,
// waiting 3 seconds before sending anything
go func(sendMsg chan<- struct{}) {
  time.Sleep(time.Second * 3)
  sendMsg <- struct{}{}
}(nobodyTalking)

// 5 times, look if a message is ready, then sleep
// for a second
for i := 0; i < 5; i++ {
  select {
  case <-nobodyTalking: // only if somebody is ready to send
    log.Printf("Got a message")
  default:
    log.Printf("Nobody's talking!")
  }
  <-time.Tick(time.Second * 1)

}
```
The [output of this program](http://play.golang.org/p/KemjPa-fDz) will be:

```
2009/11/10 23:00:00 Nobody's talking!
2009/11/10 23:00:01 Nobody's talking!
2009/11/10 23:00:02 Nobody's talking!
2009/11/10 23:00:03 Got a message
2009/11/10 23:00:04 Nobody's talking!
```

### Real use-case

In a project of mine, a bunch of worker goroutines perform actions in batches.
Between batches, they look on a config channel if the master goroutine has sent
them a new configuration to use.  This let the master reconfigure the worker
goroutines without shutting them down.

## Slow consumers

The last example was pretty obvious and is seen as a canonical use case
of `select`: avoiding to block on a receive.  You might not have thought
about the inverse case.  Can you avoid blocking on a send? You can.

Say you perform work in a loop and want to report on that work, but only
if sending that report is not going to block.

```
// A non-buffered channel
nobodyListening := make(chan struct{})

// Start a *consumer* that's quite slow,
// waiting 3 seconds before receiving anything
go func(sendMsg <-chan struct{}) {
  time.Sleep(time.Second * 3)
  <-sendMsg
}(nobodyListening)

// 5 times, look if a consumer is ready, then sleep
// for a second
for i := 0; i < 5; i++ {
  select {
  case nobodyListening <- struct{}{}: // only if somebody is ready to receive
    log.Printf("Sent a message")
  default:
    log.Printf("Nobody's listening!")
  }
  <-time.Tick(time.Second * 1)

}
```
The [output of this program](http://play.golang.org/p/-U91BOUdih) will be:

```
2009/11/10 23:00:00 Nobody's listening!
2009/11/10 23:00:01 Nobody's listening!
2009/11/10 23:00:02 Nobody's listening!
2009/11/10 23:00:03 Sent a message
2009/11/10 23:00:04 Nobody's listening!
```

### Real use-case

Same project of mine, the reporter goroutine collects results from the worker
goroutine.  It computes statistics and also offer the results to clients,
provided by an exposed goroutine in the API.

However, if nobody is consuming
those results, or if the consumer is too slow to grab them, I don't want
the reporter to stop collecting results - and thus eventually block the
worker goroutines when their reporting queue gets full.  The idea is thus
to offer the results to clients only if they are ready to receive them, and
drop them otherwise.

# Primer on channels

Here's a refresher on channels, if the above was confusing to you. If you
want to read exhaustive and authoritative sources, read [Effective Go](http://golang.org/doc/effective_go.html#concurrency).

There are two types of channels in Go, buffered and non buffered ones.

```
// A non-buffered channel
uChan := make(chan int)
// A buffered channel
vChan := make(chan int, 42)
```

### Non-buffered channels

Operations on a non-buffered channels are synchronized.  This means a
__send__ on `uChan` will block until somebody is ready to __receive__ the
value.

```
// Will block
uChan <- message
```

The inverse is also true, if you try to receive a value from `uChan`:

```
// Will block
message := <- uChan
```

### Buffered channels

Buffered channels have, like their name says, a buffer.  While the buffer
is not full, messages sent on the channel will not block.

```
// Will not block
vChan <- message
```

Inversely, while the buffer is not empty, receiving a message from the
channel will not block.

```
// Will not block
message := <- vChan
```

However, when `vChan` is empty, receiving from it will block, just like on
a non-buffered channel.  The same is true for sending on `vChan` while its buffer
is full.
