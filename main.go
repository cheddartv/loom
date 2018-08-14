package main

import (
	"fmt"

	"github.com/cheddartv/loom/config"
)

type Context struct {
	Config *config.Config
}

func main() {
	var context Context
	context.Config = config.Load()
	fmt.Printf("Got an output: %v\n", context.Config)
}
