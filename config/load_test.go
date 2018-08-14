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
})
