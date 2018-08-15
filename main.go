package main

import "log"

type Context struct {
	Config *Config
}

func main() {
	var context Context
	context.Config = Load()

	inputs := []string{}
	for _, m := range context.Config.Manifests {
		inputs = append(inputs, m.Inputs...)
	}

	evts := CreateWatcher(inputs)
	for {
		evt := <-evts
		log.Println("Got an event:", evt)
	}
}
