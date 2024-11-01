package main

import (
	"io/fs"
	"log"
	"os"
	"path"
	"path/filepath"
)

type nav struct {
	currDir string
	entries []fs.FileInfo
	cursor  int
}

func (n *nav) init() {
	n.entries = make([]fs.FileInfo, 0)
	n.chDir("/home/grig/sources/transcend")
}

func (n *nav) cursorPrev() {
	if n.cursor == 0 {
		return
	}

	n.cursor -= 1
}

func (n *nav) cursorNext() {
	if n.cursor == len(n.entries)-1 {
		return
	}

	n.cursor += 1
}

func (n *nav) upDir() {
	prevDir := n.currDir
	n.chDir(path.Dir(n.currDir))
	for i, e := range n.entries {
		if filepath.Join(n.currDir, e.Name()) == prevDir {
			n.cursor = i
			break
		}
	}
}

func (n *nav) intoDir() {
	curr := n.entries[n.cursor]
	if curr.IsDir() {
		n.chDir(filepath.Join(n.currDir, curr.Name()))
	}
	n.cursor = 0
}

func (n *nav) chDir(path string) {
	if n.currDir == path {
		return
	}

	n.currDir = path
	err := os.Chdir(n.currDir)
	if err != nil {
		log.Fatal(err)
	}

	f, err := os.Open(n.currDir)
	if err != nil {
		log.Fatal(err)
	}

	names, err := f.Readdirnames(-1)
	f.Close()

	n.entries = n.entries[:0]
	for _, fname := range names {
		lstat, err := os.Lstat(filepath.Join(path, fname))
		if err != nil {
			continue
		}
		n.entries = append(n.entries, lstat)
	}
}
