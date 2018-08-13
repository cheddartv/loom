package main

import (
	"fmt"

	"github.com/spf13/viper"
)

type manifest struct {
	Output string
	Inputs []string
}

type config struct {
	Manifests []manifest
}

type Context struct {
	Config config
}

func main() {
	viper.SetConfigName("loom")

	viper.AddConfigPath("/etc/")
	viper.AddConfigPath("./etc")

	err := viper.ReadInConfig()
	if err != nil {
		panic(fmt.Errorf("Fatal error reading config file: %s \n", err))
	}

	var context Context
	err = viper.Unmarshal(&context.Config)
	if err != nil {
		panic(fmt.Errorf("Fatal error reading config file: %s \n", err))
	}

	fmt.Printf("Got an output: %v\n", context.Config)
}
