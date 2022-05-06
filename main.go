package main

import (
	"flag"
	"fmt"
	"log"
	"time"

	"github.com/arshamalh/twigo"
	"github.com/go-co-op/gocron"
	"github.com/gofiber/fiber/v2"
	"github.com/sdomino/scribble"
)

type User struct {
	Name string
}

type TweetData struct {
	CreatedAt     time.Time
	AuthorID      string
	ID            string
	Text          string
	PublicMetrics map[int64]map[string]int
}

var tweet_authors = map[string]string{}
var bot *twigo.Client
var db *scribble.Driver
var scheduler *gocron.Scheduler

func main() {
	var bearer_token string
	var err error

	flag.StringVar(&bearer_token, "bearer", "", "Bearer token")

	bot, err = twigo.NewClient("", "", "", "", bearer_token)
	if err != nil {
		fmt.Println(err)
	}

	db, err = scribble.New("users", nil)
	if err != nil {
		fmt.Println("Error", err)
	}

	scheduler = gocron.NewScheduler(time.UTC)

	app := fiber.New()

	app.Post("/add_users_to_be_tracked", func(c *fiber.Ctx) error {
		var user_ids []string

		if err := c.BodyParser(&user_ids); err != nil {
			return err
		}

		for _, user_id := range user_ids {
			_, err := scheduler.Every(5).Minutes().Tag(fmt.Sprintf("user_track_%s", user_id)).Do(TrackUser, user_id)
			if err == nil {
				fmt.Printf("tracking user %s started\n", user_id)
			} else {
				fmt.Println(err)
			}
		}

		return c.JSON(true)
	})

	scheduler.StartAsync()
	log.Fatal(app.Listen(":3000"))
}

func CollectTweetData(tweet_id string) {
	fields := map[string]interface{}{"tweet.fields": []string{"author_id", "created_at", "public_metrics"}}
	response, _ := bot.GetTweet(tweet_id, false, fields)
	tweet := response.Data
	fmt.Printf("Collection tweet data: %s\n", tweet_id)
	item_id := fmt.Sprintf("%s_%s", tweet.AuthorID, tweet_id)
	if tweet.ID != "" {
		if _, ok := tweet_authors[tweet_id]; ok {
			tweet_data := &TweetData{}
			db.Read("users", item_id, tweet_data)
			fmt.Printf("%#v\n", tweet_data)
			tweet_data.PublicMetrics[time.Now().Unix()] = tweet.PublicMetrics
			db.Write("users", item_id, tweet_data)
		} else {
			tweet_data := &TweetData{
				CreatedAt:     tweet.CreatedAt,
				AuthorID:      tweet.AuthorID,
				ID:            tweet.ID,
				Text:          tweet.Text,
				PublicMetrics: map[int64]map[string]int{},
			}
			tweet_data.PublicMetrics[time.Now().Unix()] = tweet.PublicMetrics
			db.Write("users", item_id, tweet_data)
			tweet_authors[tweet_id] = tweet.AuthorID
		}
	}
}

func TrackUser(user_id string) {
	start_time := time.Now().UTC().Add(-5 * time.Minute)
	params := map[string]interface{}{"max_results": 5, "start_time": start_time, "tweet.fields": []string{"created_at"}}
	user_tweets, _ := bot.GetUserTweets(user_id, false, params)
	// user_tweets, _ := testGetUserTweets(user_id, false, params)
	if len(user_tweets.Data) != 0 { // It's surely not nil
		fmt.Printf("A tweet found: https://twitter.com/i/web/status/%s\n", user_tweets.Data[0].ID)
		scheduler.Every(5).Seconds().Do(CollectTweetData, user_tweets.Data[0].ID)
		if err := scheduler.RemoveByTag(fmt.Sprintf("user_track_%s", user_id)); err == nil {
			fmt.Printf("tracking user %s stopped\n", user_id)
		}
	}
}

// func testGetUserTweets(user_id string, oauth_1a bool, params map[string]interface{}) (*twigo.TweetsResponse, error) {
// 	return &twigo.TweetsResponse{
// 		Data: []twigo.Tweet{
// 			{
// 				ID:   "1519735033950470144",
// 				Text: "Hello, World!",
// 			},
// 		},
// 	}, nil
// }
