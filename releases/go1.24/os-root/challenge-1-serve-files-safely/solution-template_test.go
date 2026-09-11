package main

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func newServer(t *testing.T) FileServer {
	t.Helper()
	dir := t.TempDir()
	if err := os.WriteFile(filepath.Join(dir, "notes.txt"), []byte("hello"), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.Mkdir(filepath.Join(dir, "sub"), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(dir, "sub", "deep.txt"), []byte("deep"), 0o644); err != nil {
		t.Fatal(err)
	}
	return FileServer{BaseDir: dir}
}

func TestReadSafeReadsFilesInsideTheDirectory(t *testing.T) {
	s := newServer(t)

	got, err := s.ReadSafe("notes.txt")
	if err != nil {
		t.Fatalf("ReadSafe(notes.txt) returned an error: %v", err)
	}
	if string(got) != "hello" {
		t.Fatalf("ReadSafe(notes.txt) = %q, want %q", got, "hello")
	}
}

func TestReadSafeReadsNestedFiles(t *testing.T) {
	s := newServer(t)

	got, err := s.ReadSafe("sub/deep.txt")
	if err != nil {
		t.Fatalf("ReadSafe(sub/deep.txt) returned an error: %v", err)
	}
	if string(got) != "deep" {
		t.Fatalf("ReadSafe(sub/deep.txt) = %q, want %q", got, "deep")
	}
}

func TestReadSafeRefusesTraversal(t *testing.T) {
	s := newServer(t)

	// Prove the bug is real first: the unsafe version happily leaves the directory.
	outside := filepath.Join(filepath.Dir(s.BaseDir), "outside.txt")
	if err := os.WriteFile(outside, []byte("secret"), 0o644); err != nil {
		t.Fatal(err)
	}
	escape := "../" + filepath.Base(outside)

	if got, err := s.ReadUnsafe(escape); err != nil || string(got) != "secret" {
		t.Fatalf("ReadUnsafe should have escaped and read %q (got %q, err %v)", "secret", got, err)
	}

	// The safe version must not.
	if _, err := s.ReadSafe(escape); err == nil {
		t.Fatal("ReadSafe followed a ../ path out of the directory; os.Root must refuse it")
	}
}

func TestReadSafeRefusesDeepTraversal(t *testing.T) {
	s := newServer(t)

	for _, name := range []string{
		"../../etc/passwd",
		"sub/../../etc/passwd",
		"./../../etc/passwd",
	} {
		if _, err := s.ReadSafe(name); err == nil {
			t.Errorf("ReadSafe(%q) succeeded; it must fail", name)
		}
	}
}

func TestReadSafeReportsMissingFilesNormally(t *testing.T) {
	s := newServer(t)

	_, err := s.ReadSafe("nope.txt")
	if err == nil {
		t.Fatal("ReadSafe on a missing file should return an error")
	}
	if !os.IsNotExist(err) {
		t.Fatalf("want a not-exist error for a missing file inside the root, got %v", err)
	}
}

func TestWriteSafeWritesInsideTheDirectory(t *testing.T) {
	s := newServer(t)

	if err := s.WriteSafe("new.txt", []byte("written")); err != nil {
		t.Fatalf("WriteSafe returned an error: %v", err)
	}

	got, err := os.ReadFile(filepath.Join(s.BaseDir, "new.txt"))
	if err != nil {
		t.Fatalf("the file was not created inside BaseDir: %v", err)
	}
	if string(got) != "written" {
		t.Fatalf("file contains %q, want %q", got, "written")
	}
}

func TestWriteSafeRefusesTraversal(t *testing.T) {
	s := newServer(t)
	target := "../escaped.txt"

	if err := s.WriteSafe(target, []byte("nope")); err == nil {
		t.Fatal("WriteSafe wrote outside the directory; os.Root must refuse it")
	}

	if _, err := os.Stat(filepath.Join(filepath.Dir(s.BaseDir), "escaped.txt")); err == nil {
		t.Fatal("a file was created outside BaseDir")
	}
}

func TestErrorMentionsTheEscape(t *testing.T) {
	// Not a behaviour requirement, but a useful thing to see: os.Root says exactly
	// what went wrong rather than returning a generic permission error.
	s := newServer(t)

	_, err := s.ReadSafe("../../etc/passwd")
	if err == nil {
		t.Fatal("expected an error")
	}
	if !strings.Contains(err.Error(), "escapes") {
		t.Logf("error was %q; os.Root normally reports \"path escapes from parent\"", err)
	}
}
