package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	"github.com/christopherhanke/bootdev_gator/internal/config"
	"github.com/christopherhanke/bootdev_gator/internal/database"
	_ "github.com/lib/pq"
)

func main() {
	var currState state
	bufferCfg, err := config.Read()
	if err != nil {
		log.Fatalf("error reading config: %v\n", err)
	} else {
		currState.cfg = &bufferCfg
		fmt.Printf("Read state: %+v\n", *currState.cfg)
	}
	db, err := sql.Open("postgres", currState.cfg.DBURL)
	if err != nil {
		log.Fatalf("error open database: %v\n", err)
	}
	defer db.Close()
	dbQueries := database.New(db)
	currState.db = dbQueries

	commands := commands{
		make(map[string]func(*state, command) error),
	}
	commands.register("login", handlerLogin)
	commands.register("register", handlerRegister)
	commands.register("reset", handlerReset)

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

	cfg, err := config.Read()
	if err != nil {
		log.Fatalf("error reading config: %v\n", err)
	}
	fmt.Printf("Read config again: %+v\n", cfg)
}
