package main

import (
	"bufio"
	"errors"
	"log"
	"os"

	"github.com/grafov/m3u8"
)

type ParsedInput struct {
	Path     string
	AbsPath  string
	Include  bool
	Playlist *m3u8.MasterPlaylist
}

func ImportPlaylist(file string) (*m3u8.MasterPlaylist, error) {
	f, err := os.Open(file)
	defer f.Close()
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

func ImportInputs(dirtyPaths []string) ([]string, []ParsedInput) {
	inputStructs := make([]ParsedInput, len(dirtyPaths))
	cleanPaths := make([]string, len(dirtyPaths))
	for i, file := range dirtyPaths {
		cleanfile := CleanPath(file)
		cleanPaths[i] = cleanfile
		mp, err := ImportPlaylist(cleanfile)
		if err != nil {
			log.Print("!!ERROR: ", err)
			inputStructs[i] = ParsedInput{Path: file, AbsPath: cleanfile, Include: false, Playlist: mp}
		} else {
			inputStructs[i] = ParsedInput{Path: file, AbsPath: cleanfile, Include: true, Playlist: mp}
		}
	}

	return cleanPaths, inputStructs
}
