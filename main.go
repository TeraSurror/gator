package main

import (
	"log"
	"os"

	"github.com/TeraSurror/gator/internal/config"
)

type state struct {
	cfg *config.Config
}

func main() {
	if len(os.Args) < 2 {
		log.Fatal("usage: cli <command> [args...]")
	}

	cfg, err := config.Read()
	if err != nil {
		log.Fatalf("error reading config: %v\n", err)
	}
	log.Printf("Read config: %+v\n", cfg)

	programState := &state{cfg: &cfg}
	cmds := commands{commandFuncMap: make(map[string]func(*state, command) error)}
	cmds.register("login", loginHandler)

	userCmd := os.Args[1]
	userCmdArgs := os.Args[2:]
	cmd := command{
		Name: userCmd,
		Args: userCmdArgs,
	}

	err = cmds.run(programState, cmd)
	if err != nil {
		log.Fatalf("%v", err)
	}
}
