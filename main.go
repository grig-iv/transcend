package main

import (
	"log"
	"os"

	"github.com/gdamore/tcell/v2"
)

func main() {
	logToFile("log")

	screen, screenEventCh := initScreen()

	quit := func() {
		maybePanic := recover()
		screen.Fini()
		if maybePanic != nil {
			panic(maybePanic)
		}
	}
	defer quit()

	app := &app{}
	app.init()

	ui := &ui{}
	ui.init(screen)

	for {
		select {
		case ev := <-screenEventCh:
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
				case 'c':
					app.copySelected()
				case 'd':
					app.deleteSelected()
				case '.':
					app.toggleHidden()
				}

				switch ev.Name() {
				case "Ctrl+PgDn":
					app.cursorLast()
				case "Ctrl+PgUp":
					app.cursorFirst()
				}

			case *tcell.EventResize:
				ui.onResize()
			}
			ui.render(app)

		case watchEv := <-app.watch.events:
			redraw := app.onWatchEvent(watchEv)
			if redraw {
				ui.render(app)
			}
		}

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

func initScreen() (tcell.Screen, <-chan tcell.Event) {
	s, err := tcell.NewScreen()
	if err != nil {
		log.Fatalf("%+v", err)
	}
	if err := s.Init(); err != nil {
		log.Fatalf("%+v", err)
	}

	eventCh := make(chan tcell.Event)

	go func() {
		for {
			eventCh <- s.PollEvent()
		}
	}()

	return s, eventCh
}
