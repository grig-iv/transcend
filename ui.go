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

const (
	headerHeigth = 1
)

var (
	headerStyle        = tcell.StyleDefault.Background(tcell.ColorReset).Foreground(tcell.ColorGreen)
	fileStyle          = tcell.StyleDefault.Background(tcell.ColorReset).Foreground(tcell.ColorWhite)
	cursorFileStyle    = tcell.StyleDefault.Background(tcell.ColorWhite).Foreground(tcell.ColorBlack)
	selectedIndication = tcell.StyleDefault.Background(tcell.ColorYellow)
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
	row := index + headerHeigth

	style := fileStyle
	if index == app.nav.cursor {
		style = cursorFileStyle
	}

	for c := range ui.screenSize.width {
		ui.screen.SetContent(c, row, ' ', nil, style)
	}

	if app.isSelected(file) {
		ui.screen.SetContent(0, row, ' ', nil, selectedIndication)
	} else {
		ui.screen.SetContent(0, row, ' ', nil, style)
	}

	for c, r := range file.Name() {
		ui.screen.SetContent(c+1, row, r, nil, style)
	}
}
