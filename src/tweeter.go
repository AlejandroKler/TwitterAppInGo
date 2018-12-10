package main

import (
	"github.com/abiosoft/ishell"
	"github.com/akler/twitter/src/domain"
	"github.com/akler/twitter/src/service"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

func main() {

	tm := service.NewTweetManager(service.NewFileTweetWriter())

	ginServer := NewGinServer(tm)
	ginServer.StartGinServer()

	initConsole(tm)
}

type GinServer struct {
	tm *service.TweetManager
}

func NewGinServer(tm *service.TweetManager) (gs *GinServer)  {
	return &GinServer{tm}
}

func (gs *GinServer) StartGinServer(){
	router := gin.Default()

	router.GET("/tweet/:id",gs.getTweetById)
	router.GET("/tweet/",gs.getAllTweets)
	router.PUT("/tweet/",gs.newTextTweet)
	router.PUT("/quote/",gs.newQuoteTweet)

	go router.Run()
}

func (gs *GinServer) getTweetById(c *gin.Context) {
	id,_ := strconv.Atoi(c.Param("id"))
	c.JSON(http.StatusOK,gs.tm.GetTweetById(id))
}

func (gs *GinServer) getAllTweets(c *gin.Context) {
	c.JSON(http.StatusOK,gs.tm.GetTweets())
}

func (gs *GinServer) newQuoteTweet(c *gin.Context) {
	var request struct{
		User string `json:"user"`
		Text string `json:"text"`
		QuoteId int `json:"quote"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest,gin.H{"error":err.Error()})
		return
	}

	tw := domain.NewQuoteTweet(request.User,request.Text,gs.tm.GetTweetById(request.QuoteId))
	_,err := gs.tm.PublishTweet(tw)

	if err != nil {
		c.JSON(http.StatusBadRequest,gin.H{"error":err.Error()})
		return
	}

	c.JSON(http.StatusOK,tw)
}

func (gs *GinServer) newTextTweet(c *gin.Context) {
	var request struct{
		User string `json:"user"`
		Text string `json:"text"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest,gin.H{"error":err.Error()})
		return
	}

	tw := domain.NewTextTweet(request.User,request.Text)
	gs.tm.PublishTweet(tw)

	c.JSON(http.StatusOK,tw)
}


func (gs *GinServer) getLastTweet(c *gin.Context) {
	c.JSON(http.StatusOK, gs.tm.GetTweet())
}



func initConsole(tm *service.TweetManager){
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