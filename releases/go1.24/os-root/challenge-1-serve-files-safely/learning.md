# Learning Materials for Serve Files Safely

## Path Traversal in One Paragraph

If a path component comes from outside your program and you join it to a base
directory, the result is not guaranteed to stay in that directory. `..` walks up,
absolute paths ignore the base entirely, and a symlink can point anywhere. An attacker
who controls the name controls which file you open.

## Why the Usual Fixes Leak

### Rejecting ".."

```go
if strings.Contains(name, "..") { return errors.New("nope") }
```

Catches the obvious probe. Misses a symlink inside the directory that points at `/etc`.
Also rejects the innocent file `notes..txt`.

### Checking the prefix after joining

```go
p := filepath.Join(base, name)
if !strings.HasPrefix(p, base) { return errors.New("nope") }
```

Two problems. `/srv/files-secret` has `/srv/files` as a string prefix, so a sibling
directory passes. And the check runs *before* the open, so anything that changes the
filesystem in between (a swapped symlink) defeats it. That gap is a TOCTOU race.

Both fixes validate a string. The syscall then resolves a path. Those are different
things, and the difference is the vulnerability.

## What os.Root Does Instead

```go
root, err := os.OpenRoot(baseDir)
defer root.Close()

f, err := root.Open(name)
```

`os.Root` holds the directory **open** and resolves the path one component at a time
relative to it, rather than building a string and handing it to the kernel. On Linux
that is `openat` with `RESOLVE_BENEATH`; elsewhere Go performs the equivalent walk.

Two properties fall out:

- **Symlinks are followed and then judged.** A link whose target is outside the root
  fails. A link to a sibling inside the root still works.
- **There is no check-then-use gap.** The directory descriptor is the boundary, so
  nothing can be swapped in between validating and opening.

The error is specific:

```text
openat ../../etc/passwd: path escapes from parent
```

## The API

`*os.Root` mirrors the `os` package for operations inside the directory:

```go
root.Open(name)             root.Create(name)
root.OpenFile(name, f, m)   root.ReadFile(name)
root.WriteFile(name, b, m)  root.Stat(name)
root.Lstat(name)            root.Mkdir(name, m)
root.Remove(name)           root.Rename(old, new)
```

Go 1.25 added more (`MkdirAll`, `RemoveAll`, `Chmod`, `Chown`, `Symlink`, `Readlink`
and friends), so if a method is missing, check the version in your `go.mod`.

## Holding the Root

Open one root per served directory at start-up and keep it:

```go
type Server struct{ root *os.Root }

func New(dir string) (*Server, error) {
	r, err := os.OpenRoot(dir)
	return &Server{root: r}, err
}
```

Opening a root per request costs a syscall and a file descriptor each time, and the
open descriptor is the thing providing the guarantee. This challenge opens per call to
keep the code short, but a real server should not.

## What It Does Not Cover

- **Permissions.** If your process can write in the directory, so can any code holding
  the root.
- **Name hygiene.** Length limits, null bytes, reserved Windows names, case-insensitive
  collisions: still yours to handle.
- **What is in the file.** It controls *where* you land, not *what* you do there.

## Where This Matters Most

Archive extraction. Every pre-1.24 `zip` or `tar` extractor hand-rolled traversal
checks and a good number got them wrong, which is the "zip slip" family of CVEs. The
modern shape:

```go
root, err := os.OpenRoot(dest)
defer root.Close()

for _, entry := range archive.Files {
	f, err := root.Create(entry.Name) // refuses ../../ and symlink escapes
	...
}
```

## Further Reading

- [Go 1.24 release notes: directory-limited filesystem access](https://go.dev/doc/go1.24#directory-limited-filesystem-access)
- [Proposal: os: add Root](https://github.com/golang/go/issues/67002)
- [CWE-22: path traversal](https://cwe.mitre.org/data/definitions/22.html)
