package main

import (
	"github.com/abiosoft/ishell"
	"github.com/akler/twitter/src/domain"
	"github.com/akler/twitter/src/service"
)

func main() {

	tm := service.NewTweetManager(service.NewFileTweetWriter())

	shell := ishell.New()
	shell.SetPrompt("Tweeter >> ")
	shell.Print("Type 'help' to know commands\n")

	shell.AddCmd(&ishell.Cmd{
		Name: "publishTweet",
		Help: "Publishes a tweet",
		Func: func(c *ishell.Context) {

			defer c.ShowPrompt(true)

			c.Print("Write your tweet: ")

			tweet := domain.NewTextTweet("Usuario de shell", c.ReadLine())

			tm.PublishTweet(tweet)

			c.Print("Tweet sent\n")

			return
		},
	})

	shell.AddCmd(&ishell.Cmd{
		Name: "showTweet",
		Help: "Shows a tweet",
		Func: func(c *ishell.Context) {

			defer c.ShowPrompt(true)

			tweet := tm.GetTweet()

			c.Println(tweet.PrintableTweet())

			return
		},
	})

	shell.Run()

}
