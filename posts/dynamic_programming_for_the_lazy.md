{
    "title":"Dynamic Programming for the Impatient, Dumb and Lazy Ones",
    "author":"Antoine Grondin",
    "date":"2013-11-12T23:59:59.000Z",
    "invisible": false,
    "abstract":"So you've got this nice recursive function that you love and cherish, but that's incredibly prohibitive to compute.  How can you keep it and make it much faster?  Dynamic programming is said to be hard; but with enough lazyness, you too can do it.",
    "language":"en"
}

<script type="application/javascript" src="/assets/js/algo_convenience_hacks.js"></script>

> Preambule:  the algorithms in this page are all implemented in Javascript, and your browser is currently using them to generate the numbers you'll see below. Please note that I'm not a good Javascript dev and the code I've written is hackish at best.

## You Love Recursion

Hello Jim, I've been told that you quite enjoy recursion.

You've got this nice recursive function that you love and cherish, but oh boy it's an incredibly prohibitive one to compute! How can you keep the little recursive thing and make it much faster?

You've been told that dynamic programming is the Graal. You've been told that dynamic programming hard.

However, let me tell you that with enough lazyness, you too can do it while not changing much to your lovely recuring thing.

## The Mighty Change Giving Algorithm

Your professor Ms. Computer Science gave you a cute little thing.  It's not really important to know what this is supposed to do, but she said it computes the change needed to form an amount $n$ with $10$, $5$ and $1$ cents coins (we trust her with that):

```javascript
function getChange(n) {
    recursiveChange(n, 3) // base case is k=3
}

// helper
function recursiveChange (n, k) {

    if (k == 1) return 1;

    var big = 0;
    if (k == 3) { big = 10; }
    else        { big = 5;  }

    if (n < big) {
        return recursiveChange(n, k-1);
    }

    return recursiveChange(n - big, k) + recursiveChange(n, k-1);
}
```

<!-- This follows from above -->
<script type="text/javascript">
function getChange(n) {
    var start = new Date().getTime();
    var countOps = 0;
    function recursiveChange (n, k) {
        countOps++;
        if (k == 1) return 1;

        var big = 0;
        if (k == 3) { big = 10; }
        else        { big = 5;  }

        if (n < big) {
            return recursiveChange(n, k-1);
        }

        return recursiveChange(n - big, k) + recursiveChange(n, k-1);
    }
    return {answer: recursiveChange(n, 3),
            recursions: countOps,
            time: (new Date().getTime() - start)};
}
</script>

Whether the algorithm actually does compute the proper change (or not!) is quite irrelevant to our problem.  We just want to make this thing FASTER!  Fast like a [Rockomax "Mainsail"](http://wiki.kerbalspaceprogram.com/wiki/Rockomax_%22Mainsail%22_Liquid_Engine) rocket engine!

Wow, that's pretty fast!  How are we going to do this, Antoine?

Oh well, I'm not too intelligent, but I have a trick!  It's called __lazy loading__ and I've heard it's quite the cool and hip thing.  But first, let's look at how this thing grows!

<script type="text/javascript">
var t = VerticalTable([
    "$n$",
    "Answer",
    "Recursions done"
]);

for (var i = 1; i <= 4096; i*= 2) {
    var r = getChange(i).answer;
    var rOps = getChange(i).recursions;
    t.addEntry([
        math(i),
        math(r),
        math(rOps),
    ]);
}
puts(t.toHTML());
</script>

Well, from the look of it, I'd say this is growing pretty FAST!  But that's not the 'fast' we're looking for.  We want this _Recursions done_ column to be little, so that the algorithm itself is _FAST_!

> Note : All along the post, I've assumed that _faster_ means less recursions.  Let's just assume that this is true in this algorithm, shall we?

## So What Is Lazy Loading Anyways?

You might have seen the lazy loading pattern in other shapes.  It's quite used in Java (and everything else!) to avoid loading expensive objects until they're very needed.  Here's an example (the only part of this post that's written) in Java:

```Java
// hugeString is huge, so don't store it by default.
private String hugeString = null;

// People ask for the huge string using this method, always
public String getHugeString() {

    // If nobody ever asked for hugeString, it's null so we get it
    if (hugeString == null) {
        hugeString = Database.getHugeString();
        // from now on, hugeString is set
    }

    // In any case, we return hugeString, whether it was already there or not
    return hugeString;
}
```
That's pretty simple, right?  Let's see how this applies to our $getChange$ algorithm.


