package main

import (
	"log"
	"os"
	"path"
)

type app struct {
	nav *nav

	showHidden    bool
	selectedFiles map[string]struct{}
}

var (
	homeDir = os.Getenv("HOME") + "/"
)

func (a *app) init() {
	a.nav = &nav{}
	a.nav.init()

	a.setCursorOnVisibleFile()

	a.selectedFiles = make(map[string]struct{})
}

func (a *app) cursorPrev() { a.nav.cursorPrev(!a.showHidden) }
func (a *app) cursorNext() { a.nav.cursorNext(!a.showHidden) }

func (a *app) upDir() {
	a.nav.upDir()
	a.setCursorOnVisibleFile()
}

func (a *app) intoDir() {
	a.nav.intoDir()
	a.setCursorOnVisibleFile()
}

func (a *app) toggleSelection() {
	file := a.nav.cursorFile()
	if a.isSelected(file) {
		delete(a.selectedFiles, file.path)
	} else {
		a.selectedFiles[file.path] = struct{}{}
		a.nav.cursorNext(!a.showHidden)
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

func (a *app) toggleHidden() {
	a.showHidden = !a.showHidden
	a.setCursorOnVisibleFile()
}

func (a *app) setCursorOnVisibleFile() {
	a.nav.cursorNext(!a.showHidden)
	a.nav.cursorPrev(!a.showHidden)
}
