package main_test

import (
	"time"

	main "github.com/cheddartv/loom"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("ParsedInputsOutput", func() {
	var context main.Context
	context.Config = main.Load()

	It("parses the yml file", func() {
		inputs, _ := main.ParseInputsOutput(context.Config)
		Expect(len(inputs)).To(BeEquivalentTo(2))
	})
})

var _ = Describe("ParsePidFile", func() {
	var context main.Context
	context.Config = main.Load()

	It("parses the yml file", func() {
		pidFile := main.ParsePidFile(context.Config)
		Expect(pidFile).To(BeEquivalentTo("/var/run/loom.pid"))
	})

})

var _ = Describe("Weave", func() {
	var context main.Context
	context.Config = main.Load()
	inputs, _ := main.ParseInputsOutput(context.Config)
	stop := make(chan bool)
	It("Weave to spin on a channel", func() {
		go main.Weave(inputs[0], "./tmp/weaveindex.m3u8", stop)
		time.Sleep(1000)
		stop <- true
		Expect("./tmp/weaveindex.m3u8").Should(BeAnExistingFile())
	})

})

var _ = Describe("Main exits", func() {
	var context main.Context
	context.Config = main.Load()
	It("Should eventually exit", func() {
		stop := make(chan bool, 1)
		stop <- true
		Eventually(func() bool {
			main.SignalSafeMain(stop, context)
			return true
		}).Should(BeTrue())
	})
})
