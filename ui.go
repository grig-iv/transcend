package main

import (
	"unicode/utf8"

	"github.com/gdamore/tcell/v2"
)

type ui struct {
	screen     tcell.Screen
	screenSize size
	scrollPos  int
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

	files := app.visibleFiles()
	cursorFile := app.nav.cursorFile()

	ui.updateScrollPos(files, cursorFile)

	row := headerHeigth
	for _, f := range files[ui.scrollPos:] {
		ui.renderFile(app, f, row, f == cursorFile)
		row += 1
	}

	ui.screen.Show()
}

func (ui *ui) renderHeader(app *app) {
	for i, r := range app.nav.currDir.path {
		ui.screen.SetContent(i, 0, r, nil, headerStyle)
	}
}

func (ui *ui) renderFile(app *app, file *file, row int, isCursorRow bool) {
	style := fileStyle
	if isCursorRow {
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

	icon := getIcon(file)
	iconStyle := style
	if icon.fgColorDark != "" && isCursorRow == false {
		iconStyle = style.Foreground(tcell.GetColor(icon.fgColorDark))
	}
	r, _ := utf8.DecodeRuneInString(icon.text)
	ui.screen.SetContent(1, row, r, nil, iconStyle)

	for c, r := range file.Name() {
		ui.screen.SetContent(c+3, row, r, nil, style)
	}
}

func (ui *ui) updateScrollPos(files []*file, cursorFile *file) {
	cursorFilePos := 0
	for i, f := range files {
		if f == cursorFile {
			cursorFilePos = i
			break
		}
	}

	// top scroll handling
	if ui.scrollPos+scrolloff > cursorFilePos {
		ui.scrollPos = cursorFilePos - scrolloff
		if ui.scrollPos < 0 {
			ui.scrollPos = 0
		}
	}

	// bottom scroll handling
	viewportHeight := ui.screenSize.heigth - 1 - 1
	if ui.scrollPos+viewportHeight-3 < cursorFilePos {
		ui.scrollPos = cursorFilePos - viewportHeight + 3
	}
}

func getIcon(file *file) icon {
	if icon, ok := extToIcon[file.ext()]; ok {
		return icon
	}

	if file.IsDir() {
		return fallbackDirIcon
	}

	return fallbackFileIcon
}