## Identifying Key Informations

There are a few things you need to _just know_ about dynamic programming.  Here they are:

### 1. Thou shall need a multi-dimensional array.
$$
    dynamic\ programming \to array
$$

The answer to _question 1_ is pretty easy: need I use an array?  Yes!

### 2. Thou shall know how many dimensions thy array will have, somehow.

$$
    array[i][j][k]...[x]
$$

Let's look at _question 2_.  To do that, recall the first part of the algorithm:

```javascript
function recursiveChange (n, k) {
    ...
}
```

You can see here that function `recursiveChange` takes two arguments, `n` and `k`.  So the answer to _question 2_ is a 2-dimensional array.  That's pretty easy so far.

### 3. Thou shall know what size thy dimensions will be, somehow.
$$\begin{align}
    i&=something\\\
    j&=somethingElse\\\
    &\cdots \\\
    x&=somethingAtLast
\end{align}$$

Now, the harder part; _question 3_.  We need to know how big each argument will be, because that's how big we need to make the dimensions.  So we just look at the code and figure it out with our dumb little brains:

First let's look at the easy one.

#### How big will `k` be?

We see in the `getChange` function that `k` starts with value $3$.

```javascript
function getChange(n) {
    recursiveChange(n, 3) // base case is k=3
}
```

Now we look if anything in the helper function `recursiveChange` will ever make `k` bigger:

```javascript
function recursiveChange (n, k) {

    if (k == 1) return 1;

    var big = 0;
    if (k == 3) { big = 10; }
    else        { big = 5;  }

    if (n < big) {
        return recursiveChange(n, k-1);
    }

    return recursiveChange(n - big, k) + recursiveChange(n, k-1);
}
```
Nope, never.  So $k\leq3$ in any case.

#### How big will `n` be?

We know that the algorithm starts with initial value `n`, then recurses with equal or smaller values than `n`.  So `n` will be as big as `n` is.  This sounds weird, but what it means is that the array will be dynamically sized by a constant `3` and a varying `n`.

### Putting it 1, 2 and 3 together.

So we need a 2-dimensional array of size $n \times 3$.  Let's call this array `memo`, short for _memoization_; a fancy term intelligent people use to show that they have a fancy vocabulary.  I use it too, just to feel smart and _distingué_.  You should do that too from now on.

```javascript
memo = new Array(n);
for (var i = 0; i < memo.length; i++) {
    memo = new Array(3);
}
```
> Note: actually, it's `n + 1` and `3 + 1` because prof Ms. Computer Science likes her algorithms 1-indexed.

You can clearly see that the _space complexity_ of your yet-to-be-born algorithm will somehow be bounded below by `n`.  Actually, at least $3n$.  The fancy people say $\Omega(n)$.

## Cool Bro, But Where's My Rockomax-Fast Algorithm?

Well from now on '_bro_', it will be even easier. Apply the following recipe:

* At the beginning of `recursiveChange`, look into `memo` if you don't know the answer already.  If not, change nothing.

```javascript
if (memo[n][k]) { // we cheat; javascript considers `null` to be false.
    return memo[n][k];
}
```

* For every recursive call to `recursiveChange(i, j)`, check if `memo[i][j]` is known.
    * If __yes__, return that value instead of doing the recursion.
    * If __no__, do the recursive call, but save the value you get back into `memo[i][j]`.

```javascript
if ( !memo[i][j] ) {
    memo[i][j] = recursiveChange(i, j);
}
// use `memo[i][j]`
```

* When you're about to `return` the computed value at the very end, save it into `memo` first.

```javascript
memo[i][j] = answer;
return answer;
```

## The Quasi End Result

We do as I've said above, and replace all the access to lazy ones, and save whatever we compute at each step.  Here's the result:

