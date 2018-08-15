package main

import (
	"bufio"
	"errors"
	"log"
	"os"
	"path"
	"path/filepath"

	"github.com/grafov/m3u8"
	"github.com/rjeczalik/notify"
)

type Change struct {
	Path   string
	Remove bool
}

func CleanPath(dirty string) string {
	a, _ := filepath.Abs(dirty)

	d, err := filepath.EvalSymlinks(path.Dir(a))
	if err != nil {
		panic("Unable to evaluate symlinks on " + path.Dir(a))
	}
	return path.Join(d, path.Base(a))
}

func ImportPlaylist(file string) (*m3u8.MasterPlaylist, error) {
	f, err := os.Open(file)
	if err != nil {
		log.Printf("Provided file %s errored and was skipped", file)
		return nil, err
	}
	p, listType, err := m3u8.DecodeFrom(bufio.NewReader(f), true)
	if err != nil {
		return nil, err
	}
	if listType != m3u8.MASTER {
		log.Print("m3u8 was not a master playlist, can not weave")
		return nil, errors.New("Incorrect playlist format")
	} else {
		masterpl := p.(*m3u8.MasterPlaylist)
		log.Printf("%+v\n", masterpl)
		return masterpl, nil
	}
}

func ImportInputs(dirtyPaths []string) ([]string, map[string]*m3u8.MasterPlaylist) {
	inputStructs := make(map[string]*m3u8.MasterPlaylist)
	cleanPaths := make([]string, len(dirtyPaths))
	for _, file := range dirtyPaths {
		cleanfile := CleanPath(file)
		cleanPaths = append(cleanPaths, cleanfile)
		mp, err := ImportPlaylist(cleanfile)
		if err != nil {
			log.Print(err)
		} else {
			inputStructs[file] = mp
		}
	}
	return cleanPaths, inputStructs
}

func CreateWatcher(paths []string) <-chan Change {
	in := make(chan notify.EventInfo, 10)
	out := make(chan Change)
	paths, _ = ImportInputs(paths)

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
