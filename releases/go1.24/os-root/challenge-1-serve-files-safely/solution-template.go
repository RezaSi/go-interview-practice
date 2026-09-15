package main

import (
	"fmt"
	"os"
	"path/filepath"
)

// FileServer serves files out of a single directory
type FileServer struct {
	BaseDir string
}

// ReadUnsafe is the bug. It joins the user's name onto the base directory and
// reads whatever comes out, so "../../etc/passwd" leaves the directory.
// It is here as the reference the tests compare against. Do not change it.
func (s FileServer) ReadUnsafe(name string) ([]byte, error) {
	return os.ReadFile(filepath.Join(s.BaseDir, name))
}

// ReadSafe reads a file from BaseDir and refuses anything that resolves outside it
func (s FileServer) ReadSafe(name string) ([]byte, error) {
	// TODO: open BaseDir with os.OpenRoot, close it, and read name through the root
	return nil, nil
}

// WriteSafe writes a file inside BaseDir and refuses anything that resolves outside it
func (s FileServer) WriteSafe(name string, data []byte) error {
	// TODO: open BaseDir with os.OpenRoot, close it, and write through the root
	return nil
}

func main() {
	dir, _ := os.MkdirTemp("", "srv")
	defer os.RemoveAll(dir)
	os.WriteFile(filepath.Join(dir, "notes.txt"), []byte("hello"), 0o644)

	s := FileServer{BaseDir: dir}

	b, err := s.ReadSafe("notes.txt")
	fmt.Printf("inside: %q err=%v\n", b, err)

	_, err = s.ReadSafe("../../etc/passwd")
	fmt.Println("escape refused:", err != nil)
}
