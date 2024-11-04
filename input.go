package main

import (
	"strings"

	"github.com/gdamore/tcell/v2"
)

type input struct {
	combo string
}

func (i *input) handle(event *tcell.EventKey, app *app) {
	eventName := ""
	if event.Key() == tcell.KeyRune {
		eventName = string(event.Rune())
	} else {
		eventName = strings.ToLower(event.Name())
	}

	handler, ok := keybindings[eventName]
	if ok {
		handler(app)
		return
	} else {
		if i.combo != "" {
			i.combo += "," + eventName
		} else {
			i.combo = eventName
		}

		handler, ok := keybindings[i.combo]
		if ok {
			i.combo = ""
			handler(app)
			return
		}

		for k := range keybindings {
			if strings.HasPrefix(k, i.combo) {
				return
			}
		}

		i.combo = ""
	}
}
