package main

var keybindings = map[string]func(*app){
	"ctrl+q": func(a *app) { a.quit() },

	"left":  func(a *app) { a.upDir() },
	"right": func(a *app) { a.intoDir() },

	"up":        func(a *app) { a.cursorPrev() },
	"down":      func(a *app) { a.cursorNext() },
	"ctrl+pgdn": func(a *app) { a.cursorLast() },
	"ctrl+pgup": func(a *app) { a.cursorFirst() },

	" ": func(a *app) { a.toggleSelection() },
	"c": func(a *app) { a.copySelected() },
	"d": func(a *app) { a.deleteSelected() },

	".": func(a *app) { a.toggleHidden() },
}
