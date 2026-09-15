## The argument that took two years

Proposal [#67002](https://github.com/golang/go/issues/67002) was filed in 2024, but
the problem had been circling the issue tracker for far longer. Several earlier
attempts stalled, and the reasons they stalled are the reasons the final API looks the
way it does.

**"Just add a safe join function."** The obvious ask is
`filepath.SafeJoin(base, name)` returning an error when the result escapes. It cannot
work. A pure string function has no idea what is on disk, so it cannot see that
`uploads/logo.png` is a symlink to `/etc`. Anything that answers using only the text
is guessing.

**"Then make it consult the filesystem."** Now the function stats the path, decides it
is fine, and returns a string. Your code then opens that string. Between those two
steps the filesystem can change. That gap is the TOCTOU race, and a checking function
that returns a path cannot close it, no matter how careful the check is.

The conclusion both of those dead ends point to: the safety property cannot live in a
function that returns a **path**. It has to live in something that returns an **open
file**, because only then is there no window between deciding and doing.

That is why the API is a type with methods rather than a helper function, and why
`os.Root` holds an open directory rather than remembering a string.

## Why the name is `Root` and not `Dir` or `FS`

Two rejected shapes are worth knowing about.

`io/fs.FS` already exists and is already scoped to a tree, so why not use it?
Because `fs.FS` is read-only and abstract. It has no `Create`, no `Rename`, no
`Remove`, and its implementations are often not backed by a real directory at all.
`os.Root` is deliberately concrete: it is one real directory, held open, with the
mutating operations included.

The other shape was a `chroot`-style call that changes the process's view. That is
process-wide, usually needs privileges, and cannot be scoped to one request handler.
`os.Root` is a value you can hold per server, per tenant, per anything.

## What the type actually is

From `src/os/root_openat.go`, the platform-specific implementation is small:

```go
type root struct {
	name string

	// refs is incremented while an operation is using fd.
	// closed is set when Close is called.
	// fd is closed when closed is true and refs is 0.
	mu     sync.Mutex
	fd     sysfdType
	refs   int  // number of active operations
	closed bool // set when closed
}
```

One file descriptor, a mutex, and a reference count. The descriptor is the security
boundary; everything else exists so that `Close` during a concurrent `Open` does not
close the descriptor out from under it. Each operation calls `incref` first and
`decref` when it finishes, and the real `syscall.Close` only happens once `closed` is
set and `refs` has fallen to zero.

## The resolution loop

Every method funnels into one generic helper:

```go
func doInRoot[T any](
	r *Root,
	name string,
	openDirFunc func(parent sysfdType, name string) (sysfdType, error),
	f func(parent sysfdType, name string) (T, error),
) (ret T, err error)
```

It splits `name` into components and walks them, keeping a `dirfd` that starts at the
root's descriptor. Each intermediate component is opened *relative to the previous
descriptor*, so the kernel never sees a reassembled string. When it reaches the final
component it calls `f`, which is where the actual operation (`open`, `mkdir`, `stat`)
happens.

Three details in that loop are worth reading closely.

### `..` restarts from the root

The intuitive implementation of `..` is `openat(dirfd, "..")`. The standard library
refuses to do that, and the comment says why:

```go
// When resolving .. path components, we restart path resolution from the root.
// (We can't openat(dir, "..") to move up to the parent directory,
// because dir may have moved since we opened it.)
```

If someone renames a directory you are holding open, walking `..` from it lands
somewhere you never agreed to. So instead the loop rewrites the path, deleting the
component that the `..` cancels, resets `dirfd` back to the root descriptor, and
starts again.

The escape check falls straight out of that rewrite:

```go
count := end - i
if count > i {
	return ret, errPathEscapes
}
```

If there are more `..` components than there are components before them, the path is
trying to climb above the root, and you get `path escapes from parent`. No string
comparison, no prefix test. It is arithmetic on the component list.

### There are budgets

```go
const maxSteps = 255
const maxRestarts = 8
```

A path made of alternating symlinks and `..` components can force an unbounded number
of opens. The loop counts steps and restarts and gives up with `ENAMETOOLONG` when
both budgets are exhausted. This is a denial-of-service guard, and it is the kind of
detail that a hand-rolled traversal check in application code invariably forgets.

### Symlinks are a return value, not a special case

The final-component callback signals a symlink by returning a sentinel error, and the
loop handles the following. That keeps one code path for "resolve this name" instead
of duplicating link handling into every method, and it is why a link pointing inside
the root keeps working while a link pointing outside does not.

## The platform split, and the one honest compromise

There are three implementations in the tree:

| File | Platforms | Mechanism |
|---|---|---|
| `root_openat.go` | unix, windows, wasip1 | walk with `openat`-style calls, holds a descriptor |
| `root_noopenat.go` | platforms without it | `checkPathEscapes` before the operation |
| `root_js.go` | js/wasm | same fallback |

The fallback is documented with unusual candour in the source:

```go
// Due to the lack of openat, checkPathEscapes is subject to TOCTOU races
```

On a platform without `openat`, `os.Root` still rejects traversal, but it does so with
a check before the operation, so the race the API was designed to eliminate is back.
The standard library says so rather than pretending otherwise. On Linux, macOS,
BSD and Windows, which is where your servers are, you get the real guarantee.

## What this buys you as a reader

Two takeaways that generalise beyond this package.

**A security property belongs to whatever performs the operation.** The moment your
API returns a validated *description* of an action instead of performing it, you have
created a window. That applies to path checks, to permission checks, to signed URLs,
to anything.

**Correct code here is not obvious.** The `..`-restart, the step budget, the
descriptor refcount, the symlink sentinel. That is four non-obvious details in about a
hundred lines. Every application that hand-rolled traversal checking had to get all
four right, and the CVE record shows how that went.

## The change itself

The type landed in one commit, [`43d90c6`](https://github.com/golang/go/commit/43d90c6a14e7), reviewed as
[CL 612136](https://go-review.googlesource.com/c/go/+/612136):

```text
os: add Root

Add os.Root, a type which represents a directory and permits performing
file operations within that directory.

For #67002
```

What it touched is a good map of the design:

| File | Why it exists |
|---|---|
| `src/os/root.go` | the exported API and its documentation |
| `src/os/root_openat.go` | the real implementation: the refcounted descriptor and `doInRoot` |
| `src/os/root_noopenat.go` | the fallback for platforms without `openat` |
| `src/os/root_js.go` | the same fallback for js/wasm |
| `src/os/root_nonwindows.go` | unix specifics |
| `src/internal/syscall/windows/at_windows.go` | Windows needed new syscall plumbing to do this at all |
| `api/next/67002.txt` | the API promise, one line per exported symbol |
| `doc/next/6-stdlib/1-os-root.md` | the release note, written with the code |

Two things worth noticing there. The Windows file is new: this could not be built on
existing syscall wrappers, so the change had to add them. And the API and release note
files are part of the same commit rather than a follow-up, which is how the Go project
keeps `api/` honest.

Everything after that commit is filling in methods. `Root.Stat`, `Root.Remove`,
`Root.FS`, `OpenInRoot`, then in Go 1.25 `Root.ReadFile`, `Root.WriteFile`,
`Root.MkdirAll`, `Root.RemoveAll`, `Root.Rename`, `Root.Readlink`, `Root.Lchown` and
more. Each is a separate small CL calling into the same `doInRoot` walk, which is why
the guarantee is identical across all of them: there is one resolution path and every
method goes through it.

## Read the source

It is short and unusually readable:

- [`src/os/root.go`](https://github.com/golang/go/blob/master/src/os/root.go) — the API and docs
- [`src/os/root_openat.go`](https://github.com/golang/go/blob/master/src/os/root_openat.go) — the resolution loop
- [Proposal #67002](https://github.com/golang/go/issues/67002) — the full discussion
