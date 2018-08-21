package main

<<<<<<< HEAD
import "log"
=======
import (
	"log"
)
>>>>>>> reader and writter

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
<<<<<<< HEAD

	evts := CreateWatcher(inputs)
	for {
		evt := <-evts
		log.Println("Got an event:", evt)
=======
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

>>>>>>> reader and writter
	}
}
