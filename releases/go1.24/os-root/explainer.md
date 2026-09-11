## The bug that keeps getting written

You are serving files out of a directory. The name comes from the user.

```go
func serve(base, name string) ([]byte, error) {
	return os.ReadFile(filepath.Join(base, name))
}
```

`filepath.Join` cleans the path, so `serve("/srv/files", "../../etc/passwd")` resolves
to `/etc/passwd` and hands it straight back. This is CWE-22, it is the single most
common vulnerability class in file-serving code, and Go has never had a good answer
for it.

The usual fixes all have holes:

```go
// 1. reject ".." in the input
if strings.Contains(name, "..") { return nil, errors.New("nope") }
```

Blocks the obvious attempt. Does nothing about a **symlink** inside `base` that points
at `/etc`. Also rejects the perfectly innocent file `notes..txt`.

```go
// 2. check the joined path has the right prefix
p := filepath.Join(base, name)
if !strings.HasPrefix(p, base) { return nil, errors.New("nope") }
```

Better, and still wrong twice over. `/srv/files-secret` has `/srv/files` as a string
prefix. And the check happens *before* the open, so a symlink swapped in between the
check and the read walks straight out. That race has a name, TOCTOU, and it is not
theoretical.

The real problem is that all of these validate a **string** and then hand a different
thing, a **path resolved by the kernel**, to the syscall.

## What Go 1.24 added

`os.OpenRoot` returns an `*os.Root`: a handle to an open directory that will not let
you out of it.

```go
root, err := os.OpenRoot("/srv/files")
if err != nil {
	return err
}
defer root.Close()

f, err := root.Open(name)
```

If `name` resolves anywhere outside that directory, the call fails:

```text
openat ../../etc/passwd: path escapes from parent
```

No string comparison anywhere in your code. The check is part of the resolution.

`*os.Root` mirrors the `os` API you already know: `Open`, `Create`, `OpenFile`,
`Stat`, `Lstat`, `Mkdir`, `Remove`, `Rename`, `ReadFile`, `WriteFile`, `RemoveAll`,
`MkdirAll`, `Symlink`, `Readlink`. Go 1.25 filled in most of the remaining gaps, so
if a method you want is missing, check which version you are building with.

## Why it actually holds

This is the part worth understanding, because it explains why `os.Root` is safe in
places where a prefix check is not.

A traversal check on a string asks: *does this text look like it stays inside?*
`os.Root` asks a different question at every step: *is this component still inside the
directory I have open?*

It walks the path one component at a time, opening each relative to the previous one
rather than re-joining strings. On Linux that is `openat` with `RESOLVE_BENEATH`; on
other platforms Go does the equivalent walk itself. Two consequences fall out:

**Symlinks are followed and then judged.** A link inside the root pointing at
`/etc/hosts` resolves to a target outside the boundary, so the operation fails. You do
not have to detect links or refuse to follow them; a link to a sibling file inside the
root still works fine.

**There is no gap between checking and using.** The directory is held open as a file
descriptor. Renaming or replacing a component after the check does not help an
attacker, because there is no "after the check": the check *is* the resolution.

That is also why you should hold onto the `*os.Root` rather than reopening it per
request. The open descriptor is the security boundary.

## What it does not do

`os.Root` is about *where* the path lands, not *what* you are allowed to do there.

- It does not check permissions. If your process can write inside the directory, so
  can any handler that has the root.
- It does not stop absolute paths from being nonsense: `root.Open("/etc/passwd")` is
  interpreted relative to the root, and fails, which is what you want, but it is not
  reading your mind.
- It does not sanitise names for other purposes. A file called `../weird` on Windows,
  a name with a null byte, a name that is 4KB long: still your problem.
- It costs a file descriptor per root. Open one per served directory at start-up, not
  one per request.

## When to reach for it

Any time a path component comes from outside your program. Uploaded file names, an
archive being extracted, a template path from config, a plugin directory scan, a
static file handler.

Extraction is the case worth calling out. Every `zip` and `tar` extractor written
before 1.24 had to hand-roll traversal checks, and a good number got it wrong. The
pattern now:

```go
root, err := os.OpenRoot(destDir)
if err != nil {
	return err
}
defer root.Close()

for _, entry := range archive.Files {
	f, err := root.Create(entry.Name) // refuses ../../ and symlink escapes
	if err != nil {
		return err
	}
	// ...
}
```

If you are on Go 1.24 or newer and you are still calling `filepath.Join` on
user-supplied names, that is the code to change first.

> **Try it:** the two challenges below start with a file handler that has the classic
> bug, then go after the sneakier version of it, where the escape is a symlink planted
> inside the directory rather than a `..` in the request.
