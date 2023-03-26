package main

import (
	"flag"
	"fmt"
	"os"
	"time"

	"github.com/joho/godotenv"
	"github.com/zhienbaevsa/toggle2jira-go/internal/controller"
	"github.com/zhienbaevsa/toggle2jira-go/internal/repository/toggl"
)

const timeFormat string = "2006-01-02"

type config struct {
	Toggl struct {
		host   string
		apiKey string
	}
	Jira struct {
		host string
		user string
		pass string
	}
	Timezone string
}

func loadEnvConfig() (config, error) {
	var cfg config
	err := godotenv.Load()
	if err != nil {
		return cfg, err
	}

	cfg.Toggl.host = os.Getenv("TOGGL_HOST")
	cfg.Toggl.apiKey = os.Getenv("TOGGL_API_KEY")

	cfg.Jira.host = os.Getenv("JIRA_HOST")
	cfg.Jira.user = os.Getenv("JIRA_USER")
	cfg.Jira.pass = os.Getenv("JIRA_PASS")

	cfg.Timezone = os.Getenv("TIMEZONE")

	return cfg, nil
}

func main() {
	from := flag.String("from", "", "From date")
	to := flag.String("to", "", "To date")

	flag.Parse()

	fromDate, err := time.Parse(timeFormat, *from)
	if err != nil {
		panic(fmt.Sprintf("Could not parse from time: %v", *from))
	}

	var toDate time.Time
	if *to != "" {
		toDateParsed, err := time.Parse(timeFormat, *to)
		if err != nil {
			panic(fmt.Sprintf("Could not parse to time: %v", *to))
		}
		toDate = toDateParsed
	} else {
		now := time.Now()
		toDate = time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
	}

	cfg, err := loadEnvConfig()
	if err != nil {
		panic(err)
	}

	ts := &toggl.TogglTimeEntryStorage{
		ApiKey: cfg.Toggl.apiKey,
	}

	wu, err := controller.New(ts, cfg.Timezone, controller.JiraConfig{
		Host: cfg.Jira.host,
		User: cfg.Jira.user,
		Pass: cfg.Jira.pass,
	})
	if err != nil {
		panic(err) // TODO Proper error handling
	}
	err = wu.Start(fromDate, toDate)
	if err != nil {
		panic(err) // TODO Proper error handling
	}

	// Implement progress bar
	// bar := progressbar.Default(int64(len(ww)))
	// bar.Add(1)
}