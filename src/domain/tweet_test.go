package domain_test

import (
	"github.com/akler/twitter/src/domain"
	"testing"
)

func TestTextTweetPrintsUserAndText(t *testing.T){
	tweet := domain.NewTextTweet("user","this is my tw")

	text := tweet.PrintableTweet()

	expectedText := "@user: this is my tw"

	if text != expectedText {
		t.Errorf("Wrong expected test")
	}
}

func TestImageTweetPrintsUserTextAndUrl(t *testing.T){
	it := domain.NewImageTweet("user","text in tw","http://url.com")
	text := it.PrintableTweet()

	expectedText := "@user: text in tw <a>http://url.com</a>"

	if text != expectedText {
		t.Errorf("Wrong expected test")
	}
}

func TestQuoteTweetPrintUserTextAndQuotedTweet(t *testing.T)  {
	tweet := domain.NewTextTweet("user","this is my tw")
	quoteTweet := domain.NewQuoteTweet("another","awesome",tweet)

	expectedText := `@another: awesome "@user: this is my tw"`

	if quoteTweet.PrintableTweet() != expectedText {
		t.Errorf("Expected %s but was %s",expectedText,quoteTweet.PrintableTweet())
	}
}
