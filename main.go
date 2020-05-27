package main

import (
	"strings"

	"github.com/BurntSushi/toml"
	"github.com/mmcdole/gofeed"
	"gopkg.in/gomail.v2"
)

type Config struct {
	From     string
	To       string
	Subject  string
	Host     string
	Port     int
	Username string
	Password string
	FeedUrl  string
	ReaderUrl  string
}

func sendMail(c Config, text string) {
	m := gomail.NewMessage()
	m.SetHeader("From", c.From)
	m.SetHeader("To", c.To)
	m.SetHeader("Subject", c.Subject)
	m.SetBody("text/plain", text)

	d := gomail.NewDialer(c.Host, c.Port, c.Username, c.Password)

	if err := d.DialAndSend(m); err != nil {
		panic(err)
	}
}

func fetchFeed(c Config) string {
	fp := gofeed.NewParser()

	feed, _ := fp.ParseURL(c.FeedUrl)
	items := feed.Items

	var sb strings.Builder

	for _, item := range items {
		sb.WriteString(item.Title)
		sb.WriteString("\n")
		sb.WriteString(item.Link)
		sb.WriteString("\n")
		sb.WriteString("\n")
	}

	sb.WriteString("\n")
	sb.WriteString("Read:")
	sb.WriteString(c.ReaderUrl)

	return sb.String()
}

func readConfig() (Config, error) {
	var conf Config
	if _, err := toml.DecodeFile("config.toml", &conf); err != nil {
		return conf, err
	}
	return conf, nil
}

func main() {
	c, err := readConfig()
	if err != nil {
		panic(err)
	}

	t := fetchFeed(c)
	if t != "" {
		sendMail(c, t)
	}
}
