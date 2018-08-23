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
	output := context.Config.Manifests[0].Output
	log.Print("output is: ", output)

	AllData := []ParsedInput{}
	evts := CreateWatcher(inputs)
	setupComplete := false
	for {
		evt := <-evts
		log.Println("Got an event:", evt)

		if setupComplete {
			AllData = HandleEvent(evt, AllData)

			WriteManifest(AllData, output)
		} else {
			if evt.Type == "EndSetup" {
				setupComplete = true
				WriteManifest(AllData, output)
			} else {
				AllData = HandleEvent(evt, AllData)
			}
		}

	}
}
