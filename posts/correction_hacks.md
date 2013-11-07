{
    "title":"Correction Hacks",
    "author":"Antoine Grondin",
    "date":"2013-11-04T18:00:00.000Z",
    "invisible": false,
    "abstract":"Correcting assignments with the Blackboard software is a real pain.  Here's how I helped myself having a nicer day.",
    "language": "en"
}

Nota:

> This is certainly not the most efficient way to do this.  This is just the way I did it tonight because it made me happy.  It might not even be an intelligent way to do this.  Whatever, it's fun.

I'm doing some TA work for the first year class of C at UofO.  For some insane reason, they decided that all engineers __not__ in a computer related field would learn C in first year.

Meanwhile, all computer related engineers (software and computer engineers), along with computer scientist, learn Java as their first language.  So overall, it's a messed up curriculum, but whatever.  That's not the scope of this post.

## Correction is painful

I didn't quite know that before.  I'm sure it could be made better, but the software used at UofO to manage classes is quite painful to deal with.  It tries to keep you in its web UI and work from there.  Problem is:

* It doesn't support most filetypes.
* Its UI is unusable.
* It's just not convenient.

Because they want you to use their web app, they didn't put much care into making offline work usable.  All you get, pretty much, is a giant dump of all the 192 students' files and you have to navigate through that.

So why a hack?

## Because Blackboard S*cks

I don't want to spend more time than required doing this correction thing. So I developped a basic workflow.

Here's how it goes:

* Download all the student files
* Dump them in a folder.
* Create `done/`, `todo/` and `wip/` folders.
* Send everything to `todo/`
* For every student ID:
    * Move its files to `wip/`
    * Correct them (compile, run in a VM, look at the report)
    * Move the files to `done/`

It sounds simple, but think about the cost of every loop iteration (that not any real programming language):

```bash
mv *studentID* wip/
cd wip/
if [ `cwd`.Contains(zipFile) ]
    unzip zipFile
open *.doc *.docx *.odt *.pdf
cc *.c -Wall
./a.out
doTheCorrectionDance()
rm a.out
cd ..
mv wip/* done/
```

That's quite a lot of typing in your terminal, and it's annoying.  There are things you can inline:

```bash
unzip *.zip &&
    open *.doc *.docx *.odt *.pdf &&
    cc *.c -Wall                  &&
    ./a.out                       &&
```
So that at every iteration, you just `arrow-up` to this command and run it.  You're in a VM so you don't really care about running random code.  If something fucks up, you just kill the VM.

The biggest problem remain, there are a lot of little stupid commands to do in order to move the proper student's files into the proper folders at the right time.

## Is There A Solution Out Of This Nightmare

When you download a dump of assignments from Blackboard, the students' files share a common filename format: they all contain their student number in the prefix.

```
Assignment\ 2_XXXXXXX_Tentative_2013-10-27-17-28-14_asg\ 2\ final.docx
```

A simple tool to help you do the correction would pick all the file in `wip/`, move them to `done/`, then pick the files with the `XXXXXXX` student ID string in `todo/` and move them into `wip/`.  Just doing that would save you all these commands:

```
mv *studentID* wip/
cd wip/
cd ..
mv wip/* done/
```

Which aren't easily inlinable in a simple `bash` onliner.  I mean, I'm sure there's a way to do this with some bash magic.  I'm not a bash magician, so I wrote a Go _script_ to do this.

I'll use a couple of Go packages, here they are:

```go
import (
    // To print on the screen
    "fmt"
    // A package I made to colorize strings
    "github.com/aybabtme/color/brush"
    // To rename files
    "os"
    // To create filenames and clean them
    "path"
    // To walk the directories
    "path/filepath"
    // To match the studentIDs in the filenames
    "regexp"
)
```

First I'll need to keep track of what are the files in `wip/` and `todo/`.  I'll index those by `studentID` so I can easily get all the files of a specific student.

```go
todoFiles := make(map[string][]string)
wipFiles := make(map[string][]string)
```

Next, I need to have a way to get the student IDs that are in play.  This simple regexp will do, since student IDs are all 7 digits long.

```go
regexpStdID := regexp.MustCompile(`\d{7}`)
```

Now, I need to collect the files that lie in `todo/`, indexed by `studentID`:

```go
filepath.Walk("todo/", func(path string, fi os.FileInfo, err error) error {
    if !regexpStdID.Match([]byte(fi.Name())) {
        return nil
    }
    stdID := regexpStdID.FindString(fi.Name())
    todoFiles[stdID] = append(todoFiles[stdID], "todo/"+fi.Name())
    return nil
})
```

We do the same thing with `wip/`:

```go
// Collect all the files in `wip`
filepath.Walk("wip/", func(path string, fi os.FileInfo, err error) error {
    // Do the same thing
})
```

Then, we move all the files we've found in `wip/` to `done/`.

```go
// Move all the files in `wip` to `done`
for _, filenamesToMove := range wipFiles {
    for _, wipFilename := range filenamesToMove {
        doneFilename := path.Join("done", path.Base(wipFilename))
        err := os.Rename(wipFilename, doneFilename)
        if err != nil {
            fmt.Printf("Didn't work, %v\n", err)
            continue
        }
    }
}
```
> _In fact, this part could simply have been to move anything in `wip/` to `done/`, there's no need to actually collect the student IDs for `wip/`._

Now that `wip/` is emptied, it's time to find the next student to grade:

```go
// Find the student number of the first document in `todo`
var firstTodo string
for todoKey := range todoFiles {
    firstTodo = todoKey
    break
}
```

If we don't have anything to do, print it in <span style="font-color: red;">red</span>!  So we clearly know we're done.

```go
// Didn't find anything, then we're done
if firstTodo == "" {
    fmt.Println(brush.Red("Nothing to do"))
    return
}
```

Then move all those student's files to `wip/`:

```go
// For every file with this student ID, move it from `todo` to `wip`
for _, filename := range todoFiles[firstTodo] {

    destination := path.Join("wip", path.Base(filename))

    err := os.Rename(filename, destination)
    if err != nil {
        fmt.Printf("Didn't work, %v\n", brush.Red(err.Error()))
        continue
    }
}
```

That's it!  Now, all you have to do is run this script everytime you're done with a student.
