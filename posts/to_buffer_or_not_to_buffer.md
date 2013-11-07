{
    "title":"On Files: To Buffer Or Not To Buffer",
    "author":"Antoine Grondin",
    "date":"2013-11-07T00:16:30.000Z",
    "invisible": true,
    "abstract":"How should I read or write to a file? Buffering, what is buffering?"
}

## Summary

Here, I'll review quickly and roughly how files work.  Then I'll show some experimental data to explain what are buffered reads and writes.

If you know this stuff, you can just skip the text and look at the graphs, they're nice looking.

## Review of file mechanics

When accessing a file, you typically obtain a sort of socket, or _stream_, or file descriptor, from which to read or write.  Normally, a file has a beginning and an end.

In general, the process of dealing with files goes like this :

* Ask the OS to __create__ or __open__ a certain file.  Depending on the programming environment, you will get various things that represent that file.
* Use that _thing_ to __write__ and __read__ from the file in question.
* You can read and write in two fashion:
    * Sequentialy : you start from the beginning and you carry on until the end.
    * Randomly : you start from anywhere and jump to anywhere into the file.
* When you're done, tell the OS that you wish to __close__ that file.

Now, let's see that in code.  Here's a pseudo workflow (imagine that I'm handling the errors), starting by opening the file:

```go
file, err := os.Open("myfile.data")  // or os.Create(...)

// Sequential access operations
myData := make([]byte, dataLen)
n, err := file.Write(myData)

yourData := make([]byte, dataLen)
m, err := file.Read(yourData)

// Random access operations
var offset int64 = 42

myData := make([]byte, dataLen)
n, err := file.WriteAt(myData, offset)

yourData := make([]byte, dataLen)
m, err := file.ReadAt(yourData, offset)
```

We will not concern ourselves with _random accesses_ for this post, because they comes with its own performance issues.

You notice that each time we `Read` or `Write` with a file, we need to provide a piece of memory which the file will read from or write to.  Oddly, the methods are taking arrays of bytes `[]byte`.  That implies that somehow, we need to create a `[]byte` of some size.  But what size?  What if you don't know the size of what you're about to read?  Can you just pass an array of size 1 and access every part of the file one byte at a time?  Sure you can do that, but it might be a bad idea.

```go
yourData := make([]byte, 1)
n, err := file.Read(yourData)
// do something with yourData[0]
```

Alright, so if it's a bad idea, how do we read a file for which we don't know the size in advance?  What if I read the next 4 bytes but there were only 2 bytes left to read?

Let's remember how we read some part of `file` into `yourData`:

```go
n, err := file.Read(yourData)
// do something with yourData[:n]
```

In Go, `yourData[:n]` means everything in `yourData` up to `n - 1`. Notice that we don't use _all_ of `yourData`.  We only use how much `Read` said was put into `yourData`. That is, `n` bytes.

That's much more burden onto you, isn't it?  Now you need to remember not to use all of `yourData` but only its `n`th first values.  Why would you do that?

## Performance!

Alright, alright.  "_Premature optimization is the root of all evil._" Sure it is, but in our case, it's not so much premature.  Disk IO is one of the most expensive thing your computer does.  Like, orders of magnitude more expensive than any other operation on your machine... aside perhaps network accesses.  So doing your file access _right_ is kind of essential if you want to have somehow acceptable performance.

So, why should we read many bytes at once instead of one at a time?  Let's look at this graph, measured on my Macbook Pro (SSD, 256GB):

![Read/Write speed for a file of 1.0MB](/assets/data/to_buffer_or_not_to_buffer/mbpr_256GB_ssd_bench_1.0MB.svg "As the size of the data increases, the speed of access also increase")

In the above graph, we see that ass $accessSize$ increases, the speed at which we read a 1 MB file decreases.  And this decrease is exponential (see the logarithm scales).

