package service_test

import (
	"strings"
	"testing"
	"github.com/akler/twitter/src/domain"
	"github.com/akler/twitter/src/service"
)

func TestPublishedTweetIsSaved(t *testing.T) {
	tm := service.NewTweetManager(service.NewMemoryTweetWriter())

	var tweet domain.Tweet
	user := "akler"
	text := "mi primer tw"
	tweet = domain.NewTextTweet(user, text)

	tm.PublishTweet(tweet)

	publishedTweet := tm.GetTweet()
	if publishedTweet.ObtenerUser() != user || publishedTweet.ObtenerText() != text {
		t.Errorf("Expected tweet is %s: %v \nbut is %s: %s", user, tweet, publishedTweet.ObtenerUser(), publishedTweet.ObtenerText())
		return
	}
	if publishedTweet.ObtenerDate() == nil {
		t.Error("Expected date can not be nil")
		return
	}
}

func TestTweetWithoutUserIsNotPublished(t *testing.T) {
	tm := service.NewTweetManager(service.NewMemoryTweetWriter())

	var tweet domain.Tweet
	var user string

	text := "contenido"
	tweet = domain.NewTextTweet(user, text)

	var e error
	_,e = tm.PublishTweet(tweet)

	if e == nil || e.Error() != "user is required" {
		t.Error("expected error user is required")
	}

}

func TestTweetWithoutTextIsNotPublished(t *testing.T) {
	tm := service.NewTweetManager(service.NewMemoryTweetWriter())

	var tweet domain.Tweet
	user := "user"

	text := ""
	tweet = domain.NewTextTweet(user, text)

	var e error
	_,e = tm.PublishTweet(tweet)

	if e == nil || e.Error() != "text is required" {
		t.Error("expected error text is required")
	}
}

func TestTweetWhichExceeding140CharactersIsNotPubliched(t *testing.T) {
	tm := service.NewTweetManager(service.NewMemoryTweetWriter())

	var tweet domain.Tweet
	user := "usuario"

	text := `estfsfdfgsdgfsdfgsdfg
	dsfgsdfgsdgsdfgsdggdsdsdfgsdgssgfsg
	dfgsdfgsdgsdfgsdfgsdfgsdgdsfgsdfggsdgfg
	sddfgsdfgsdgsdgsdgsdgsdgsdgsdfgsdg
	fgdfgsdgdsfgsdfgsdg876865tiu6tu6u6u6uu66u6`
	tweet = domain.NewTextTweet(user, text)

	var e error
	_,e = tm.PublishTweet(tweet)

	if e == nil || e.Error() != "max lenght" {
		t.Error("expected error max lenght")
	}
}

func TestCanPublishAndRetrieveMoreThanOneTweet(t *testing.T){
	tm := service.NewTweetManager(service.NewMemoryTweetWriter())
	var tweet, secondTweet domain.Tweet
	user1 := "user1"
	user2 := "user2"
	text1 := "text 1"
	text2 := "text 2"

	tweet = domain.NewTextTweet(user1,text1)
	secondTweet = domain.NewTextTweet(user2,text2)

	tm.PublishTweet(tweet)
	tm.PublishTweet(secondTweet)

	publishedTweets := tm.GetTweets()
	if len(publishedTweets) != 2 {
		t.Errorf("Expected size 2 but was %d",len(publishedTweets))
		return
	}
	firstPublishedTweet := publishedTweets[0]
	secondPublishedTweet := publishedTweets[1]

	if !isValidTweet(firstPublishedTweet,user1,text1){
		t.Errorf("Tweet 1 invalido")
		return
	}
	if !isValidTweet(secondPublishedTweet,user2,text2){
		t.Errorf("Tweet 2 invalido")
		return
	}

}

func isValidTweet(tw domain.Tweet, user string, text string) bool {
	if tw.ObtenerText() == text && tw.ObtenerUser() == user {
		return true
	}
	return false
}

func TestCanRetrieveTweetById(t *testing.T){
	tm := service.NewTweetManager(service.NewMemoryTweetWriter())

	var tweet domain.Tweet
	var id int

	user := "usuario"
	text := "holaa "

	tweet = domain.NewTextTweet(user,text)

	id,_ = tm.PublishTweet(tweet)

	publishedTweet := tm.GetTweetById(id)

	if publishedTweet.ObtenerText() != text || publishedTweet.ObtenerUser() != user {
		t.Errorf("Error on id")
	}
}

func TestCanCountTheTweetsSentByAnUser(t *testing.T) {
	tm := service.NewTweetManager(service.NewMemoryTweetWriter())

	var tweet, secondTweet, thirdTweet domain.Tweet
	user := "usuario"
	anotherUser := "usuario 2"

	tweet = domain.NewTextTweet(user,"algo")
	secondTweet = domain.NewTextTweet(user,"segundo tw")
	thirdTweet = domain.NewTextTweet(anotherUser,"algo")

	tm.PublishTweet(tweet)
	tm.PublishTweet(secondTweet)
	tm.PublishTweet(thirdTweet)

	count := tm.CountTweetsByUser(user)
	if count != 2 {
		t.Errorf("Should be 2 but was %d",count)
	}
}

func TestCanRetrieveTheTweetsSentByAnUser(t *testing.T) {
	tm := service.NewTweetManager(service.NewMemoryTweetWriter())

	var tweet, secondTweet, thirdTweet domain.Tweet
	user := "usuario"
	anotherUser := "usuario 2"

	tweet = domain.NewTextTweet(user,"algo")
	secondTweet = domain.NewTextTweet(user,"segundo tw")
	thirdTweet = domain.NewTextTweet(anotherUser,"algo")

	tm.PublishTweet(tweet)
	tm.PublishTweet(secondTweet)
	tm.PublishTweet(thirdTweet)

	tweets := tm.GetTweetsByUser(user)

	if len(tweets) != 2 {
		t.Errorf("Len should be 2")
		return
	}
	if tweets[0] != tweet || tweets[1] != secondTweet {
		t.Errorf("Len should be 2")
		return
	}

}


func TestPublishedTweetIsSavedToExternalResource(t *testing.T)  {
	var tweetWriter service.TweetWriter
	tweetWriter = service.NewMemoryTweetWriter()
	tweetManager := service.NewTweetManager(tweetWriter)

	var tweet domain.Tweet
	tweet = domain.NewTextTweet("ale","jajaja esto es go")
	id,_ := tweetManager.PublishTweet(tweet)

	memoryWriter := (tweetWriter).(*service.MemoryTweetWriter)
	savedTweet := memoryWriter.GetLastSavedTweet()

	if savedTweet == nil {
		t.Errorf("shuld not be nil")
		return
	}

	if savedTweet.ObtenerId() != id {
		t.Errorf("Wrong tweet id")
		return
	}
}

func TestCanSearchForTweetContainingText(t *testing.T)  {
	tm := service.NewTweetManager(service.NewMemoryTweetWriter())

	tm.PublishTweet(domain.NewTextTweet("user1","tweet first tweet"))
	tm.PublishTweet(domain.NewTextTweet("user1","second tweet"))
	tm.PublishTweet(domain.NewTextTweet("user1","third tweet"))
	tm.PublishTweet(domain.NewTextTweet("user1","tweet"))

	searchResult := make(chan domain.Tweet)
	query := "first"
	tm.SearchTweetsContaining(query,searchResult)

	foundTweet := <- searchResult

	if foundTweet == nil{
		t.Errorf("Should not be nil")
		return
	}
	if !strings.Contains(foundTweet.ObtenerText(),query) {
		t.Errorf("Should cointein query")
		return
	}

}