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

var _ = Describe("longest existing path", func() {
	real_path, _ := filepath.EvalSymlinks(os.Getenv("PWD"))
	long_path := real_path + "/along/the/fake/path/we/go.tar"
	It("returns the longest real file path", func() {
		Expect(main.LongestExistingPath(long_path)).To(Equal(real_path))
	})
})

var _ = Describe("generating clean paths", func() {
	workingDir, _ := filepath.EvalSymlinks(os.Getenv("PWD"))

	It("panics if provided path does not exist", func() {
		Expect(func() {
			main.CleanPath(path.Join("tmp", name(), "index.m3u8"))
		}).ToNot(Panic())
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

var _ = Describe("ProcessFsEvent", func() {
	ev := mockEventInfo{Type: notify.Create, BasePath: "primary/sample.m3u8"}
	It("returns a change object", func() {
		Expect(main.ProcessFsEvent(ev, "tmp/last.m3u8").Type).To(BeEquivalentTo("Create"))
	})
	It("returns a change object", func() {
		Expect(main.ProcessFsEvent(ev, "tmp/last.m3u8").Path).To(BeEquivalentTo("tmp/last.m3u8"))
	})
	It("returns a change object", func() {
		Expect(main.ProcessFsEvent(ev, "tmp/last.m3u8").AbsPath).To(BeEquivalentTo(ev.Path()))
	})
})
