package service

import (
	"fmt"
	"github.com/akler/twitter/src/domain"
	"strings"
)

type TweetManager struct {
	tweets []domain.Tweet
	tweetsByUser map[string][]domain.Tweet
	writer TweetWriter
}

func NewTweetManager(writer TweetWriter) (tm *TweetManager){
	tm = &TweetManager{
		tweets : make([]domain.Tweet,0),
		tweetsByUser : make(map[string][]domain.Tweet),
		writer: writer,
	}
	return
}

func (tm *TweetManager) GetTweets() []domain.Tweet {
	return tm.tweets
}

func (tm *TweetManager) GetTweet() domain.Tweet {
	cant := len(tm.tweets)
	if cant >= 1 {
		return tm.tweets[len(tm.tweets)-1]
	}
	return nil
}

func (tm *TweetManager) PublishTweet(t domain.Tweet) (int,error) {
	if t.ObtenerUser() == "" {
		return 0,fmt.Errorf("user is required")
	}
	if t.ObtenerText() == "" {
		return 0,fmt.Errorf("text is required")
	}
	if len(t.ObtenerText()) > 140 {
		return 0,fmt.Errorf("max lenght")
	}
	tm.tweets = append(tm.tweets, t)
	tm.tweetsByUser[t.ObtenerUser()] = append(tm.tweetsByUser[t.ObtenerUser()], t)
	t.SetId(len(tm.tweets))
	tm.writer.WriteTweet(t)
	return len(tm.tweets),nil
}

func (tm *TweetManager) GetTweetById(id int) (tw domain.Tweet) {
	if len(tm.tweets) < 1 || id > len(tm.tweets)  {
		return nil
	}
	return tm.tweets[id-1]
}

func (tm *TweetManager) CountTweetsByUser(user string) (count int)  {
	for _,tw := range tm.tweets  {
		if tw.ObtenerUser() == user {
			count++
		}
	}
	return count
}

func (tm *TweetManager) GetTweetsByUser(user string) []domain.Tweet {
	return tm.tweetsByUser[user]
}

func (tm *TweetManager) SearchTweetsContaining(query string, c chan domain.Tweet) {
	go func() {
		for _,tw := range tm.tweets {
			if strings.Contains(tw.ObtenerText(),query){
				fmt.Printf(tw.ObtenerText())
				c <- tw
			}
		}
	}()
}
