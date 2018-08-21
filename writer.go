package main

import (
	"log"

	"github.com/jinzhu/copier"
)

func FindStructIndexByPath(path string, cs []ParsedInput) int {
	for i, c := range cs {
		if c.Path == path {
			return i
		}
	}
	return -1
}

func HandleEvent(event Change, cs []ParsedInput) []ParsedInput {
	i := FindStructIndexByPath(event.Path, cs)
	log.Printf("i: %v, event.path: %v", i, event.Path)
	if i < 0 {
		if event.Type == "Create" {
			mp, err := ImportPlaylist(event.AbsPath)
			if err != nil {

			} else {
				cs = append(cs, ParsedInput{Path: event.Path, AbsPath: event.AbsPath, Include: true, Playlist: mp})
			}
		}
		return cs
	} else {
		if event.Type == "Remove" {
			cs[i].Include = false
		} else {
			mp, err := ImportPlaylist(event.AbsPath)
			if err != nil {
				log.Print("an update created an error with the playlist", err)
				log.Print("removing the playlist from the viable set")
				cs[i].Include = false
			} else {
				// log.Printf("updating %v: ", cs[i])
				cs[i].Include = true
				cs[i].Playlist = mp
			}
		}
	}
	return cs
}

func WriteManifest(manifests []ParsedInput) {

	master := ParsedInput{}
	copier.Copy(&master, &manifests[0])
	playlists := manifests[1:]

	for _, v := range playlists {
		master.Playlist.Variants = append(master.Playlist.Variants, v.Playlist.Variants...)
	}
	log.Print("\nTrue Master\n")
	log.Print(master.Playlist)
	log.Print("\nEnd True Master\n")
}
