package main_test

import (
	"math/rand"
	"os"
	"path"
	"path/filepath"

	main "github.com/cheddartv/loom"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

func name() string {
	b := make([]rune, 6)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}

var _ = Describe("generating clean paths", func() {
	workingDir, _ := filepath.EvalSymlinks(os.Getenv("PWD"))

	It("panics if provided path does not exist", func() {
		Expect(func() {
			main.CleanPath(path.Join("tmp", name(), "index.m3u8"))
		}).To(Panic())
	})

	It("absolutifies and unsymlinks the path(s) it's given", func() {
		Expect(main.CleanPath("tmp/index.m3u8")).To(Equal(workingDir + "/tmp/index.m3u8"))
	})
})

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
		playlist, err := main.ImportPlaylist("./example/primary.m3u8")
		Expect(playlist.Variants[0].URI).To(Equal("primary/1.m3u8"))
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
		_, structMap := main.ImportInputs(paths)
		Expect(structMap[paths[0]]).To(BeNil())
	})

})
