package service

import (
	"github.com/akler/twitter/src/domain"
	"os"
)

type TweetWriter interface {
	WriteTweet(tweet domain.Tweet)
	GetLastSavedTweet() domain.Tweet
}

type MemoryTweetWriter struct {
	tweets []domain.Tweet
}

func (m *MemoryTweetWriter) WriteTweet(tweet domain.Tweet)  {
	m.tweets = append(m.tweets, tweet)
}

func (m *MemoryTweetWriter) GetLastSavedTweet() domain.Tweet {
	length := len(m.tweets)
	if length < 1{
		return nil
	}
	return m.tweets[length-1]
}


func NewMemoryTweetWriter() (tw TweetWriter)  {
	tw = &MemoryTweetWriter{
		tweets: make([]domain.Tweet,0),
	}
	return
}


type FileTweetWriter struct {
	file *os.File
}

func (writer *FileTweetWriter) WriteTweet(tweet domain.Tweet)  {
	go func() {
		if writer.file != nil {
			byteSlide := []byte(tweet.PrintableTweet()+"\n")
			writer.file.Write(byteSlide)
		}
	}()
}

func (writer *FileTweetWriter) GetLastSavedTweet() domain.Tweet {
	return nil
}


func NewFileTweetWriter() (writer *FileTweetWriter)  {
	file,_ := os.OpenFile(
		"tweet.txt",
		os.O_WRONLY|os.O_TRUNC|os.O_CREATE,
		0666)
	writer = &FileTweetWriter{
		file:file,
	}
	return
}

