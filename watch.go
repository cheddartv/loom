package main

import (
	"log"
	"path"
	"path/filepath"

	"github.com/rjeczalik/notify"
)

type Change struct {
	Path    string
	AbsPath string
	Type    string
}

func CleanPath(dirty string) string {
	a, _ := filepath.Abs(dirty)

	d, err := filepath.EvalSymlinks(path.Dir(a))
	if err != nil {
		panic("Unable to evaluate symlinks on " + path.Dir(a))
	}
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
		d := path.Dir(p)
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
