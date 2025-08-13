package commands

import (
	"database/sql"
	"log"

	"github.com/Kenedy228/blog-aggregator/internal/config"
	"github.com/Kenedy228/blog-aggregator/internal/database"
)

type State struct {
	Cfg       *config.Config
	DBQueries *database.Queries
}

type Command struct {
	Name string
	Args []string
}

type Commands struct {
	cmds map[string]func(*State, Command) error
}

func (c *Commands) Run(s *State, cmd Command) error {
	err := c.cmds[cmd.Name](s, cmd)

	if err != nil {
		return err
	}

	return nil
}

func (c *Commands) Register(name string, f func(*State, Command) error) {
	c.cmds[name] = f
}

func NewState() State {
	cfg, err := config.Read()

	if err != nil {
		log.Fatalf("%v\n", err)
	}

	db, err := sql.Open("postgres", cfg.DBurl)

	if err != nil {
		log.Fatalf("db dsn error: %v\n", err)
	}

	dbQueries := database.New(db)

	state := State{Cfg: &cfg, DBQueries: dbQueries}

	return state
}

func NewCommand(name string, args []string) Command {
	cmd := Command{Name: name, Args: args}

	return cmd
}

func NewCommands() Commands {
	commands := Commands{cmds: make(map[string]func(*State, Command) error)}
	return commands
}
