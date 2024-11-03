package main

type app struct {
	nav *nav

	selectedFiles map[string]struct{}
}

func (a *app) init() {
	a.nav = &nav{}
	a.nav.init()

	a.selectedFiles = make(map[string]struct{})
}

func (a *app) toggleSelection(file *file) {
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
