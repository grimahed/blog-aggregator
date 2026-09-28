package main

import (
"blog-aggregator/internal/database"
"github.com/google/uuid"
"errors"
"fmt"
"time"
"context"
"log"
)


func handlerLogin(s *state, cmd command) error {
        if len(cmd.args) == 0 {
                return errors.New("arrrg where are my arrrgs")
        }

        if _, err := s.db.GetUser(context.Background(), cmd.args[0]); err != nil {
                return errors.New("user does not exist :<")
        }

         _, err := s.cfg.SetUser(cmd.args[0])
        if err != nil {
                return err
        }

        fmt.Printf("user %v has logged in!\n", s.cfg.Current_user_name)
return nil
}

func handlerRegisterUser(s *state, cmd command) error {
        if len(cmd.args) == 0 {
                return errors.New("What pirate doesn't say aarrrg")
        }

        user := database.CreateUserParams {
                        ID: uuid.New(),
                        CreatedAt: time.Now(),
                        UpdatedAt: time.Now(),
                        Name:      cmd.args[0],
                }

        builtUser, err := s.db.CreateUser(context.Background(), user)
        if err != nil {
                return err
        }
        if _, err := s.cfg.SetUser(builtUser.Name); err != nil {
                return err
        }

        log.Println(builtUser)
        return nil
}

func handlerUsers(s *state, cmd command) error {
        users, err := s.db.GetUsers(context.Background())
        if err != nil {
                return err
        }

        for _, user := range users {
                if user.Name == s.cfg.Current_user_name {
                        fmt.Printf("* %s (current)\n", user.Name)
                } else {
                fmt.Printf("* %s\n", user.Name)
                }
        }

return nil
}
func handlerReset(s *state, cmd command) error {
        if err := s.db.DeleteAll(context.Background()); err != nil {
                return err
        }
return nil
}

func handlerAgg(s *state, cmd command) error {
        dummyURL := "https://www.wagslane.dev/index.xml"

        feed, err := fetchFeed(context.Background(), dummyURL)
        if err != nil {
                return err
        }
        fmt.Println(feed)
return nil
}

func handlerAddFeed(s *state, cmd command) error {
	if len(cmd.args) < 2 {
		return errors.New("Not enough args")
	}

	var feedName string = cmd.args[0]
	var feedURL string = cmd.args[1]

	user, err := s.db.GetUser(context.Background(), s.cfg.Current_user_name)
	if err != nil {
		return err
	}

	whyTheFUCKcanINotDoThisInTheSig := database.CreateFeedParams{
                        ID: uuid.New(),
                        CreatedAt: time.Now(),
                        UpdatedAt: time.Now(),
                        Name: feedName,
                        Url: feedURL,
			UserID: user.ID,
                }
	builtfeed, err := s.db.CreateFeed(context.Background(),
			whyTheFUCKcanINotDoThisInTheSig) //because you'd give everyone a
	if err != nil {					 //headache trying to read it
		return err
	}

	fmt.Println(builtfeed)
return nil
}
