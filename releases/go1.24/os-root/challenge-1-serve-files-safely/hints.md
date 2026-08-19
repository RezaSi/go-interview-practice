# Hints for Serve Files Safely

## Hint 1: Open the Directory First

`os.OpenRoot` takes a directory path and gives you a handle:

```go
root, err := os.OpenRoot(s.BaseDir)
if err != nil {
    return nil, err
}
defer root.Close()
```

## Hint 2: The Root Mirrors the os Package

Once you have a root, the methods are the ones you already know, just scoped:

```go
root.ReadFile(name)
root.WriteFile(name, data, 0o644)
root.Open(name)
root.Create(name)
```

## Hint 3: ReadSafe Is Three Lines Plus the Open

Open the root, defer the close, return `root.ReadFile(name)`. There is nothing to
validate, which is the point.

## Hint 4: Return the Error Unchanged

`TestReadSafeReportsMissingFilesNormally` calls `os.IsNotExist(err)`. If you wrap the
error in your own type or replace it with `errors.New("not found")`, that check fails.
Pass it straight through.

## Hint 5: If the Write Test Fails

`root.WriteFile(name, data, 0o644)` creates the file relative to the root. If you used
`os.WriteFile` with a joined path anywhere, the traversal test will catch it.

## Hint 6: See the Error for Yourself

Add a `fmt.Println(err)` in the traversal case and run the tests. It reads
`openat ../../etc/passwd: path escapes from parent`, which tells you the check
happened during resolution rather than in a string comparison you wrote.
