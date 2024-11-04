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

	input := input{}

	for {
		select {
		case <-app.quitChan:
			app.watch.close()
			return

		case ev := <-screenEventCh:
			switch ev := ev.(type) {
			case *tcell.EventKey:
				input.handle(ev, app)
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
