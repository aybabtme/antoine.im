{
   "title": "Someone's Wrong! Membership testing in practical cases",
   "date": "2014-03-22T09:12:08.407670904-04:00",
   "author": "Antoine Grondin",
   "invisible": false,
   "abstract": "Someone was wrong on the internet. They said arrays were preferable over maps for membership testing. They were wrong.",
   "language": "en"
}
Obviously, when someone's wrong on the internet, you've gotta do something about it.  In this case, an individual on StackOverflow mentionned that using maps for membership testing was slower than simply iterating over an array, for $n$ small enough.

Their argument was that the cost of hashing the value was higher than doing a lookup in an array, for most practical sizes.  In their words, practical size meant _until over a million values_.

Now you might be wondering why I'm not directly linking to the comment in question.  The reason is that I can't find it anymore.  I only remember the comment troubled me enough to make me shout _"Someone's wrong on the internet!"_

![Someone's wrong on the internet!](http://imgs.xkcd.com/comics/duty_calls.png)

But I had doubts; their argument could have made sense.  I mean, _maybe_ the cost of hashing is greater than iterating and comparing for $n$ smaller than $something$.  I had a gut feeling it was crap, but I'm not pretentious enough to affirm it was crap before actually verifying it was crap.

All in all, finding the comment (and telling the wrongdoer they're wrong!)  doesn't matter.  What matters is __The Truth__.  In this post, we explore the truth using the [Go Programming Language](http://golang.org/), the greatest language of all (no hyperbole here).

## TL;DR

Obviously this individual was wrong.  Verified for $n>1$, membership testing on a map is always faster than on an slice. This means, in all cases you should _not_ use a slice instead of a map.

```<insert fancy graph here>```

__update__: the tests that follow assume sets of `string` types. The same tests with `int` types reveal that slices are slightly faster until $n \approx 30$.

## Problem

Membership testing consists of asking a datastructure whether it contains a value or not.  Of the many ways to implement this, two are discussed here:

### Map

Use a `map[value]bool`, then check if a value is in the map:

```go
func isMapMember(m map[string]bool, key string) bool {
  _, ok := m[key]
  return ok
}
```

### Slice

Use a `[]value`, then iterate over all the values to check if one of them is in
the slice:

```go
func isSliceMember(s []string, key string) bool {
  for _, entry := range s {
    if key == entry {
      return true
    }
  }
  return false
}
```

## Question

Which one of them is fastest?

## Hypothesis

My hypothesis is that using a slice will be faster.  I like trying to prove I'm wrong.

## Prediction

If my hypothesis is indeed right, there will be a $n$ for which using a slice will be faster than using a map.  This will mean the individual was right.  For many use cases, membership testing will be done on sets that contain a few values.

## Testing

I can think of ~~3~~ 4 dimensions that might affect the results.

1. $n$, the size of the set being tested.  The claim here is that for $n$ small enough, a slice will be faster.
2. $valSize$, the size of the individual values stored in the set. As $valSize$ increase, it is possible that the structures will perform differently than with smaller $valSize$.
3. Whether or not the entry is in the set. It could be that map are faster at determining non-membership.  Or slices.  Who knows!
4. __update__: the type of the values held in the set.

### Methodology

Taking into consideration the above dimensions, we will benchmark the two methods of testing for membership.

* For a set of $n$ unique entries, take a random entry that is in the set. Measure how much time it takes to assert that this entry is a member of the set.
* For a set of $n$ unique entries, take a random entry that is _not_ in the set. Measure how much time it takes to assert that this entry is _not_ a member of the set.

The first step is to create a set of $n$ unique entries, each randomly generated strings of length $valSize$. The actual benchmarking code separates those in two funcs, so to run one type of benchmark at a time.

```go
func makeRandomSet(n, valSize int) (map[string]bool, []string) {
  mapset := make(map[string]bool, n)
  sliceset := make([]string, 0, n)

  for len(mapset) != n {
    str := makeRandomString(valSize)
    if _, ok := mapset[str]; !ok {
      mapset[str] = true
      sliceset = append(sliceset, str)
    }
  }
  return mapset, sliceset
}
```

Then, we have two paths: asserting membership, asserting non-membership.

#### Membership

Take a sample from the set:

```go
func sampleInSlice(sliceset []string) {
  randIdx := rand.Intn(len(sliceset))
  return sliceset[randIdx]
}
```

Sampling in the `map` is done the same way, but you first need to extract the map keys into a slice.  For benchmarking purposes, this function actually returns a slice of unique samples, and then measure how much time it takes to assert the membership of each.

We will use this helper to benchmark maps:

```go
func benchMapMembers(b *testing.B, size int, keySize int) {
  m := makeRandomMap(size, keySize)
  samples := sampleFromMap(m, b.N)

  var member string
  var ok bool
  b.ResetTimer()
  for i := 0; i < b.N; i++ {
    member = samples[i]
    ok = isMapMember(m, member)
  }
  _ = ok
}
```

And this one to benchmark slices:

```go
func benchSliceMembers(b *testing.B, size int, keySize int) {
  s := makeRandomSlice(size, keySize)
  samples := sampleFromSlice(s, b.N)

  var member string
  var ok bool
  b.ResetTimer()
  for i := 0; i < b.N; i++ {
    member = samples[i]
    ok = isSliceMember(s, member)
  }
  _ = ok
}
```

#### Non-membership

Find entries that aren't in the set:

```go
func sampleNotInMap(mapset map[string]bool, valSize int) (sample string) {
  for {
    str := makeRandomString(valSize)
    if !mapset[str] {
      return str
    }
  }
}
```

Sampling in the `slice` is done the same way.  Then we also create benchmark helpers that are quite similar to the _membership_ helpers.

For maps:

```go
func benchMapNotMembers(b *testing.B, size int, keySize int) {
  m := makeRandomMap(size, keySize)
  samples := sampleNotInMap(m, b.N, keySize)

  var member string
  var ok bool
  b.ResetTimer()
  for i := 0; i < b.N; i++ {
    member = samples[i]
    ok = isMapMember(m, member)
  }
  _ = ok
}
```

For slices:

```go
func benchSliceMembers(b *testing.B, size int, keySize int) {
  s := makeRandomSlice(size, keySize)
  samples := sampleFromSlice(s, b.N)

  var member string
  var ok bool
  b.ResetTimer()
  for i := 0; i < b.N; i++ {
    member = samples[i]
    ok = isSliceMember(s, member)
  }
  _ = ok
}
```

#### Faceoff!

We run the membership helpers using a series of benchmarks, with $n$ in $(2,3,4,5,6,7,8,9,10, 100, 1000, 10000, 100000, 1000000)$ and $valSize$ in $(10, 100)$.

```go
func BenchmarkMap_2key_10bytes(b *testing.B)        { benchMapMembers(b, 2, 10) }
func BenchmarkMap_3key_10bytes(b *testing.B)        { benchMapMembers(b, 3, 10) }
...
func BenchmarkMap_1000000key_100bytes(b *testing.B) { benchMapMembers(b, 1000000, 100) }

func BenchmarkSlice_2key_10bytes(b *testing.B)        { benchSliceMembers(b, 2, 10) }
func BenchmarkSlice_3key_10bytes(b *testing.B)        { benchSliceMembers(b, 3, 10) }
...
func BenchmarkSlice_1000000key_100bytes(b *testing.B) { benchSliceMembers(b, 1000000, 100) }
```

Testing for non-membership is much more expensive than testing for membership, since we need to create random strings that aren't in the original set.  This is very time consuming to setup, so I've limited this part of the test to a few values, and assume that the behavior of non-membership will be consistent with the membership ones.  This is an acceptable assumption given that the experiment I did measure suggest that they are indeed.

Also, Go will kill benchmarks running for more than 600s (10min). The above combinations run is just short of 600s (584s).

```go
func BenchmarkMapNot_10key_10bytes(b *testing.B)      { benchMapNotMembers(b, 10, 10) }
func BenchmarkMapNot_1000key_10bytes(b *testing.B)    { benchMapNotMembers(b, 1000, 10) }
func BenchmarkMapNot_1000000key_10bytes(b *testing.B) { benchMapNotMembers(b, 1000000, 10) }

func BenchmarkSliceNot_10key_10bytes(b *testing.B)      { benchSliceNotMembers(b, 10, 10) }
func BenchmarkSliceNot_10000key_10bytes(b *testing.B)   { benchSliceNotMembers(b, 10000, 10) }
func BenchmarkSliceNot_1000000key_10bytes(b *testing.B) { benchSliceNotMembers(b, 1000000, 10) }
```

## Results

The code to run the benchmarks yourself is on [github](https://gist.github.com/aybabtme/9653488c4f910097b109).  The results of running this benchmark on my Mac follow.

### Membership

<table>
<tr>
  <th>$n$</th>
  <th>$valSize$</th>
  <th>Measurements (slice)</th>
  <th>Measurements (map)</th>
  <th>ns/op (slice)</th>
  <th>ns/op (map)</th>
</tr>
<tr>  <td>2</td>        <td>10</td>   <td>100000000</td>  <td>100000000</td>  <td>17.0</td>     <td>14.4</td>   </tr>
<tr>  <td>3</td>        <td>10</td>   <td>100000000</td>  <td>100000000</td>  <td>23.6</td>     <td>18.8</td>   </tr>
<tr>  <td>4</td>        <td>10</td>   <td>100000000</td>  <td>100000000</td>  <td>28.6</td>     <td>22.5</td>   </tr>
<tr>  <td>5</td>        <td>10</td>   <td>50000000</td>   <td>100000000</td>  <td>33.5</td>     <td>26.3</td>   </tr>
<tr>  <td>6</td>        <td>10</td>   <td>50000000</td>   <td>100000000</td>  <td>40.8</td>     <td>30.1</td>   </tr>
<tr>  <td>7</td>        <td>10</td>   <td>50000000</td>   <td>50000000</td>   <td>45.4</td>     <td>33.5</td>   </tr>
<tr>  <td>8</td>        <td>10</td>   <td>50000000</td>   <td>50000000</td>   <td>52.4</td>     <td>36.0</td>   </tr>
<tr>  <td>9</td>        <td>10</td>   <td>50000000</td>   <td>50000000</td>   <td>54.2</td>     <td>36.6</td>   </tr>
<tr>  <td>10</td>       <td>10</td>   <td>50000000</td>   <td>50000000</td>   <td>60.4</td>     <td>38.1</td>   </tr>
<tr>  <td>100</td>      <td>10</td>   <td>5000000</td>    <td>50000000</td>   <td>492</td>      <td>39.3</td>   </tr>
<tr>  <td>1000</td>     <td>10</td>   <td>500000</td>     <td>50000000</td>   <td>4829</td>     <td>39.3</td>   </tr>
<tr>  <td>10000</td>    <td>10</td>   <td>50000</td>      <td>50000000</td>   <td>48172</td>    <td>48.4</td>   </tr>
<tr>  <td>100000</td>   <td>10</td>   <td>5000</td>       <td>50000000</td>   <td>480773</td>   <td>77.0</td>   </tr>
<tr>  <td>1000000</td>  <td>10</td>   <td>500</td>        <td>20000000</td>   <td>4853658</td>  <td>149</td>    </tr>
<tr>  <td>2</td>        <td>100</td>  <td>100000000</td>  <td>100000000</td>  <td>16.9</td>     <td>12.4</td>   </tr>
<tr>  <td>3</td>        <td>100</td>  <td>100000000</td>  <td>100000000</td>  <td>23.9</td>     <td>15.9</td>   </tr>
<tr>  <td>4</td>        <td>100</td>  <td>100000000</td>  <td>100000000</td>  <td>29.4</td>     <td>18.0</td>   </tr>
<tr>  <td>5</td>        <td>100</td>  <td>50000000</td>   <td>100000000</td>  <td>36.2</td>     <td>20.8</td>   </tr>
<tr>  <td>6</td>        <td>100</td>  <td>50000000</td>   <td>100000000</td>  <td>42.0</td>     <td>22.6</td>   </tr>
<tr>  <td>7</td>        <td>100</td>  <td>50000000</td>   <td>100000000</td>  <td>45.4</td>     <td>24.6</td>   </tr>
<tr>  <td>8</td>        <td>100</td>  <td>50000000</td>   <td>100000000</td>  <td>50.9</td>     <td>26.8</td>   </tr>
<tr>  <td>9</td>        <td>100</td>  <td>50000000</td>   <td>50000000</td>   <td>56.5</td>     <td>67.7</td>   </tr>
<tr>  <td>10</td>       <td>100</td>  <td>50000000</td>   <td>50000000</td>   <td>61.3</td>     <td>68.0</td>   </tr>
<tr>  <td>100</td>      <td>100</td>  <td>5000000</td>    <td>50000000</td>   <td>518</td>      <td>68.8</td>   </tr>
<tr>  <td>1000</td>     <td>100</td>  <td>500000</td>     <td>50000000</td>   <td>4940</td>     <td>70.2</td>   </tr>
<tr>  <td>10000</td>    <td>100</td>  <td>50000</td>      <td>20000000</td>   <td>49396</td>    <td>81.0</td>   </tr>
<tr>  <td>100000</td>   <td>100</td>  <td>5000</td>       <td>10000000</td>   <td>608481</td>   <td>172</td>    </tr>
<tr>  <td>1000000</td>  <td>100</td>  <td>200</td>        <td>10000000</td>   <td>6361138</td>  <td>199</td>    </tr>
</table>

### Non-membership

<table>
<tr>
  <th>$n$</th>
  <th>$valSize$</th>
  <th>Measurements (slice)</th>
  <th>Measurements (map)</th>
  <th>ns/op (slice)</th>
  <th>ns/op (map)</th>
</tr>
<tr>  <td>10</td>       <td>10</td> <td>20000000</td> <td>50000000</td>  <td>97.7</td>    <td>36.5</td> </tr>
<tr>  <td>10000</td>    <td>10</td> <td>20000</td>    <td>50000000</td>  <td>85053</td>   <td>36.9</td> </tr>
<tr>  <td>1000000</td>  <td>10</td> <td>200</td>      <td>20000000</td>  <td>9101328</td> <td>123</td>  </tr>
</table>

The results are clear, my hypothesis was wrong. For $n > 1$, membership testing is always faster using a map.

#### So clearly, someone was wrong on the Internet!!!!

## Update

The above was written considering `string` as the type held by the set.  It turns out that for sets of `int`, slices are slightly faster than maps until $n \approx 30$.

### Integer membership

<table>
<tr>
  <th>$n$</th>
  <th>Measurements (slice)</th>
  <th>Measurements (map)</th>
  <th>ns/op (slice)</th>
  <th>ns/op (map)</th>
</tr>
<tr><td>2</td><td>200000000</td><td>100000000</td><td>9.02</td><td>11.3</td></tr>
<tr><td>3</td><td>100000000</td><td>100000000</td><td>11.1</td><td>14.0</td></tr>
<tr><td>4</td><td>100000000</td><td>100000000</td><td>12.3</td><td>16.0</td></tr>
<tr><td>5</td><td>100000000</td><td>100000000</td><td>13.2</td><td>16.4</td></tr>
<tr><td>6</td><td>100000000</td><td>100000000</td><td>13.7</td><td>17.4</td></tr>
<tr><td>7</td><td>100000000</td><td>100000000</td><td>14.5</td><td>19.4</td></tr>
<tr><td>8</td><td>100000000</td><td>100000000</td><td>15.1</td><td>20.5</td></tr>
<tr><td>9</td><td>100000000</td><td>50000000</td><td>16.0</td><td>29.9</td></tr>
<tr><td>10</td><td>100000000</td><td>50000000</td><td>16.7</td><td>29.9</td></tr>
<tr><td>20</td><td>100000000</td><td>50000000</td><td>24.6</td><td>29.8</td></tr>
<tr><td>30</td><td>50000000</td><td>50000000</td><td>31.1</td><td>28.5</td></tr>
<tr><td>40</td><td>50000000</td><td>50000000</td><td>35.3</td><td>31.6</td></tr>
<tr><td>50</td><td>50000000</td><td>50000000</td><td>39.5</td><td>30.7</td></tr>
<tr><td>100</td><td>50000000</td><td>50000000</td><td>56.2</td><td>30.6</td></tr>
<tr><td>1000</td><td>5000000</td><td>50000000</td><td>340</td><td>29.8</td></tr>
<tr><td>10000</td><td>500000</td><td>50000000</td><td>3212</td><td>32.6</td></tr>
<tr><td>100000</td><td>50000</td><td>50000000</td><td>31051</td><td>40.4</td></tr>
<tr><td>1000000</td><td>5000</td><td>50000000</td><td>331630</td><td>74.7</td></tr>
</table>

### Integer non-membership

<table>
<tr>
  <th>$n$</th>
  <th>Measurements (slice)</th>
  <th>Measurements (map)</th>
  <th>ns/op (slice)</th>
  <th>ns/op (map)</th>
</tr>
<tr><td>10</td><td>100000000</td><td>100000000</td><td>18.0</td><td>25.4</td></tr>
<tr><td>10000</td><td>500000</td><td>100000000</td><td>6220</td><td>25.5</td></tr>
<tr><td>1000000</td><td>5000</td><td>20000000</td><td>718006</td><td>80.7</td></tr>
</table>

#### Someone was not so clearly wrong on the Internet!!!!

<script type="text/javascript" src="http://cdn.mathjax.org/mathjax/latest/MathJax.js?config=TeX-AMS-MML_HTMLorMML"></script>

<script type="text/javascript">
// Single $ for inline LaTeX
MathJax.Hub.Config({
  tex2jax: {inlineMath: [['$','$'], ['\\(','\\)']]}
});
</script>
