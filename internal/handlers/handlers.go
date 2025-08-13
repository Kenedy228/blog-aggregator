package handlers

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/Kenedy228/blog-aggregator/internal/commands"
	"github.com/Kenedy228/blog-aggregator/internal/database"
	"github.com/google/uuid"
	"time"
)

func ResolveCommand(cmd commands.Command) error {
	var err error

	state := commands.NewState()
	cmds := commands.NewCommands()
	cmds.Register("login", handlerLogin)
	cmds.Register("register", handlerRegister)
	cmds.Register("reset", handlerReset)
	cmds.Register("users", handlerUsers)

	switch cmd.Name {
	case "login", "register", "reset", "users":
		err = cmds.Run(&state, cmd)
	default:
		err = fmt.Errorf("unknown command")
	}

	return err
}

func handlerLogin(s *commands.State, cmd commands.Command) error {
	if len(cmd.Args) == 0 {
		return fmt.Errorf("provide username")
	}

	if len(cmd.Args) > 1 {
		return fmt.Errorf("provide only one username")
	}

	ctx := context.Background()

	dbUser, err := s.DBQueries.GetUserByName(ctx, cmd.Args[0])

	if err != nil {
		if err == sql.ErrNoRows {
			return fmt.Errorf("User is not registered")
		}

		return fmt.Errorf("failed to login user: %v", err)
	}

	err = s.Cfg.SetUser(dbUser.Name)

	if err != nil {
		return err
	}

	fmt.Println("User has been set succesfully")
	return nil
}

func handlerRegister(s *commands.State, cmd commands.Command) error {
	if len(cmd.Args) == 0 {
		return fmt.Errorf("provide username")
	}

	if len(cmd.Args) > 1 {
		return fmt.Errorf("provide only one username")
	}

	ctx := context.Background()

	_, err := s.DBQueries.GetUserByName(ctx, cmd.Args[0])
	var dbUser database.User

	if err != nil {
		if err == sql.ErrNoRows {
			params := database.CreateUserParams{
				ID:        uuid.New(),
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
				Name:      cmd.Args[0],
			}

			dbUser, err = s.DBQueries.CreateUser(ctx, params)

			if err != nil {
				return fmt.Errorf("failed to register user: %v", err)
			}

			err = s.Cfg.SetUser(dbUser.Name)

			if err != nil {
				return fmt.Errorf("failed to set reigistered user: %v", err)
			}

			fmt.Println("User has been registered succesfully")
			fmt.Printf("{\n\tid: %v,\n\tcreated_at: %v,\n\tupdated_at: %v,\n\tname: %v,\n}\n",
				dbUser.ID,
				dbUser.CreatedAt,
				dbUser.UpdatedAt,
				dbUser.Name)

			return nil
		} else {
			return fmt.Errorf("failed register user: %v", err)
		}
	}

	return fmt.Errorf("User already exists")
}

func handlerReset(s *commands.State, cmd commands.Command) error {
	if len(cmd.Args) > 0 {
		return fmt.Errorf("reset command must have no arguments")
	}

	ctx := context.Background()
	err := s.DBQueries.DeleteAll(ctx)

	if err != nil {
		return fmt.Errorf("reset command error: %v", err)
	}

	fmt.Println("reset command successs")
	return nil
}

func handlerUsers(s *commands.State, cmd commands.Command) error {
	if len(cmd.Args) > 0 {
		return fmt.Errorf("users command must have no arguments")
	}

	ctx := context.Background()
	dbUsers, err := s.DBQueries.GetAllUsers(ctx)

	if err != nil {
		return fmt.Errorf("users command error: %v", err)
	}

	if len(dbUsers) == 0 {
		fmt.Println("There is no registered users")
		return nil
	}

	for _, dbUser := range dbUsers {
		data := fmt.Sprintf("* %s", dbUser.Name)

		if dbUser.Name == s.Cfg.Username {
			data += " (current)"
		}

		fmt.Println(data)
	}

	return nil
}
