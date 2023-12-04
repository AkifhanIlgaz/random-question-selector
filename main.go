package main

import (
	"fmt"

	"github.com/AkifhanIlgaz/random-question-selector/cfg"
)

func main() {
	config, err := cfg.LoadConfig(".")
	if err != nil {
		panic(err)
	}

	

	fmt.Println(config)
}
