package main

import (
	"github.com/spf13/viper"
)

type Context struct {
}

func main() {
	viper.SetConfigName("loom")
	viper.AddConfigPath("/etc/")
}
