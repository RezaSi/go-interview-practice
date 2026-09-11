# Challenge 1: Serve Files Safely

## The situation

`FileServer` serves files out of one directory. `ReadUnsafe` has the bug that has
shipped in a thousand Go services:

```go
func (s FileServer) ReadUnsafe(name string) ([]byte, error) {
	return os.ReadFile(filepath.Join(s.BaseDir, name))
}
```

`filepath.Join` cleans the path, so `"../../etc/passwd"` resolves to `/etc/passwd` and
gets read. The test suite proves this happens before it asks you to fix it.

## The task

Implement two methods using `os.OpenRoot`.

### `ReadSafe`

```go
func (s FileServer) ReadSafe(name string) ([]byte, error)
```

Read `name` from `BaseDir`. Nested paths like `sub/deep.txt` must still work. Anything
that resolves outside `BaseDir` must return an error.

### `WriteSafe`

```go
func (s FileServer) WriteSafe(name string, data []byte) error
```

Same rules, for writing. Use mode `0o644`.

Leave `ReadUnsafe` alone. The tests use it as the "before" picture.

## Rules

- No string validation. Do not check for `..`, do not compare prefixes. That is the
  whole point: `os.Root` makes those checks unnecessary and it is stricter than they
  are.
- Close the root you open.
- A missing file *inside* the directory should still report a normal not-exist error,
  so `os.IsNotExist` keeps working. Do not turn every failure into a generic error.

## What the tests check

| Test | What it is really about |
|---|---|
| `TestReadSafeReadsFilesInsideTheDirectory` | the happy path still works |
| `TestReadSafeReadsNestedFiles` | subdirectories are fine, only escapes are not |
| `TestReadSafeRefusesTraversal` | proves `ReadUnsafe` escapes, then that `ReadSafe` does not |
| `TestReadSafeRefusesDeepTraversal` | three spellings of the same attack |
| `TestReadSafeReportsMissingFilesNormally` | you did not swallow the real error |
| `TestWriteSafe*` | the same rules apply to writes |

## Run it

```bash
go test -v ./...
```
