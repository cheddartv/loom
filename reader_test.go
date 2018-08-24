package main_test

import (
	"os"
	"path/filepath"

	main "github.com/cheddartv/loom"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("import playlist from file", func() {
	It("return nil playlist and error when playlist doesn't exist", func() {
		playlist, err := main.ImportPlaylist("playlist.m3u8")
		Expect(playlist).To(BeNil())
		Expect(err).Should(HaveOccurred())
	})
	It("return nil playlist and error when its not a playlist", func() {
		playlist, err := main.ImportPlaylist("./watch_test.go")
		Expect(playlist).To(BeNil())
		Expect(err).Should(HaveOccurred())
	})
	It("return nil playlist and error, when its not a master playlist", func() {
		playlist, err := main.ImportPlaylist("./example/1.m3u8")
		Expect(playlist).To(BeNil())
		Expect(err).Should(HaveOccurred())
	})
	It("returns a struct of the masterplaylist", func() {
		playlist, err := main.ImportPlaylist("./example/primary/primary.m3u8")
		Expect(playlist.Variants[0].URI).To(Equal("1.m3u8"))
		Expect(err).ShouldNot(HaveOccurred())
	})
})
var _ = Describe("importInputs playlists from paths", func() {
	workingDir, _ := filepath.EvalSymlinks(os.Getenv("PWD"))
	It("cleans the paths", func() {
		paths := []string{"example/primary.m3u8", "example/backup.m3u8"}
		c, _ := main.ImportInputs(paths)
		Expect(c[0]).To(Equal(workingDir + "/" + paths[0]))
		Expect(c[1]).To(Equal(workingDir + "/" + paths[1]))
	})
	It("ImportPlaylist gets a bad path, it doesn't add the file to the map", func() {
		paths := []string{"example/primary_dne.m3u8"}
		_, structs := main.ImportInputs(paths)
		Expect(structs[0].Playlist).To(BeNil())
	})
})
