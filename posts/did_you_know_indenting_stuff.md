{
   "title": "Did You Know You Could Indent Stuff",
   "date": "2014-05-30T20:05:53.176518334-04:00",
   "author": "Antoine Grondin",
   "invisible": false,
   "abstract": "JSON is great you could eat that for breakfast and skip coffee.",
   "language": "en"
}

Let's list all the things that are fun about `encoding/json`.

You can marshal your things and indent them at the same time. That's crazy?

```go
// 4 spaces to split the world evenly between angry and happy people
data, err := json.MarshalIndent(myPreciousData, "", strings.Repeat(" ", 4))
```

When you `println(string(data))`, you will see a nicely-4-spaces-indented JSON
document.

Incredible, right?

## wat?

Is that it? Is that a complete blog post?

Yes. Originally it was. But since we're at JSON, did you know that: I made a package
called [`fatherhood`](https://github.com/aybabtme/fatherhood). It's not fancy but it's fast.

Recently, in our comfortable tech-echo-chamber, were presented a variety of JSON encoders that
are all faster than light in a vacuum. All this is great, since I haven't bothered writing
an encoder for [`fatherhood`](https://github.com/aybabtme/fatherhood) so having a fast encoder
should turn out to be very convenient. But I did write a decoder!

This decoder is super fast. Like really fast. It uses nothing but [megajson's](https://github.com/benbjohnson/megajson)
scanner. So it's fast like it megajson (faster than light in a vacuum).  Worst, working on fatherhood,
I got to submit PRs to megajson, rendering fatherhood almost irrelevant. The
scanner can pretty much do everything you'd need, without code generation.

## is this getting somewhere?

Nope. That's all. Now you know, [megajson's](https://github.com/benbjohnson/megajson) was
there, it's still there with gems inside that we can carve out with things like
[`fatherhood`](https://github.com/aybabtme/fatherhood).

## indenting stuff?

Yes, you can indent your JSON. You can also [indent XML](http://golang.org/pkg/encoding/xml/#MarshalIndent).
That's scary, right?

```go
v := &Person{
    Id: 13,
    FirstName: "John",
    LastName: "Doe",
    Age: 42,
}
// 16 spaces to make them regret
output, err := xml.MarshalIndent(v, "", strings.Repeat(" ", 16))
if err != nil {
    fmt.Printf("error: %v\n", err)
}
```

I haven't been in this tech-echo-chamber for a long time, only a couple years. But I was
told, very young, that XML is the devil. So I'm always happy to repeat jokes about it.
So you should indent your XML by 16 chars, at the very least. Or try this out:



```go
v := &Person{
    Id: 13,
    FirstName: "John",
    LastName: "Doe",
    Age: 42,
}
// 16 spaces to make them regret
output, err := xml.MarshalIndent(v, "", strings.Repeat("pain", 4))
if err != nil {
    fmt.Printf("error: %v\n", err)
}
// <person id="13">
// painpainpainpain<name>
// painpainpainpainpainpainpainpain<first>John</first>
// painpainpainpainpainpainpainpain<last>Doe</last>
// painpainpainpain</name>
// painpainpainpain<age>42</age>
// painpainpainpain<Married>false</Married>
// painpainpainpain<City>Hanga Roa</City>
// painpainpainpain<State>Easter Island</State>
// painpainpainpain<!-- Need more details. -->
// </person>
```

Incredible, right?

## what should I do tomorrow?

Explore the standard library again, more.

## did you really just find out about this?

Of course not! `:too-proud:`

> <small>no really I didn't.</small>

## what's going on?

I really enjoy writing Go.

Have a nice weekend folks!
