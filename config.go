package main

import (
	"cmp"
	"strings"
)

var (
	scrolloff = 3
)

var hiddenFileNames = map[string]struct{}{
	"LICENSE":    {},
	"flake.lock": {},
}

var hiddenPaths = map[string]struct{}{
	"~/go":    {},
	"~/steam": {},
}

func (f *file) isHidden() bool {
	if strings.HasPrefix(f.Name(), ".") {
		return true
	}

	if _, ok := hiddenFileNames[f.Name()]; ok {
		return true
	}

	if _, ok := hiddenPaths[f.path]; ok {
		return true
	}

	return false
}

func sortFunc(lhs, rhs *file) int {
	// dirs goes first
	if lhs.IsDir() != rhs.IsDir() {
		if lhs.IsDir() {
			return -1
		} else {
			return 1
		}
	}

	// then sort by extensions
	extCmp := cmp.Compare(lhs.ext(), rhs.ext())
	if extCmp != 0 {
		return extCmp
	}

	// then by names
	return cmp.Compare(lhs.Name(), rhs.Name())
}
