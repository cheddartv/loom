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
	Path     string
	AbsPath  string
	Remove   bool
	Playlist *m3u8.MasterPlaylist
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
		log.Printf("parsing error\n")
		return nil, err
	}
	if listType != m3u8.MASTER {
		log.Print("m3u8 was not a master playlist, can not weave")
		return nil, errors.New("Incorrect playlist format")
	} else {
		masterpl := p.(*m3u8.MasterPlaylist)
		return masterpl, nil
	}
}

func ImportInputs(dirtyPaths []string) ([]string, []Change) {
	inputStructs := make([]Change, len(dirtyPaths))
	cleanPaths := make([]string, len(dirtyPaths))
	for i, file := range dirtyPaths {
		cleanfile := CleanPath(file)
		cleanPaths[i] = cleanfile
		mp, err := ImportPlaylist(cleanfile)
		if err != nil {
			log.Print("!!ERROR: ", err)
			inputStructs[i] = Change{Path: file, AbsPath: cleanfile, Remove: true, Playlist: nil}
		} else {
			inputStructs[i] = Change{Path: file, AbsPath: cleanfile, Remove: false, Playlist: mp}
		}
	}

	return cleanPaths, inputStructs
}

func FindStructIndexByPath(path string, cs []Change) int {
	for i, c := range cs {
		if c.AbsPath == path {
			return i
		}
	}
	return -1
}

func HandleEvent(event notify.EventInfo, cs []Change) []Change {
	p := event.Path()
	i := FindStructIndexByPath(p, cs)
	if i < 0 {
		return cs
	} else {
		if event.Event() == notify.Remove {
			cs[i].Remove = true
		} else {
			mp, err := ImportPlaylist(p)
			if err != nil {
				log.Print("an update created an error with the playlist", err)
				log.Print("removing the playlist from the viable set")
				cs[i].Remove = true
			} else {
				cs[i].Remove = false
				cs[i].Playlist = mp
			}
		}
	}
	return cs
}

func CreateWatcher(paths []string) <-chan []Change {
	in := make(chan notify.EventInfo, 10)
	out := make(chan []Change)
	paths, allData := ImportInputs(paths)

	go func() {
		for _, p := range paths {
			d := path.Dir(p)
			if err := notify.Watch(d+"/...", in, notify.Create, notify.Remove, notify.Write); err != nil {
				log.Fatal(err)
			}
		}
		defer notify.Stop(in)
		out <- allData
		for {
			ei := <-in
			log.Println("Raw event: ", ei)
			allData = HandleEvent(ei, allData)
			out <- allData
		}
	}()

	return out
}
