package main

import "strings"

func newHiddenFileNames(fileNames ...string) map[string]struct{} {
	hasSet := make(map[string]struct{}, len(fileNames))
	for _, f := range fileNames {
		hasSet[f] = struct{}{}
	}
	return hasSet
}

func newHiddenPaths(paths ...string) map[string]struct{} {
	hasSet := make(map[string]struct{}, len(paths))
	for _, p := range paths {
		if strings.HasPrefix(p, "~/") {
			p = strings.Replace(p, "~/", homeDir, 1)
		}
		hasSet[p] = struct{}{}
	}
	return hasSet
}
