package main

import (
	"github.com/gdamore/tcell/v2"
)

type ui struct {
	screen     tcell.Screen
	screenSize size
}

type size struct {
	width, heigth int
}

var (
	headerStyle     = tcell.StyleDefault.Background(tcell.ColorReset).Foreground(tcell.ColorGreen)
	fileStyle       = tcell.StyleDefault.Background(tcell.ColorReset).Foreground(tcell.ColorWhite)
	cursorFileStyle = tcell.StyleDefault.Background(tcell.ColorWhite).Foreground(tcell.ColorBlack)
)

func (ui *ui) init(screen tcell.Screen) {
	ui.screen = screen
}

func (ui *ui) onResize() {
	width, height := ui.screen.Size()
	ui.screenSize = size{width, height}
}

func (ui *ui) render(app *app) {
	ui.screen.Clear()

	ui.renderHeader(app)

	for i, f := range app.nav.files {
		ui.renderFile(app, f, i)
	}

	ui.screen.Show()
}

func (ui *ui) renderHeader(app *app) {
	for i, r := range app.nav.currDir.path {
		ui.screen.SetContent(i, 0, r, nil, headerStyle)
	}
}

func (ui *ui) renderFile(app *app, file *file, index int) {
	for c, r := range file.Name() {
		if index != app.nav.cursor {
			ui.screen.SetContent(c, index+1, r, nil, fileStyle)
		} else {
			ui.screen.SetContent(c, index+1, r, nil, cursorFileStyle)
		}
	}
}
