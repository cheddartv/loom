package main

import (
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"
)

type Context struct {
	Config *Config
}

func ParseInputsOutput(cfg *Config) ([][]string, []string) {
	inputs := [][]string{}
	outputs := []string{}
	for _, m := range cfg.Manifests {
		inputs = append(inputs, m.Inputs)
		outputs = append(outputs, m.Output)
	}
	return inputs, outputs
}

func Weave(inputs []string, output string, stop chan bool) {
	AllData := []ParsedInput{}
	evts := CreateWatcher(inputs)
	for _, p := range inputs {
		AllData = HandleEvent(Change{Path: p, AbsPath: CleanPath(p), Type: "Create"}, AllData)
	}
	WriteManifest(AllData, output)
	for {
		select {
		case evt := <-evts:
			AllData = HandleEvent(evt, AllData)
			WriteManifest(AllData, output)
		case <-stop:
			log.Print("Stopping our weave")
			return
		}
	}
}

func SignalSafeMain(osStop chan bool) {
	var wg sync.WaitGroup
	var context Context
	context.Config = Load()

	inputs, outputs := ParseInputsOutput(context.Config)
	workers := len(outputs)
	stopChannel := make(chan bool)

	for i := 0; i < workers; i++ {
		wg.Add(1)
		go Weave(inputs[i], outputs[i], stopChannel)
		wg.Done()
	}
	wg.Wait()
}

func main() {
	osStop := make(chan os.Signal, 1)
	closing := make(chan bool, 1)
	signal.Notify(osStop, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		sig := <-osStop
		log.Print("recieved a signal, killing all workers: ", sig)
		closing <- true
	}()
	SignalSafeMain(closing)
	<-closing
}
