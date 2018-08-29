package main

import (
	"io/ioutil"
	"log"
	"path/filepath"
	"sort"

	"github.com/grafov/m3u8"
)

func FindStructIndexByPath(abspath string, cs []ParsedInput) int {
	for i, c := range cs {
		if c.AbsPath == abspath {
			return i
		}
	}
	return -1
}

func AddPlaylist(event Change, cs []ParsedInput) []ParsedInput {
	mp, err := ImportPlaylist(event.AbsPath)
	if err != nil {
		log.Printf("an error occured parsing the playlist %v ; It has not be tracked", event.Path)
	} else {
		cs = append(cs, ParsedInput{Path: event.Path, AbsPath: event.AbsPath, Include: true, Playlist: mp})
	}
	return cs
}

func AdjustExistingPlaylist(event Change, cs []ParsedInput, index int) []ParsedInput {
	if event.Type == "Remove" {
		cs[index].Include = false
	} else {
		mp, err := ImportPlaylist(event.AbsPath)
		if err != nil {
			log.Print("an update created an error with the playlist", err)
			log.Print("removing the playlist from the viable set")
			cs[index].Include = false
		} else {
			cs[index].Include = true
			cs[index].Playlist = mp
		}
	}
	return cs
}

func HandleEvent(event Change, cs []ParsedInput) []ParsedInput {
	i := FindStructIndexByPath(event.AbsPath, cs)
	if i < 0 {
		if event.Type == "Create" {
			return AddPlaylist(event, cs)
		} else {
			return cs
		}
	} else {
		return AdjustExistingPlaylist(event, cs, i)
	}
}

type byBandwidth []variantWithPosition

func (s byBandwidth) Len() int {
	return len(s)
}
func (s byBandwidth) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}
func (s byBandwidth) Less(i, j int) bool {
	if s[i].variant.VariantParams.Bandwidth == s[j].variant.VariantParams.Bandwidth {
		if s[i].inputPosition == s[j].inputPosition {
			return s[i].manifestPosition < s[j].manifestPosition
		} else {
			return s[i].inputPosition < s[j].inputPosition
		}
	} else {
		return s[i].variant.VariantParams.Bandwidth > s[j].variant.VariantParams.Bandwidth
	}
}

type variantWithPosition struct {
	inputPosition    int
	manifestPosition int
	variant          *m3u8.Variant
}

func WriteManifest(manifests []ParsedInput, output string) {
	variants := []variantWithPosition{}
	variantToPath := make(map[*m3u8.Variant]string)
	for i, input := range manifests {
		if input.Include {
			for j, v := range input.Playlist.Variants {
				variantToPath[v] = input.AbsPath
				variants = append(variants, variantWithPosition{inputPosition: i, manifestPosition: j, variant: v})
			}
		}
	}
	sort.Sort(byBandwidth(variants))

	outputManifest := m3u8.NewMasterPlaylist()
	for _, v := range variants {
		rel, _ := filepath.Rel(filepath.Dir(output), filepath.Dir(variantToPath[v.variant]))
		outputManifest.Append(rel+"/"+v.variant.URI, v.variant.Chunklist, v.variant.VariantParams)
	}

	d1 := []byte(outputManifest.Encode().String())
	_ = ioutil.WriteFile(output, d1, 0644)

}
