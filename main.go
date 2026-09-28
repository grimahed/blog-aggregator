package main

import(
"fmt"
"os"
"database/sql"
"blog-aggregator/internal/config"
"blog-aggregator/internal/database"
_ "github.com/lib/pq"
)



func main() {
	cfg, err := config.ReadFile()
	if err != nil {
		fmt.Println(err)
		return
	}

	db, err := sql.Open("postgres", cfg.Db_url)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	dbQueries := database.New(db)

	s := state{
		db: dbQueries,
		cfg: &cfg,
	}
	cmd := commands{
		Commands: make(map[string]func(*state, command) error),
		}
	com := command{}

	cmd.register("login", handlerLogin)
	cmd.register("register", handlerRegisterUser)
	cmd.register("reset", handlerReset)
	cmd.register("users", handlerUsers)
	cmd.register("agg", handlerAgg)
	cmd.register("addfeed", handlerAddFeed)
	if len(os.Args) < 2 {
		fmt.Println("I NEED a command first")
		os.Exit(1)
	}

	com = command {
		name: os.Args[1],
		args: os.Args[2:],
		}

	if err := cmd.run(&s, com); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
