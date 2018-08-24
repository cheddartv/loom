package main_test

import (
	main "github.com/cheddartv/loom"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("loading configuration file", func() {
	var c *main.Config

	BeforeEach(func() {
		c = main.Load()
	})

	It("creates the correct number of manifests", func() {
		Expect(len(c.Manifests)).To(Equal(2))
	})

	It("configures the correct output", func() {
		Expect(c.Manifests[0].Output).To(Equal("example/index.m3u8"))
	})

	It("configures the correct inputs", func() {
		Expect(c.Manifests[0].Inputs).To(Equal([]string{"example/primary/primary.m3u8", "example/backup/backup.m3u8"}))
	})
})
