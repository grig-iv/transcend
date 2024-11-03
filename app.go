package main

import (
	"log"
	"path"
)

type app struct {
	nav *nav

	selectedFiles map[string]struct{}
}

func (a *app) init() {
	a.nav = &nav{}
	a.nav.init()

	a.selectedFiles = make(map[string]struct{})
}

func (a *app) cursorPrev() { a.nav.cursorPrev() }
func (a *app) cursorNext() { a.nav.cursorNext() }
func (a *app) upDir()      { a.nav.upDir() }
func (a *app) intoDir()    { a.nav.intoDir() }

func (a *app) toggleSelection() {
	file := a.nav.cursorFile()
	if a.isSelected(file) {
		delete(a.selectedFiles, file.path)
	} else {
		a.selectedFiles[file.path] = struct{}{}
		a.nav.cursorNext()
	}
}

func (a *app) isSelected(file *file) bool {
	_, ok := a.selectedFiles[file.path]
	return ok
}

func (a *app) copySelected() {
	for sourcePath := range a.selectedFiles {
		sourceName := path.Base(sourcePath)
		destPath := path.Join(a.nav.currDir.path, sourceName)
		err := copyFile(sourcePath, destPath)
		if err != nil {
			log.Println(err)
		}
	}
}
