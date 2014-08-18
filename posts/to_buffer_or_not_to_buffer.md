{
    "title":"On Files: To Buffer Or Not To Buffer",
    "author":"Antoine Grondin",
    "date":"2013-11-07T00:16:30.000Z",
    "invisible": false,
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

We will not concern ourselves with _random accesses_ for this post, because they come with their own performance issues.

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

So, why should we read many bytes at once instead of one at a time?  Let's look at this graph, measured on my MacBook Pro<sub>1</sub>:

<img src="/assets/data/to_buffer_or_not_to_buffer/mbpr_256GB_ssd_bench_1.0MB.png" style="width:50%;"/>

In the above graph, we see that as $access\ size$ increases, the time it takes to read a 1 MB file decreases.  And this decrease is exponential (see the two logarithm scales).

So this is on my fast SSD.  But your regular, cheap-o web instance won't have a fast SSD (or most likely not), so what will performance look like?  Well, even worst !  Let's look at the same benchmark run on an [AWS EC2 t1.micro instance](http://docs.aws.amazon.com/AWSEC2/latest/UserGuide/concepts_micro_instances.html):

<img src="/assets/data/to_buffer_or_not_to_buffer/t1_micro_bench_1.0MB.png"  style="width:50%;"/>

You can see that for accesses using small buffers, the decrease in performance is 10 times that of my laptop SSD, while with buffering, the difference is not as significant (although the instance's disk - an [EBS](https://aws.amazon.com/ebs/) - has pretty terrible performances).

So, I hope you're convinced now of the importance of doing buffered disk accesses.

## Sum It Up

How to do buffered reads and writes (ignoring all error handling):

```go
buf := bytes.NewBuffer(nil)
// Choose a decent size, the Go standard lib defines bytes.MinRead
pageSize := bytes.MinRead
data := make([]byte, pageSize)
n, _ := file.Read(data)

m, err := buf.Write(data[:n])
if err == io.EOF {
    // We're done reading, buf contains everything
}
```

Here I only handle `err` to check for the `io.EOF` that is returned when there is no more data to be read from the file.  When you implement your own file logic in Go, do three things:

* Handle all of your errors
* If you just want all the bytes, use [`ioutil.ReadAll(...)`](http://golang.org/pkg/io/ioutil/#ReadAll)
* If you want to do something in real time with the data (decode its JSON content, gunzip it on the fly, ...), don't consume the actual data byte by byte.  Instead, chain it with decoders:

```go
file, err := os.Open(...)
gzRd, err := gzip.NewReader(file)
jsonDec, err := json.NewDecoder(gzRd)

for {
    err := jsonDec.Decode(&yourStruct)
    if err == io.EOF {
        // We're done
    }
}
```

In three line, you made an on-the-fly gzip JSON decoder.  Go is pretty awesome.

## Run It Yourself

The [code to generate these graph](https://gist.github.com/aybabtme/7348714) is on my [Github](https://github.com/aybabtme/).

Please note that the code is written in a script-alike way with very bad error handling.  Also, the IO code is pretty weird because of the need to compute time measurements and artificially write with various buffer size.

This is not idiomatic, good Go code.

## Next

Right now, I don't have a post on the _random access_ side of the picture.  The day might come where I will write such a post.

<sub>1: Retina, Mid 2012 with 256GB SSD</sub>

<script type="text/javascript" src="http://cdn.mathjax.org/mathjax/latest/MathJax.js?config=TeX-AMS-MML_HTMLorMML"></script>

<script type="text/javascript">
// Single $ for inline LaTeX
MathJax.Hub.Config({
  tex2jax: {inlineMath: [['$','$'], ['\\(','\\)']]}
});
</script>
