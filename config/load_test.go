package config_test

import (
	"github.com/cheddartv/loom/config"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("loading configuration file", func() {
	var c *config.Config

	BeforeEach(func() {
		c = config.Load()
	})

	It("creates the correct number of manifests", func() {
		Expect(len(c.Manifests)).To(Equal(2))
	})

	It("configures the correct output", func() {
		Expect(c.Manifests[0].Output).To(Equal("a.m3u8"))
	})

	It("configures the correct inputs", func() {
		Expect(c.Manifests[0].Inputs).To(Equal([]string{"a/primary.m3u8", "a/backup.m3u8"}))
	})
})
