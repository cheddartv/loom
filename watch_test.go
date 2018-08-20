package main_test

import (
	"math/rand"
	"os"
	"path"
	"path/filepath"

	main "github.com/cheddartv/loom"
	"github.com/rjeczalik/notify"

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
		_, structs := main.ImportInputs(paths)
		Expect(structs[0].Playlist).To(BeNil())
	})

})

var _ = Describe("FindStructIndexByPath", func() {
	input := []main.Change{
		{Path: "primary.m3u8", AbsPath: "/root/primary.m3u8", Remove: false, Playlist: nil},
		{Path: "primary2.m3u8", AbsPath: "/root/primary2.m3u8", Remove: false, Playlist: nil},
	}

	It("returns the index of the struct", func() {
		Expect(main.FindStructIndexByPath("/root/primary.m3u8", input)).To(Equal(0))
	})

	It("returns -1 if the file is not present", func() {
		Expect(main.FindStructIndexByPath("/not/present/primary.m3u8", input)).To(Equal(-1))
	})

})

type mockEventInfo struct {
	Type     notify.Event
	BasePath string
}

func (e mockEventInfo) Event() notify.Event {
	return e.Type
}

func (e mockEventInfo) Path() string {
	workingDir, _ := filepath.EvalSymlinks(os.Getenv("PWD"))
	return workingDir + e.BasePath
}

func (e mockEventInfo) Sys() interface{} {
	return nil
}

var _ = Describe("HandleEvent", func() {
	workingDir, _ := filepath.EvalSymlinks(os.Getenv("PWD"))
	input := []main.Change{
		{Path: "primary.m3u8", AbsPath: workingDir + "/example/primary.m3u8", Remove: false, Playlist: nil},
		{Path: "backup.m3u8", AbsPath: workingDir + "/example/backup.m3u8", Remove: true, Playlist: nil},
		{Path: "1.m3u8", AbsPath: workingDir + "/example/1.m3u8", Remove: false, Playlist: nil},
	}

	It("Updates the data struct on file creation", func() {
		event := mockEventInfo{Type: notify.Create, BasePath: "/example/backup.m3u8"}

		Expect(main.HandleEvent(event, input)[1].Remove).To(BeFalse())
	})

	It("Updates the data struct on file removal", func() {
		event := mockEventInfo{Type: notify.Remove, BasePath: "/example/primary.m3u8"}

		Expect(main.HandleEvent(event, input)[0].Remove).To(BeTrue())
	})

	It("non-tracked files do not change data struct", func() {
		event := mockEventInfo{Type: notify.Remove, BasePath: "/example/not_tracked.m3u8"}

		Expect(main.HandleEvent(event, input)).Should(BeEquivalentTo(input))
	})

	It("Removes a playlist if an update makes it fail parsing", func() {
		event := mockEventInfo{Type: notify.Write, BasePath: "/example/1.m3u8"}

		Expect(main.HandleEvent(event, input)[2].Remove).To(BeTrue())
	})

})

var _ = Describe("CreateWatcher", func() {
	paths := []string{"example/primary.m3u8", "example/backup.m3u8"}
	out := main.CreateWatcher(paths)

	It("Should write to out", func() {
		Eventually(out).Should(Receive())
	})

})
