package main

import (
	"log"
	"os"
	"path/filepath"
	"slices"
	"strings"
)

type nav struct {
	currDir   *file
	files     []*file
	cursorPos int
}

func (n *nav) init() {
	n.files = make([]*file, 0)
	n.chDir("/home/grig/sources/transcend")
}

func (n *nav) cursorPrev(skipHidden bool) {
	for i := n.cursorPos - 1; i >= 0; i-- {
		file := n.files[i]
		if file.isHidden() == false || skipHidden == false {
			n.cursorPos = i
			break
		}
	}
}

func (n *nav) cursorNext(skipHidden bool) {
	for i := n.cursorPos + 1; i < len(n.files); i++ {
		file := n.files[i]
		if file.isHidden() == false || skipHidden == false {
			n.cursorPos = i
			break
		}
	}
}

func (n *nav) upDir() {
	prevDir := n.currDir
	n.chDir(n.currDir.parentPath())
	for i, f := range n.files {
		if f.path == prevDir.path {
			n.cursorPos = i
			break
		}
	}
}

func (n *nav) intoDir() {
	curr := n.files[n.cursorPos]
	if curr.IsDir() {
		n.chDir(filepath.Join(n.currDir.path, curr.Name()))
	}
	n.cursorPos = 0
}

func (n *nav) cursorFile() *file {
	return n.files[n.cursorPos]
}

func (n *nav) chDir(path string) {
	if n.currDir != nil && n.currDir.path == path {
		return
	}

	var err error

	n.currDir, err = newFile(path)
	if err != nil {
		log.Fatal(err)
	}

	err = os.Chdir(path)
	if err != nil {
		log.Fatal(err)
	}

	n.files, err = readdir(path)
	if err != nil {
		log.Fatal(err)
	}
}

func readdir(path string) ([]*file, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	names, err := f.Readdirnames(-1)
	f.Close()

	files := make([]*file, 0, len(names))
	for _, fname := range names {
		file, err := newFile(filepath.Join(path, fname))
		if err == nil {
			files = append(files, file)
		}
	}

	slices.SortFunc(files, sortFunc)

	return files, err
}

type file struct {
	os.FileInfo
	path string
}

func newFile(path string) (*file, error) {
	lstat, err := os.Lstat(path)
	if err != nil {
		return nil, err
	}

	return &file{
		FileInfo: lstat,
		path:     path,
	}, nil
}

func (f *file) parentPath() string {
	return filepath.Dir(f.path)
}

// extension without dot
func (f *file) ext() string {
	ext := filepath.Ext(f.Name())
	if ext == f.Name() {
		return ""
	}
	return strings.TrimPrefix(ext, ".")
}
