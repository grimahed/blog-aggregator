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

func handlerAddFeed(s *state, cmd command, user database.User) error {
	if len(cmd.args) < 2 {
		return errors.New("Not enough args")
	}

	var feedName string = cmd.args[0]
	var feedURL string = cmd.args[1]

	newFeed := database.CreateFeedParams{
                        ID: uuid.New(),
                        CreatedAt: time.Now(),
                        UpdatedAt: time.Now(),
                        Name: feedName,
                        Url: feedURL,
			UserID: user.ID,
                }
	builtFeed, err := s.db.CreateFeed(context.Background(), newFeed)
	if err != nil {
		return err
	}

	folParams := database.CreateFeedFollowParams{
			ID: uuid.New(),
			CreatedAt: newFeed.CreatedAt,
			UpdatedAt: newFeed.UpdatedAt,
			UserID: newFeed.UserID,
			FeedID: newFeed.ID,
	}

	if _, err := s.db.CreateFeedFollow(context.Background(), folParams); err != nil {
		return err
	}
	fmt.Println(builtFeed)
return nil
}


func handlerFeeds(s *state, cmd command) error {

	feeds, err := s.db.GetAllFeedsWithUserName(context.Background())
	if err != nil {
		return err
	}

	for i := range feeds {
		fmt.Printf("%+v\n", feeds[i])
	}

return nil
}

func handlerFollow(s *state, cmd command, user database.User) error {
	if len(cmd.args) == 0 {
		return errors.New("Will Argue about anything for $1")
	}

	feed, err := s.db.GetFeed(context.Background(), cmd.args[0])
	if err != nil {
		return err
	}

	params := database.CreateFeedFollowParams{
			ID: uuid.New(),
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
			UserID: user.ID,
			FeedID: feed.ID,
	}

	followRow, err := s.db.CreateFeedFollow(context.Background(), params)
	if err != nil {
		return err
	}

	fmt.Printf("%v, %s\n", followRow.NameOfFeed, followRow.NameOfUser)
return nil
}

func handlerFollowing(s *state, cmd command, user database.User) error {

	following, err := s.db.GetFeedFollowsForUser(context.Background(), user.ID)
	if err != nil {
		return err
	}
	for i := range following {
	fmt.Println(following[i].FeedName)
	}
return nil
}

func handlerUnfollow(s *state, cmd command, user database.User) error {

if len(cmd.args) == 0 {
	return errors.New("I will make you argue")
}

	feed, err := s.db.GetFeed(context.Background(), cmd.args[0])
	if err != nil {
		return err
	}

	params := database.UnfollowParams{
		UserID: user.ID,
		FeedID: feed.ID,
	}

	s.db.Unfollow(context.Background(), params)

return nil
}
