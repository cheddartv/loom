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

	AllData := make([]ParsedInput, 0)
	evts := CreateWatcher(inputs)
	setupComplete := false
	for {
		evt := <-evts
		log.Println("Got an event:", evt)

		if setupComplete {
			AllData = HandleEvent(evt, AllData)
			// log.Println("Alldata is: ", AllData)
			WriteManifest(AllData)
		} else {
			if evt.Type == "EndSetup" {
				setupComplete = true
				WriteManifest(AllData)
			} else {
				AllData = HandleEvent(evt, AllData)
			}
		}

	}
}