```javascript
function getChangeDynamic(n) {
    var k = 3;
    // create our memo array
    var memo = new Array(n);
    for (var i = 0; i < memo.length; i++) {
        memo[i] = new Array(k);
    }
    // call the recursive function as usual
    recursiveChange(n, k, memo)
}

// a helper to clean up the code a bit
function lazyGet(i, j, memo) {
    if ( !memo[i][j] ) {
        memo[i][j] = recursiveChange(i, j, memo);
    }
    return memo[i][j];
}

// Keep the recursive function as is, minus the use of the memoization
// array, for lazy loading
function recursiveChange (n, k, memo) {

    if (k === 1) {
        return 1;
    }

    // if we know the answer, don't compute anything
    if ( memo[n][k] ) {
        return memo[n][k];
    }

    var big = 0;

    if (k === 3) {
        big = 10;
    } else {
        big = 5;
    }

    // lazily compute the values
    if (n < big) {
        return lazyGet(n, k-1, memo);
    }

    var withoutBig = lazyGet(n-big, k, memo);
    var withBig = lazyGet(n, k-1, memo);

    // save answers we have had to compute the long way
    memo[n][k] = withoutBig + withBig
    return memo[n][k];
}
```

<script type="text/javascript">
// Woa man, you inline Javascript in your HTML jsut like that?  Like, wtf dude?
// - I don't care, this is my blog.  Wtv.

function getChangeDynamicSlow(n) {
    var start = new Date().getTime();
    var k = 3;
    // create our memo array
    var memo = new Array(n+1);
    for (var i = 0; i < memo.length; i++) {
        memo[i] = new Array(k+1);
    }

    // a helper to clean up the code a bit
    function lazyGet(i, j) {
        if ( !memo[i][j] ) {
            memo[i][j] = recursiveChange(i, j);
        }
        return memo[i][j];
    }

    var countOps = 0;

    // Keep the recursive function as is, minus the use of the memoization
    // array, for lazy loading
    function recursiveChange (n, k) {
        countOps++;
        if (k === 1) {
            return 1;
        }

        // if we know the answer, don't compute anything
        if ( memo[n][k] ) {
            return memo[n][k];
        }

        var big = 0;

        if (k === 3) {
            big = 10;
        } else {
            big = 5;
        }

        // lazily compute the values
        if (n < big) {
            return lazyGet(n, k-1);
        }

        var withoutBig = lazyGet(n-big, k);
        var withBig = lazyGet(n, k-1);

        // save answers we have had to compute the long way
        memo[n][k] = withoutBig + withBig
        return memo[n][k];
    }

    return {answer: recursiveChange(n, k, memo),
            memo: memo,
            recursions: countOps,
            time: (new Date().getTime() - start)}
}
</script>

So that was our poor-woman and poor-man universal '_dynamization_' technique: use the plain recursive algorithm, and plug in some lazy-loading everywhere.  Is this cheating?  No it's not.  Is this a elegant way?  Hmm maybe, maybe not... but it works incredibly well!

## Show Me The Mumbers!

Starting off, 'mumbers' is not a word.  Now that this is out of the way, let's indeed look at some numbers.  We will want to look at two things:

* __The computed answer__: we want to make sure our algorithm is still computing the right thing, don't you agree?
* __The number of recursions__: we want to see if we've made the thing faster.

I've copy-pasted the algorithm above with some minor changes into the HTML of this page.  Along with some helpers, I'm now saying this:

> "Javascript, compute thy numbers!!!"

Here it goes :

<script type="text/javascript">
var t = VerticalTable([
    "$n$",
    "Original($n$)",
    "Dynamic($n$)",
    "Original recursions",
    "Dynamic recursions"
]);

for (var i = 0; i <= 20; i++) {
    var r = getChange(i).answer;
    var rOps = getChange(i).recursions;
    var d = getChangeDynamicSlow(i).answer;
    var dOps = getChangeDynamicSlow(i).recursions;
    t.addEntry([
        math(i),
        math(r),
        math(d),
        math(rOps),
        math(dOps),
    ]);
}

for (var i = 40; i <= 4096; i*= 2) {
    var r = getChange(i).answer;
    var rOps = getChange(i).recursions;
    var d = getChangeDynamicSlow(i).answer;
    var dOps = getChangeDynamicSlow(i).recursions;
    t.addEntry([
        math(i),
        math(r),
        math(d),
        math(rOps),
        math(dOps),
    ]);
}
puts(t.toHTML());
</script>


# If You Are Truly Lazy, Stop There.

That was it.

# If You Want Faster Than A Rockomax, Carry On.

Now that we've seen some numbers, we can see a pattern.  The number of operations performed changes only every multiple of $5$.  That should be a hint that there's some wastage going on.

Let's see what the `memo` array looks like.  Now, let me warn you of two things:

