package handlers

import (
	"fmt"
	"github.com/Kenedy228/blog-aggregator/internal/commands"
)

func middlewareLoggedIn(handler func(*commands.State, commands.Command) error) func(*commands.State, commands.Command) error {
	return func(s *commands.State, cmd commands.Command) error {
		if s.Cfg.Username == "" {
			return fmt.Errorf("login first")
		}

		return handler(s, cmd)
	}
}
