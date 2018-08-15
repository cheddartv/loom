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
