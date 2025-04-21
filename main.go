package main

import (
	"database/sql"
	"log"
	"os"

	"github.com/TeraSurror/gator/internal/config"
	"github.com/TeraSurror/gator/internal/database"

	_ "github.com/lib/pq"
)

type state struct {
	db  *database.Queries
	cfg *config.Config
}

func main() {
	if len(os.Args) < 2 {
		log.Fatal("usage: cli <command> [args...]")
	}

	// Read configuration from config files
	cfg, err := config.Read()
	if err != nil {
		log.Fatalf("error reading config: %v\n", err)
	}

	// Create database connection
	db, err := sql.Open("postgres", cfg.DbUrl)
	if err != nil {
		log.Fatalf("could not connect to the database: %v\n", err)
	}

	// Create program state
	programState := &state{
		db:  database.New(db),
		cfg: &cfg,
	}

	// Register commands
	cmds := commands{commandFuncMap: make(map[string]func(*state, command) error)}
	cmds.register("login", loginHandler)
	cmds.register("register", registerHandler)
	cmds.register("reset", resetHandler)
	cmds.register("users", userListHandler)
	cmds.register("agg", aggHandler)
	cmds.register("addfeed", middlewareLoggedIn(addFeedHandler))
	cmds.register("feeds", feedHandler)
	cmds.register("follow", middlewareLoggedIn(followHandler))
	cmds.register("following", middlewareLoggedIn(followsHandler))

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
