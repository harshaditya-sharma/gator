package main

import (
	"database/sql"
	"fmt"
	"os"

	"github.com/harshaditya-sharma/gator/internal/commands"
	"github.com/harshaditya-sharma/gator/internal/config"
	"github.com/harshaditya-sharma/gator/internal/database"
	_ "github.com/lib/pq"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("error: not enough arguments. usage: gator <command>")
		os.Exit(1)
	}
	cfg, err := config.Read()
	if err != nil {
		fmt.Printf("error reading config: %v\n", err)
		os.Exit(1)
	}

	dbUrl := cfg.Db_url
	db, err := sql.Open("postgres", dbUrl)
	if err != nil {
		fmt.Printf("error opening DB: %v\n", err)
		os.Exit(1)
	}
	defer db.Close()

	if err := db.Ping(); err != nil {
		fmt.Printf("error pinging DB: %v\n", err)
		os.Exit(1)
	}

	dbQueries := database.New(db)

	state := &config.State{Db: dbQueries, Cfg: cfg}

	cmds := &commands.Commands{Handlers: make(map[string]func(*config.State, commands.Command) error)}
	cmds.Register("login", commands.HandlerLogin)
	cmds.Register("register", commands.HandlerRegister)
	cmds.Register("reset", commands.HandlerReset)
	cmds.Register("users", commands.HandlerUsers)
	cmds.Register("agg", commands.HandlerAgg)
	cmds.Register("feeds", commands.HandlerFeeds)
	cmds.Register("addfeed", commands.MiddlewareLoggedIn(commands.HandlerAddFeed))
	cmds.Register("follow", commands.MiddlewareLoggedIn(commands.HandlerFollow))
	cmds.Register("following", commands.MiddlewareLoggedIn(commands.HandlerFollowing))
	cmds.Register("unfollow", commands.MiddlewareLoggedIn(commands.HandlerUnfollow))
	cmds.Register("browse", commands.MiddlewareLoggedIn(commands.HandlerBrowse))
	cmds.Register("search", commands.MiddlewareLoggedIn(commands.HandlerSearch))
	cmds.Register("like", commands.MiddlewareLoggedIn(commands.HandlerLike))
	cmds.Register("unlike", commands.MiddlewareLoggedIn(commands.HandlerUnlike))
	cmds.Register("liked", commands.MiddlewareLoggedIn(commands.HandlerLiked))
	cmdName := os.Args[1]
	args := os.Args[2:]

	cmd := commands.Command{Name: cmdName, Args: args}

	if err := cmds.Run(state, cmd); err != nil {
		fmt.Printf("error: %v\n", err)
		os.Exit(1)
	}
}
