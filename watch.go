package main

import (
	"log"
	"path"
	"path/filepath"

	"github.com/rjeczalik/notify"
)

type Change struct {
	Path   string
	Remove bool
}

func CleanPaths(dirty []string) []string {
	cleaned := make([]string, len(dirty))
	for i, p := range dirty {
		a, err := filepath.Abs(p)
		if err != nil {
			log.Fatal("Unable to generate absolute path from ", p)
		}

		d, err := filepath.EvalSymlinks(path.Dir(a))
		if err != nil {
			log.Fatal("Unable to evaluate symlinks on ", path.Dir(a))
		}

		cleaned[i] = path.Join(d, path.Base(a))
	}

	return cleaned
}

func CreateWatcher(paths []string) <-chan Change {
	in := make(chan notify.EventInfo, 10)
	out := make(chan Change)
	paths = CleanPaths(paths)

	go func() {
		for _, p := range paths {
			d := path.Dir(p)
			if err := notify.Watch(d+"/...", in, notify.Create, notify.Remove, notify.Write); err != nil {
				log.Fatal(err)
			}
		}
		defer notify.Stop(in)

		for {
			ei := <-in
			log.Println("Raw event: ", ei)
			for _, p := range paths {
				if ei.Path() == p {
					out <- Change{Path: p, Remove: ei.Event() == notify.Remove}
				}
			}
		}
	}()

	return out
}