* I'm NOT a javascript ninja.
* This code is most likely insane.
* I'm not saying that this is the best possible algorithm, and I don't care about the best possible algorithm to _~compute change~_.
* What follows assume you have a bit more of a brain that the previous part.  Which means, I won't do lengthy explanations of every line.

That was not two but four things, great!  Carry on.

Let's abuse Javascript a little bit and instrument our `memo` array to make it more convenient to work with.  Please close your eyes:


```javascript
var memo = new Array(n+1);
for (var i = 0; i < memo.length; i++) {
    memo[i] = new Array(k+1);
};

memo.lazyGet = function(i, j) {

    if (!this[i][j]) {
        this[i][j] = getChangeDynamic(i, j);
    }
    return this[i][j];
}

memo.get = function(i, j) {
    return this[i][j];
}

memo.set = function(i, j, val) {
    this[i][j] = val
    return val; // for chaining
}
```

Alright, now it will be easier to hack around the algorithm and change things.

## Show Me More Numbers!

We want to see more data!  Looking at the `memo` array for some $n$, we see this:

<!-- Wow man, wtf this is insane, why u do this? -->
<table>
<script type="text/javascript">

var n = 21;
var dynamicSlowWith21 = getChangeDynamicSlow(n);

var memo = dynamicSlowWith21.memo;

for (var i = 0; i < memo.length; i++) {
    if (i === 0) {
        puts("<tr><th>$n="+n+"$ </th>")
        for (var j = 0; j < memo[i].length; j++) {
            puts("<th>"+j+"</th>");
        };
        puts("</th>")
    }
    puts("<tr>");
    puts("<th>"+i+"</th>");
    for (var j = 0; j < memo[i].length; j++) {
        puts("<td>"+memo[i][j]+"</td>");
    };
    puts("</tr>");
};
</script>
</table>

This thing is filled with unused cells!

We can see that indeed, every non multiple of $5$ is unused.  We can thus change every access to the `memo` array to only use multiples of $5$, like this:

```javascript
var len = Math.floor(n/5) + 1 // 1-indexed
var memo = new Array(len);
for (var i = 0; i < memo.length; i++) {
    memo[i] = new Array(3 + 1); // 1-indexed
};

memo.lazyGet = function(i, j) {
    var realI = Math.floor(i/5);

    if (!this[realI][j]) {
        this[realI][j] = getChangeDynamic(i, j);
    }
    return this[realI][j];
}

memo.get = function(i, j) {
    var realI = Math.floor(i/5);
    return this[realI][j];
}

memo.set = function(i, j, val) {
    var realI = Math.floor(i/5);
    this[realI][j] = val
    return val; // for chaining
}
```

<script>

function getChangeDynamicMultiple5(n) {

    var len = Math.floor(n/5)
    var memo = new Array(len + 1);
    for (var i = 0; i < memo.length; i++) {
        memo[i] = new Array(4);
    };

    memo.lazyGet = function(i, j) {
        var realI = Math.floor(i/5);

        if (!this[realI][j]) {
            this[realI][j] = lambda(i, j);
        }
        return this[realI][j];
    }

    memo.get = function(i, j) {
        var realI = Math.floor(i/5);

        var row = this[realI];
        if (!row) { return null; }
        return row[j];
    }

    memo.set = function(i, j, val) {
        var realI = Math.floor(i/5);
        this[realI][j] = val
        return val; // for chaining
    }

    var recurCount = 0;

    function lambda (n, k) {
        if (k === 1) {
            return 1;
        }

        if ( memo.get(n,k) ) {
            return memo.get(n,k);
        }

        recurCount++;


        var big = 0;

        if (k === 3) {
            big = 10;
        } else {
            big = 5;
        }

        if (n < big) {
            return memo.lazyGet(n, k-1);
        }

        var withoutBig = memo.lazyGet(n-big, k);
        var withBig = memo.lazyGet(n, k-1);

        return memo.set(n,k, withoutBig + withBig);
    }

    return {answer: lambda(n, 3), memo: memo, recursions: recurCount}
}
</script>

Let's look at the result :

<table>
<script type="text/javascript">

var n = 21;
var dynamicMultiple5 = getChangeDynamicMultiple5(n);

var memo = dynamicMultiple5.memo;

