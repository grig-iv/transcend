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
	headerStyle        = tcell.StyleDefault.Background(tcell.ColorReset).Foreground(tcell.ColorGreen)
	entryStyle         = tcell.StyleDefault.Background(tcell.ColorReset).Foreground(tcell.ColorWhite)
	selectedEntryStyle = tcell.StyleDefault.Background(tcell.ColorLightSkyBlue).Foreground(tcell.ColorBlack)
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
	ui.renderEntries(app)

	ui.screen.Show()
}

func (ui *ui) renderHeader(app *app) {
	for i, r := range app.nav.currDir.path {
		ui.screen.SetContent(i, 0, r, nil, headerStyle)
	}
}

func (ui *ui) renderEntries(app *app) {
	for i, f := range app.nav.files {
		style := entryStyle
		if i == app.nav.cursor {
			style = selectedEntryStyle
			for c := range ui.screenSize.width {
				ui.screen.SetContent(c, i+1, ' ', nil, style)
			}
		}
		for c, r := range f.Name() {
			ui.screen.SetContent(c, i+1, r, nil, style)
		}
	}
}
