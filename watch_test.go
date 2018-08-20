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

var _ = Describe("EventToString", func() {
	It("returns Create for create", func() {
		Expect(main.EventToString(notify.Create)).To(BeEquivalentTo("Create"))
	})
	It("returns Write for Write", func() {
		Expect(main.EventToString(notify.Write)).To(BeEquivalentTo("Write"))
	})
	It("returns Remove for remove", func() {
		Expect(main.EventToString(notify.Remove)).To(BeEquivalentTo("Remove"))
	})
	It("returns emptyString for the rest", func() {
		Expect(main.EventToString(notify.Rename)).To(BeEquivalentTo(""))
	})
})

var _ = Describe("CreateWatcher", func() {
	paths := []string{"example/primary.m3u8", "example/backup.m3u8"}
	out := main.CreateWatcher(paths)

	It("Should write to out", func() {
		Eventually(out).Should(Receive(Equal(main.Change{"", "", "EndSetup"})))
	})

})
