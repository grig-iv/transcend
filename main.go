package main

import (
	"log"
	"os"

	"github.com/gdamore/tcell/v2"
)

func main() {
	logToFile("log")

	s, err := tcell.NewScreen()
	if err != nil {
		log.Fatalf("%+v", err)
	}
	if err := s.Init(); err != nil {
		log.Fatalf("%+v", err)
	}

	quit := func() {
		maybePanic := recover()
		s.Fini()
		if maybePanic != nil {
			panic(maybePanic)
		}
	}
	defer quit()

	app := &app{}
	app.init()

	ui := &ui{}
	ui.init(s)

	for {
		ev := s.PollEvent()
		switch ev := ev.(type) {
		case *tcell.EventKey:
			switch ev.Key() {
			case tcell.KeyCtrlQ:
				return
			case tcell.KeyDown:
				app.cursorNext()
			case tcell.KeyUp:
				app.cursorPrev()
			case tcell.KeyLeft:
				app.upDir()
			case tcell.KeyRight:
				app.intoDir()
			}
			switch ev.Rune() {
			case ' ':
				app.toggleSelection()
			}
		case *tcell.EventResize:
			ui.onResize()
		}

		ui.render(app)
	}
}

func logToFile(path string) func() {
	f, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_APPEND|os.O_TRUNC, 0666)
	if err != nil {
		log.Fatalf("error opening file: %v", err)
	}

	log.SetOutput(f)

	return func() { f.Close() }
}
