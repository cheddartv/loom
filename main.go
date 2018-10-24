package main

import (
	"flag"
	"log"
	"os"
	"os/signal"
	"strings"
	"sync"
	"syscall"

	"github.com/Wang/pid"
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

func ParsePidFile(cfg *Config) string {
	if cfg.PidFile == "" {
		return "/var/run/loom.pid"
	} else {
		return cfg.PidFile
	}
}

func Weave(inputs []string, output string, stop chan bool) {
	AllData := []ParsedInput{}
	evts := CreateWatcher(inputs)
	output = CleanPath(output)
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

func SignalSafeMain(osStop chan bool, context Context) {
	var wg sync.WaitGroup

	inputs, outputs := ParseInputsOutput(context.Config)
	stopChannel := make(chan bool)
	log.Print(inputs)
	log.Print("we are spinning up worker groups")
	for i, out := range outputs {
		wg.Add(1)
		go Weave(inputs[i], out, stopChannel)
		wg.Done()
	}
	wg.Wait()
}

func main() {
	pathPtr := flag.String("configPath", "", "config path")
	flag.Parse()
	log.Print("path has been set to ", *pathPtr)
	confPath := strings.Replace(*pathPtr, "loom.yml", "", 1)
	var context Context
	context.Config = Load(confPath)
	pidFile := ParsePidFile(context.Config)
	osStop := make(chan os.Signal, 1)
	closing := make(chan bool, 1)
	signal.Notify(osStop, syscall.SIGINT, syscall.SIGTERM)
	_, err := pid.Create(pidFile)
	if err != nil {
		log.Printf("create pid:%s", err.Error())
		os.Exit(1)
	}
	go func() {
		sig := <-osStop
		log.Print("recieved a signal, killing all workers: ", sig)
		os.Remove(pidFile)
		closing <- true
	}()
	SignalSafeMain(closing, context)
	<-closing
}
