package main

import (
	"fmt"

	"github.com/AkifhanIlgaz/random-question-selector/cfg"
	services "github.com/AkifhanIlgaz/random-question-selector/services/user"
)

func main() {
	config, err := cfg.LoadConfig(".")
	if err != nil {
		panic(err)
	}

	fmt.Println(config)

	
}
