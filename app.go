package main

type app struct {
	nav *nav
}

func (a *app) init() {
	a.nav = &nav{}
	a.nav.init()
}