for (var i = 0; i < memo.length; i++) {
    if (i === 0) {
        puts("<tr><th>$n="+n+"$ </th>")
        for (var j = 0; j < memo[i].length; j++) {
            puts("<th>"+j+"</th>");
        };
        puts("</th>")
    }
    puts("<tr>");
    puts("<th>"+i+"</th>");
    for (var j = 0; j < memo[i].length; j++) {
        puts("<td>"+memo[i][j]+"</td>");
    };
    puts("</tr>");
};
</script>
</table>

Wow, that's much better!  But we still see there's a lot of extra `undefined`.  Let's remove it by making `memo` $n \times k-1$, and redirecting all the accesses from `j` to `j-1`:

```javascript
var len = Math.floor(n/5) + 1
var memo = new Array(len);
for (var i = 0; i < memo.length; i++) {
    memo[i] = new Array(3);  // remove the extra 1
};

memo.lazyGet = function(i, j) {
    var realI = Math.floor(i/5);

    if (!this[realI][j-1]) {
        this[realI][j-1] = getChangeDynamic(i, j);
    }
    return this[realI][j-1];
}

memo.get = function(i, j) {
    var realI = Math.floor(i/5);
    return this[realI][j-1];
}

memo.set = function(i, j, val) {
    var realI = Math.floor(i/5);
    this[realI][j-1] = val
    return val; // for chaining
}
```

<script type="text/javascript">
function getChangeDynamicLessOneCol(n) {

    var len = Math.floor(n/5)
    var memo = new Array(len + 1);
    for (var i = 0; i < memo.length; i++) {
        memo[i] = new Array(3);
    };

    memo.lazyGet = function(i, j) {
        var realI = Math.floor(i/5);

        if (!this[realI][j-1]) {
            this[realI][j-1] = lambda(i, j);
        }
        return this[realI][j-1];
    }

    memo.get = function(i, j) {
        var realI = Math.floor(i/5);

        var row = this[realI];
        if (!row) { return null; }
        return row[j-1];
    }

    memo.set = function(i, j, val) {
        var realI = Math.floor(i/5);
        this[realI][j-1] = val
        return val; // for chaining
    }

    var recurCount = 0;

    function lambda (n, k) {
        if (k === 1) {
            return 1;
        }

        if ( memo.get(n,k) ) {
            return memo.get(n,k);
        }

        recurCount++;


        var big = 0;

        if (k === 3) {
            big = 10;
        } else {
            big = 5;
        }

        if (n < big) {
            return memo.lazyGet(n, k-1);
        }

        var withoutBig = memo.lazyGet(n-big, k);
        var withBig = memo.lazyGet(n, k-1);

        return memo.set(n,k, withoutBig + withBig);
    }

    return {answer: lambda(n, 3),
            memo: memo,
            recursions: recurCount}
}
</script>

Let's look at the result :

<table>
<script type="text/javascript">

var n = 21;
var dynamicLessOneCol = getChangeDynamicLessOneCol(n);

var memo = dynamicLessOneCol.memo;

for (var i = 0; i < memo.length; i++) {
    if (i === 0) {
        puts("<tr><th>$n="+n+"$ </th>")
        for (var j = 0; j < memo[i].length; j++) {
            puts("<th>"+j+"</th>");
        };
        puts("</th>")
    }
    puts("<tr>");
    puts("<th>"+i+"</th>");
    for (var j = 0; j < memo[i].length; j++) {
        puts("<td>"+memo[i][j]+"</td>");
    };
    puts("</tr>");
};
</script>
</table>

Cleaned-up that undefined column!

Now we see those silly entries that are invariably $1$.  Let's get rid of them:

```javascript
var len = Math.floor(n/5)
var memo = new Array(len);   // removed the extra 1
for (var i = 0; i < memo.length; i++) {
    memo[i] = new Array(3-1); // removed the extra 1
};

memo.lazyGet = function(i, j) {
    var realI = Math.floor(i/5) - 1;

    // add cases where answer is invariably 1
    if (realI == -1) { return 1;} // realI will be -1 for access to entry 0
    if (j == 1)      { return 1;}

    if (!this[realI][j-2]) {
        this[realI][j-2] = getChangeDynamic(i, j);
    }
    return this[realI][j-2];
}

memo.get = function(i, j) {
    var realI = Math.floor(i/5) - 1;

    var row = this[realI];
    if (!row) { return null; }
    return row[j-2];
}

memo.set = function(i, j, val) {
    var realI = Math.floor(i/5);
    this[realI-1][j-2] = val
    return val; // for chaining
}
```
<script>

