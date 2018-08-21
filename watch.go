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

func CreateWatcher(paths []string) <-chan Change {
	in := make(chan notify.EventInfo, 10)
	out := make(chan Change)
	pathsMap := make(map[string]string)

	go func() {
		for _, p := range paths {
			pathsMap[CleanPath(p)] = p
			out <- Change{Path: p, AbsPath: CleanPath(p), Type: EventToString(notify.Create)}
			d := path.Dir(p)
			if err := notify.Watch(d+"/...", in, notify.Create, notify.Remove, notify.Write); err != nil {
				log.Fatal(err)
			}
		}
		out <- Change{Path: "", AbsPath: "", Type: "EndSetup"}
		defer notify.Stop(in)
		for {
			ei := <-in
			log.Println("Raw event: ", ei)
			if p, ok := pathsMap[ei.Path()]; ok {
				out <- Change{Path: p, AbsPath: ei.Path(), Type: EventToString(ei.Event())}
			}
		}
	}()

	return out
}
