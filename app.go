package main

import (
	"io"
	"log"
	"os"
	"path"

	"github.com/fsnotify/fsnotify"
)

type app struct {
	nav   *nav
	watch *watch

	showHidden    bool
	selectedPaths map[string]struct{}

	quitChan chan struct{}
}

var (
	homeDir = os.Getenv("HOME") + "/"
)

func (a *app) init() {
	a.nav = &nav{}
	a.nav.init()

	a.watch = newWatch()
	a.watch.add(a.nav.currDir.path)

	a.selectedPaths = make(map[string]struct{})

	a.quitChan = make(chan struct{})

	a.setCursorOnVisibleFile()
}

func (a *app) cursorPrev()  { a.nav.cursorPrev(!a.showHidden) }
func (a *app) cursorNext()  { a.nav.cursorNext(!a.showHidden) }
func (a *app) cursorLast()  { a.nav.cursorLast(!a.showHidden) }
func (a *app) cursorFirst() { a.nav.cursorFirst(!a.showHidden) }

func (a *app) upDir() {
	a.watch.remove(a.nav.currDir.path)
	a.nav.upDir()
	a.watch.add(a.nav.currDir.path)
	a.setCursorOnVisibleFile()
}

func (a *app) intoDir() {
	a.watch.remove(a.nav.currDir.path)
	a.nav.intoDir()
	a.watch.add(a.nav.currDir.path)
	a.setCursorOnVisibleFile()
}

func (a *app) toggleSelection() {
	file := a.nav.cursorFile()
	if a.isSelected(file) {
		delete(a.selectedPaths, file.path)
	} else {
		a.selectedPaths[file.path] = struct{}{}
		a.nav.cursorNext(!a.showHidden)
	}
}

func (a *app) isSelected(file *file) bool {
	_, ok := a.selectedPaths[file.path]
	return ok
}

func (a *app) copySelected() {
	for sourcePath := range a.selectedPaths {
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

func (a *app) visibleFiles() []*file {
	files := make([]*file, 0)
	for _, f := range a.nav.files {
		if f.isHidden() && a.showHidden == false {
			continue
		}

		files = append(files, f)
	}
	return files
}

func (a *app) deleteSelected() {
	for p := range a.selectedPaths {
		err := os.RemoveAll(p)
		if err != nil {
			log.Println("deleteSelected", err)
		}
		delete(a.selectedPaths, p)
	}
}

func copyFile(srcPath, dstPath string) error {
	srcFile, err := os.OpenFile(srcPath, os.O_RDONLY, 0)
	if err != nil {
		return err
	}

	defer srcFile.Close()

	srcStat, err := srcFile.Stat()
	if err != nil {
		return err
	}

	dstFile, err := os.OpenFile(dstPath, os.O_WRONLY|os.O_CREATE|os.O_EXCL, srcStat.Mode())
	if err != nil {
		return err
	}

	defer dstFile.Close()

	_, err = io.Copy(dstFile, srcFile)
	if err != nil {
		return err
	}

	return nil
}

func (a *app) onWatchEvent(event fsnotify.Event) bool {
	if event.Has(fsnotify.Create|fsnotify.Remove) &&
		path.Dir(event.Name) == a.nav.currDir.path {
		a.nav.refreshCurrDir()
		a.setCursorOnVisibleFile()
		return true
	}

	return false
}

func (a *app) quit() {
	go func() {
		a.quitChan <- struct{}{}
	}()
}