function getChangeDynamic(n) {
    var start = new Date().getTime();
    var len = Math.floor(n/5)
    var memo = new Array(len);
    for (var i = 0; i < memo.length; i++) {
        memo[i] = new Array(2);
    };

    memo.lazyGet = function(i, j) {
        var realI = Math.floor(i/5) - 1;

        if (realI == -1) { return 1;}
        if (j == 1)     { return 1;}

        if (!this[realI][j-2]) {
            this[realI][j-2] = lambda(i, j);
        }
        return this[realI][j-2];
    }

    memo.get = function(i, j) {
        var realI = Math.floor(i/5) - 1;

        var row = this[realI];
        if (!row) { return null; }
        return row[j-2];
    }

    memo.set = function(i, j, val) {
        var realI = Math.floor(i/5);
        this[realI-1][j-2] = val
        return val; // for chaining
    }

    var recurCount = 0;

    function lambda (n, k) {
        if (k === 1) {
            return 1;
        }

        if ( memo.get(n,k) ) {
            return memo.get(n,k);
        }

        recurCount++;


        var big = 0;

        if (k === 3) {
            big = 10;
        } else {
            big = 5;
        }

        if (n < big) {
            return memo.lazyGet(n, k-1);
        }

        var withoutBig = memo.lazyGet(n-big, k);
        var withBig = memo.lazyGet(n, k-1);

        return memo.set(n,k, withoutBig + withBig);
    }

    return {answer: lambda(n, 3),
            memo: memo,
            recursions: recurCount,
            time: (new Date().getTime() - start)}
}
</script>

And at last, let's look at the resulting table:

<table>
<script type="text/javascript">

var nIs20 = getChangeDynamic(21);

var memo = nIs20.memo;

for (var i = 0; i < memo.length; i++) {
    if (i === 0) {
        puts("<tr><th></th>")
        for (var j = 0; j < memo[i].length; j++) {
            puts("<th>"+j+"</th>");
        };
        puts("</th>")
    }
    puts("<tr>");
    puts("<th>"+i+"</th>");
    for (var j = 0; j < memo[i].length; j++) {
        puts("<td>"+memo[i][j]+"</td>");
    };
    puts("</tr>");
};
</script>
</table>

## How Is This Better Than The First Algorithm?

Let's compare the original algorithm with the initial dynamic algorithm and this new, shiny one.  We compare it's values with the original recursive only algorithm, but we omit the initial dynamic version (we've already shown it equivalent).

<button id="long_computation">Compute this expensive table!</button>

<div id="expensive_table"></div>

<script type="text/javascript">
document.getElementById("long_computation").onclick = function() {
    var t = VerticalTable([
        "$n$",
        "Recursive",
        "Dynamic Fast",
        "Original cost",
        "Dynamic slow cost",
        "Dynamic fast cost",
        "Original $dT$",
        "Dynamic slow $dT$",
        "Dynamic fast $dT$"
    ]);
    for (var i = 1; i <= 40000; i*= 2) {
        var r = getChange(i);
        var d = getChangeDynamicSlow(i);
        var dFast = getChangeDynamic(i);
        t.addEntry([
            math(i),
            math(r.answer),
            math(dFast.answer),
            math(r.recursions),
            math(d.recursions),
            math(dFast.recursions),
            r.time + "ms",
            d.time + "ms",
            dFast.time + "ms"
        ]);
    }
    document.getElementById("expensive_table").innerHTML = t.toHTML();
    var table = document.getElementById("expensive_table");
    MathJax.Hub.Queue(["Typeset",MathJax.Hub,table]);

};

</script>


Yep, much better!  Both the initial `getDynamicChange` and the optimized one are growing by a factor of $n$.  However, the inital one grows by $n/2$ while the optimized one grows by $n/4$.  It's a 100% speed-up.


<script type="text/javascript" src="http://cdn.mathjax.org/mathjax/latest/MathJax.js?config=TeX-AMS-MML_HTMLorMML"></script>

<script type="text/javascript">
// Single $ for inline LaTeX
MathJax.Hub.Config({
  tex2jax: {inlineMath: [['$','$'], ['\\(','\\)']]}
});
</script>
