package main

import (
"blog-aggregator/internal/config"
"blog-aggregator/internal/database"
"github.com/google/uuid"
"errors"
"fmt"
"time"
"context"
"log"
"net/http"
"io"
"html"
"encoding/xml"
)

var client = &http.Client{Timeout: 10 * time.Second}

type state struct {
	db  *database.Queries
	cfg *config.Config
}

type RSSFeed struct {
	Channel struct {
		Title	     string      `xml:"title"`
		Link	     string      `xml:"link"`
		Description  string      `xml:"description"`
		Item	     []RSSItem   `xml:"item"`
	} `xml:"channel"`
}

type RSSItem struct {
	Title        string `xml:"title"`
	Link         string  `xml:"link"`
	Description  string `xml:"description"`
	PubDate	     string `xml:"pubDate"`
}

type command struct {
	name string
	args []string
}

type commands struct {
	Commands map[string]func(*state, command) error
}


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
			Name:	   cmd.args[0],
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

func (c *commands) run(s *state, cmd command) error {
	handler, exists := c.Commands[cmd.name]
	if !exists {
		return errors.New("command or handler does not exist")
	}
	if err := handler(s, cmd); err != nil {
		return err
	}
	return nil
}

func (c *commands) register(name string, f func(*state, command) error) {
	_, exists := c.Commands[name]
	if exists {
		fmt.Println("This command already exists")
	} else {
	c.Commands[name] = f
	}
}

func fetchFeed(ctx context.Context, feedURL string) (*RSSFeed, error) {
	req, err := http.NewRequestWithContext(
			ctx, http.MethodGet, feedURL, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("User-Agent", "gator")

	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	if res.StatusCode > 299 {
		return nil, fmt.Errorf("WRONG NUMBER CODE %d", res.StatusCode)
	}
	defer res.Body.Close()

	bytes, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	//create the thing
	feed := RSSFeed{}
	//use thing below

	if err := xml.Unmarshal(bytes, &feed); err != nil {
		return nil, err
	}

	feed.Channel.Title = html.UnescapeString(feed.Channel.Title)
	feed.Channel.Description = html.UnescapeString(feed.Channel.Description)
	for i := range feed.Channel.Item {
		feed.Channel.Item[i].Title = html.UnescapeString(feed.Channel.Item[i].Title)
		feed.Channel.Item[i].Title = html.UnescapeString(feed.Channel.Item[i].Description)
	}
return &feed, nil
}
