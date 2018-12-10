package domain

import (
	"time"
)

type Tweet interface {
	PrintableTweet() string
	ObtenerId() int
	ObtenerUser() string
	ObtenerText() string
	ObtenerDate() *time.Time
	SetId(id int)
}

type TextTweet struct {
	Tweet
	Id int
	User string
	Text string
	Date *time.Time
}

type ImageTweet struct {
	TextTweet
	Url string
}

type QuoteTweet struct {
	TextTweet
	Quoted Tweet
}

func NewTextTweet(user string, text string) (p *TextTweet) {
	d := time.Now()
	t := TextTweet{
		User: user,
		Text: text,
		Date: &d,
	}
	p = &t
	return
}

func NewImageTweet(user string, text string, url string) (p *ImageTweet) {
	t := ImageTweet{
		TextTweet: *NewTextTweet(user,text),
		Url: url,
	}
	p = &t
	return
}

func NewQuoteTweet(user string, text string, quoted Tweet) (p *QuoteTweet) {
	t := QuoteTweet{
		TextTweet: *NewTextTweet(user,text),
		Quoted: quoted,
	}
	p = &t
	return
}

func (t *TextTweet) SetId(id int)  {
	t.Id = id
}

func (t *TextTweet) PrintableTweet() string {
	return "@"+t.User+": "+t.Text
}

func (t *TextTweet) ObtenerId() int {
	return t.Id
}

func (t *TextTweet) ObtenerUser() string {
	return t.User
}
func (t *TextTweet) ObtenerText() string {
	return t.Text
}

func (t *TextTweet) ObtenerDate() *time.Time {
	return t.Date
}

func (t *ImageTweet) PrintableTweet() string {
	return "@"+t.User+": "+t.Text+" <a>"+t.Url+"</a>"
}

func (t *QuoteTweet) PrintableTweet() string {
	return "@"+t.User+": "+t.Text+` "`+t.Quoted.PrintableTweet()+`"`
}

func (t *TextTweet) String() string {
	return t.PrintableTweet()
}