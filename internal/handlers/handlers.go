package handlers

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/Kenedy228/blog-aggregator/internal/commands"
	"github.com/Kenedy228/blog-aggregator/internal/database"
	"github.com/Kenedy228/blog-aggregator/internal/rss"
	"github.com/Kenedy228/blog-aggregator/internal/utility"
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
	cmds.Register("users", middlewareLoggedIn(handlerUsers))
	cmds.Register("agg", middlewareLoggedIn(handlerAgg))
	cmds.Register("addfeed", middlewareLoggedIn(handlerAddFeed))
	cmds.Register("feeds", middlewareLoggedIn(handlerFeeds))
	cmds.Register("follow", middlewareLoggedIn(handlerFollow))
	cmds.Register("following", middlewareLoggedIn(handlerFollowing))
	cmds.Register("unfollow", middlewareLoggedIn(handlerUnfollowing))

	if !cmds.CheckCommand(cmd.Name) {
		return fmt.Errorf("unknown command")
	}

	err = cmds.Run(&state, cmd)

	return err
}

func handlerLogin(s *commands.State, cmd commands.Command) error {
	if len(cmd.Args) == 0 {
		return fmt.Errorf("provide username")
	}

	if len(cmd.Args) > 1 {
		return fmt.Errorf("provide only one username")
	}

	ctx, cancel := utility.GenerateContextWithTimeout(utility.DB)
	defer cancel()

	dbUser, err := s.DBQueries.GetUserByName(ctx, cmd.Args[0])

	if err != nil {
		if err == sql.ErrNoRows {
			return fmt.Errorf("User is not registered")
		}

		if err == context.DeadlineExceeded {
			return fmt.Errorf("server is busy")
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

	ctx, cancel := utility.GenerateContextWithTimeout(utility.DB)
	defer cancel()

	params := database.CreateUserParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Name:      cmd.Args[0],
	}

	dbUser, err := s.DBQueries.CreateUser(ctx, params)

	if err != nil {
		if err == context.DeadlineExceeded {
			return fmt.Errorf("server is busy")
		}
		return fmt.Errorf("user already registered")
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
}

func handlerReset(s *commands.State, cmd commands.Command) error {
	if len(cmd.Args) > 0 {
		return fmt.Errorf("reset command must have no arguments")
	}

	ctx, cancel := utility.GenerateContextWithTimeout(utility.DB)
	defer cancel()

	err := s.DBQueries.DeleteAll(ctx)

	if err != nil {
		if err == context.DeadlineExceeded {
			return fmt.Errorf("server is busy")
		}
		return fmt.Errorf("reset command error: %v", err)
	}

	err = s.Cfg.DeleteUser()

	if err != nil {
		return err
	}

	fmt.Println("reset command successs")
	return nil
}

func handlerUsers(s *commands.State, cmd commands.Command) error {
	if len(cmd.Args) > 0 {
		return fmt.Errorf("users command must have no arguments")
	}

	ctx, cancel := utility.GenerateContextWithTimeout(utility.DB)
	defer cancel()

	dbUsers, err := s.DBQueries.GetAllUsers(ctx)

	if err != nil {
		if err == context.DeadlineExceeded {
			return fmt.Errorf("server is busy")
		}
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

func handlerAgg(s *commands.State, cmd commands.Command) error {
	if len(cmd.Args) == 0 {
		return fmt.Errorf("provide url")
	}

	if len(cmd.Args) > 1 {
		return fmt.Errorf("agg command receives one argument")
	}

	ctx, cancel := utility.GenerateContextWithTimeout(utility.HTTP)
	defer cancel()

	feed, err := rss.FetchFeed(ctx, cmd.Args[0])

	if err != nil {
		return err
	}

	fmt.Printf("received feed:\n%v", feed)

	return nil
}

func handlerAddFeed(s *commands.State, cmd commands.Command) error {
	if len(cmd.Args) < 2 {
		return fmt.Errorf("provide feed name and feed url for addfeed command")
	}

	if len(cmd.Args) > 2 {
		return fmt.Errorf("provide only one feed name and only one feed url for addfeed command")
	}

	if s.Cfg.Username == "" {
		return fmt.Errorf("login first")
	}

	reqCtx, reqCancel := utility.GenerateContextWithTimeout(utility.HTTP)
	defer reqCancel()

	_, err := rss.FetchFeed(reqCtx, cmd.Args[1])

	if err != nil {
		return err
	}

	dbCtx, dbCancel := utility.GenerateContextWithTimeout(utility.DB)
	defer dbCancel()

	params := database.CreateFeedParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Name:      cmd.Args[0],
		Url:       cmd.Args[1],
		Username:  s.Cfg.Username,
	}

	userID, err := s.DBQueries.CreateFeed(dbCtx, params)

	if err != nil {
		return fmt.Errorf("%v", err)
	}

	fmt.Printf("{\n\tid: %v,\n\tcreated_at: %v,\n\tupdated_at: %v,\n\tname: %v,\n\turl: %v,\n\tuser_id: %v,\n}\n",
		params.ID,
		params.CreatedAt,
		params.UpdatedAt,
		params.Name,
		params.Url,
		userID.UUID,
	)

	return nil
}

func handlerFeeds(s *commands.State, cmd commands.Command) error {
	if len(cmd.Args) > 0 {
		return fmt.Errorf("command feeds must have no parameters")
	}

	ctx, cancel := utility.GenerateContextWithTimeout(utility.DB)
	defer cancel()

	feeds, err := s.DBQueries.GetFeedsWithUsers(ctx)

	if err != nil {
		return fmt.Errorf("there is no data")
	}

	for _, v := range feeds {
		fmt.Printf("* feed=%v, url=%v, username=%v\n", v.Name, v.Url, v.Name_2)
	}

	return nil
}

func handlerFollow(s *commands.State, cmd commands.Command) error {
	if len(cmd.Args) == 0 {
		return fmt.Errorf("provide url to follow")
	}

	if len(cmd.Args) > 1 {
		return fmt.Errorf("provide only one url to follow")
	}

	if s.Cfg.Username == "" {
		return fmt.Errorf("login first")
	}

	ctx, close := utility.GenerateContextWithTimeout(utility.DB)
	defer close()

	params := database.CreateFeedFollowParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		UserName:  s.Cfg.Username,
		Url:       cmd.Args[0],
	}

	row, err := s.DBQueries.CreateFeedFollow(ctx, params)

	if err != nil {
		if err == context.DeadlineExceeded {
			return fmt.Errorf("server is busy")
		}

		return fmt.Errorf("error performing query handlerFollow: %v", err)
	}

	fmt.Printf("{\n\tid: %vm\n\tcreated_at: %v,\n\tupdated_at: %v,\n\tusername: %v,\n\turl: %v,\n}\n",
		row.ID, row.CreatedAt, row.UpdatedAt, row.UserName, row.FeedName)
	return nil
}

func handlerFollowing(s *commands.State, cmd commands.Command) error {
	if len(cmd.Args) > 0 {
		return fmt.Errorf("command following must have no args")
	}

	ctx, close := utility.GenerateContextWithTimeout(utility.DB)
	defer close()

	rows, err := s.DBQueries.GetFeedFollowsForUser(ctx, s.Cfg.Username)

	if err != nil {
		if err == context.DeadlineExceeded {
			return fmt.Errorf("server is busy")
		}

		return fmt.Errorf("error performing handlerFollowing: %v", err)
	}

	fmt.Printf("* user: %v\n", s.Cfg.Username)

	if len(rows) == 0 {
		fmt.Printf("no feeds follow")
		return nil
	}

	for _, row := range rows {
		fmt.Printf("- %v %v\n", row.FeedName, row.FeedUrl)
	}

	return nil
}

func handlerUnfollowing(s *commands.State, cmd commands.Command) error {
	if len(cmd.Args) == 0 {
		return fmt.Errorf("provide a feed url to unfollow")
	}

	if len(cmd.Args) > 1 {
		return fmt.Errorf("provide only one feed url to unfollow")
	}

	params := database.DeleteFeedFollowsForUserByUrlParams{
		Url:      cmd.Args[0],
		UserName: s.Cfg.Username,
	}

	ctx, close := utility.GenerateContextWithTimeout(utility.DB)
	defer close()

	err := s.DBQueries.DeleteFeedFollowsForUserByUrl(ctx, params)

	if err != nil {
		if err == context.DeadlineExceeded {
			return fmt.Errorf("server is busy")
		}

		return fmt.Errorf("error handlerUnfollowing: %v", err)
	}

	fmt.Printf("user %v succesfully unfollowed feed by url %v", s.Cfg.Username, cmd.Args[0])
	return nil
}
