package main

type app struct {
	nav *nav

	selectedFiles map[*file]struct{}
}

func (a *app) init() {
	a.nav = &nav{}
	a.nav.init()

	a.selectedFiles = make(map[*file]struct{})
}

func (a *app) toggleSelection(file *file) {
	if a.isSelected(file) {
		delete(a.selectedFiles, file)
	} else {
		a.selectedFiles[file] = struct{}{}
	}
}

func (a *app) isSelected(file *file) bool {
	_, ok := a.selectedFiles[file]
	return ok
}
