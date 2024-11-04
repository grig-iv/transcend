package main

import (
	"log"
	"time"

	"github.com/fsnotify/fsnotify"
)

type watch struct {
	watcher *fsnotify.Watcher

	watchigPaths map[string]struct{}
	syncTimer    *time.Timer

	events chan fsnotify.Event
}

func newWatch() *watch {
	var err error

	w := &watch{}

	w.watcher, err = fsnotify.NewWatcher()
	if err != nil {
		log.Fatal(err)
	}

	w.watchigPaths = map[string]struct{}{}
	w.syncTimer = time.NewTimer(0)

	w.events = make(chan fsnotify.Event)

	go w.loop()

	return w
}

func (w *watch) loop() {
	for {
		select {
		case event, ok := <-w.watcher.Events:
			if !ok {
				return
			}
			w.events <- event
		case err, ok := <-w.watcher.Errors:
			if !ok {
				return
			}
			log.Println("watch error", err)
		case <-w.syncTimer.C:
			w.sync()
		}
	}
}

func (w *watch) add(path string) {
	w.watchigPaths[path] = struct{}{}
	w.syncTimer.Reset(1 * time.Millisecond)
}

func (w *watch) remove(path string) {
	delete(w.watchigPaths, path)
	w.syncTimer.Reset(1 * time.Millisecond)
}

func (w *watch) sync() {
	for _, p := range w.watcher.WatchList() {
		_, shouldWatch := w.watchigPaths[p]
		if shouldWatch == false {
			w.watcher.Remove(p)
		}
	}

	for p := range w.watchigPaths {
		w.watcher.Add(p)
	}
}

func (w *watch) close() {
	w.watcher.Close()
}
