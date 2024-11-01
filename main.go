package main

import (
	"log"
	"os"

	"github.com/gdamore/tcell/v2"
)

func main() {
	f, err := os.OpenFile("log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("error opening file: %v", err)
	}
	defer f.Close()
	f.Truncate(0)

	log.SetOutput(f)

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
				app.nav.cursorNext()
			case tcell.KeyUp:
				app.nav.cursorPrev()
			case tcell.KeyLeft:
				app.nav.upDir()
			case tcell.KeyRight:
				app.nav.intoDir()
			}
		case *tcell.EventResize:
			ui.onResize()
		}

		ui.render(app)
	}
}
