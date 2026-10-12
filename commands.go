package main

import (
"blog-aggregator/internal/config"
"blog-aggregator/internal/database"
"github.com/google/uuid"
"github.com/lib/pq"
"database/sql"
"errors"
"fmt"
"time"
"log"
"context"
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
		feed.Channel.Item[i].Description = html.UnescapeString(feed.Channel.Item[i].Description)
	}
return &feed, nil
}

func LoggedIn(handler func(s *state, cmd command, user database.User) error) func(*state, command) error {
	return func(s *state, cmd command) error {
		user, err := s.db.GetUser(context.Background(), s.cfg.Current_user_name)
		if err != nil {
			log.Fatal("failed to fetch user", err)
		}
		return handler(s, cmd, user)
	}
}

func scrapeFeeds(s *state) error {


		feed, err := s.db.GetNextFeedToFetch(context.Background())
		if err != nil {
			return err
		}

		markedFeed, err := s.db.MarkFeedFetched(context.Background(), feed.ID)
		if err != nil {
			return err
		}


		fetchedFeed, err := fetchFeed(context.Background(), markedFeed.Url)
		if err != nil {
			return err
		}

		for j := range fetchedFeed.Channel.Item {
			nullDesc := sql.NullString{}

			if fetchedFeed.Channel.Item[j].Description != "" {
					nullDesc.String = fetchedFeed.Channel.Item[j].Description
					nullDesc.Valid = true
			}

			nullTime := sql.NullTime{}
			parsedTime, err := time.Parse(time.DateTime, fetchedFeed.Channel.Item[j].PubDate)
			if err == nil {
				nullTime.Time = parsedTime
				nullTime.Valid = true
			}

			newPost := database.CreatePostParams{
					ID: uuid.New(),
					CreatedAt: time.Now(),
					UpdatedAt: time.Now(),
					Title: fetchedFeed.Channel.Item[j].Title,
					Description: nullDesc,
					PublishedAt: nullTime,
					FeedID: feed.ID,
			}

		_, err = s.db.CreatePost(context.Background(), newPost)
		if err != nil {
			pgErr, kindaValid := err.(*pq.Error)
			if kindaValid && pgErr.Code == "23505" {
				continue
			} else {
			log.Println(err)
			}
		}
		}
return nil
}

