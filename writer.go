package main

import (
	"io/ioutil"
	"log"
	"path/filepath"
	"sort"

	"github.com/grafov/m3u8"
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
				log.Printf("an error occured parsing the playlist %v ; It has not be tracked", event.Path)
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
				cs[i].Include = true
				cs[i].Playlist = mp
			}
		}
	}
	return cs
}

type byBandwidth []*m3u8.Variant

func (s byBandwidth) Len() int {
	return len(s)
}
func (s byBandwidth) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}
func (s byBandwidth) Less(i, j int) bool {
	return s[i].VariantParams.Bandwidth > s[j].VariantParams.Bandwidth
}

func WriteManifest(manifests []ParsedInput, output string) {
	variants := []*m3u8.Variant{}
	for _, input := range manifests {
		if input.Include {
			for _, v := range input.Playlist.Variants {
				rel, _ := filepath.Rel(filepath.Dir(output), filepath.Dir(input.AbsPath))
				v.URI = rel + "/" + v.URI
				variants = append(variants, v)
			}
		}
	}
	sort.Sort(byBandwidth(variants))

	outputManifest := m3u8.NewMasterPlaylist()
	for _, v := range variants {
		outputManifest.Append(v.URI, v.Chunklist, v.VariantParams)
	}

	d1 := []byte(outputManifest.Encode().String())
	_ = ioutil.WriteFile(output, d1, 0644)

}
