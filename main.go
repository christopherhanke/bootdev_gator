package main

import (
	"fmt"
	"log"
	"os"

	"github.com/christopherhanke/bootdev_gator/internal/config"
)

func main() {
	cfg, err := config.Read()
	if err != nil {
		log.Fatalf("error reading config: %v\n", err)
	}
	fmt.Printf("Read config: %+v\n", cfg)

	var currState state
	bufferCfg, err := config.Read()
	if err != nil {
		log.Fatalf("error reading config: %v\n", err)
	} else {
		currState.cfg = &bufferCfg
		fmt.Printf("Read state: %+v\n", *currState.cfg)
	}
	commands := commands{
		make(map[string]func(*state, command) error),
	}
	commands.register("login", handlerLogin)

	args := os.Args
	if len(args) < 2 {
		log.Fatalf("Not enough arguments to run. Args given: %v\n", len(args))
		return
	}

	cmd := command{
		name: args[1],
		args: args[2:],
	}

	err = commands.run(&currState, cmd)
	if err != nil {
		fmt.Printf("Error running command: %v\n", cmd.name)
		log.Fatal(err)
	}

	cfg, err = config.Read()
	if err != nil {
		log.Fatalf("error reading config: %v\n", err)
	}
	fmt.Printf("Read config again: %+v\n", cfg)
}
