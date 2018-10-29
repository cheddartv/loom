package main

import (
	"log"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/rjeczalik/notify"
)

type Change struct {
	Path    string
	AbsPath string
	Type    string
}

func LongestExistingPath(dir string) string {
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		splitPath := strings.Split(dir, "/")
		trimmedPath := strings.Join(splitPath[0:len(splitPath)-1], "/")
		return LongestExistingPath(trimmedPath)
	}
	if dir == "" {
		return "/"
	} else {
		return dir
	}
}

func CleanPath(dirty string) string {
	return CleanPaths(dirty, "")
}

func CleanPaths(dirty string, suffix string) string {
	a, _ := filepath.Abs(dirty)
	d, err := filepath.EvalSymlinks(path.Dir(a))
	if err != nil {
		splitPath := strings.Split(dirty, "/")
		trimmedPath := strings.Join(splitPath[0:len(splitPath)-2], "/") + "/" + splitPath[len(splitPath)-1]
		suffix = splitPath[len(splitPath)-2] + "/" + suffix
		return CleanPaths(trimmedPath, suffix)
	}
	d = d + "/" + suffix
	return path.Join(d, path.Base(a))
}

func EventToString(event notify.Event) string {
	switch event {
	case notify.Create:
		return "Create"
	case notify.Write:
		return "Write"
	case notify.Remove:
		return "Remove"
	default:
		return ""
	}

}

func InitWatchers(paths []string, in chan notify.EventInfo) map[string]string {
	pathsMap := make(map[string]string)
	for _, p := range paths {
		pathsMap[CleanPath(p)] = p
		d := LongestExistingPath(path.Dir(p))
		if err := notify.Watch(d+"/...", in, notify.Create, notify.Remove, notify.Write); err != nil {
			log.Fatal(err)
		}
	}
	return pathsMap
}

func ProcessFsEvent(event notify.EventInfo, path string) Change {
	return Change{Path: path, AbsPath: event.Path(), Type: EventToString(event.Event())}
}

func CreateWatcher(paths []string) <-chan Change {
	in := make(chan notify.EventInfo, 10)
	out := make(chan Change)
	var pathsMap map[string]string

	go func() {
		pathsMap = InitWatchers(paths, in)
		defer notify.Stop(in)
		for {
			ei := <-in
			if p, ok := pathsMap[ei.Path()]; ok {
				out <- ProcessFsEvent(ei, p)
			}
		}
	}()

	return out
}
