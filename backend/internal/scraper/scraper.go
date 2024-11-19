package scraper

import (
	"context"
	"crypto/sha1"
	"fmt"
	"log"

	colly "github.com/gocolly/colly/v2"
)

type ScraperFunc func(context.Context, *log.Logger, *log.Logger) []Article

type Article struct {
	Title   string
	Topics  []string
	URL     string
	Paywall bool
}

type collectorQuery struct {
	url             string
	querySelector   string
	queryCallbackFn func(*[]Article) colly.HTMLCallback
}

func scrape(_ context.Context, infoLog, errorLog *log.Logger, q collectorQuery) []Article {
	// TODO: Handle a cancellation of the context.

	collector := colly.NewCollector()

	collector.OnRequest(func(r *colly.Request) {
		infoLog.Printf("Visiting: %s", r.URL)
	})
	collector.OnError(func(r *colly.Response, err error) {
		errorLog.Printf("An error occurred while scraping %s: %s", r.Request.URL, err)
	})

	articles := make([]Article, 0, 200)
	collector.OnHTML(q.querySelector, q.queryCallbackFn(&articles))
	collector.Visit(q.url)

	return articles
}

func (a Article) CreateHash() (string, error) {
	if a.Title == "" || a.URL == "" {
		return "", fmt.Errorf("error generating hash for article, either title or url are missing")
	}

	articleInformation := fmt.Sprintf("%s%s", a.Title, a.URL)
	h := sha1.New()
	h.Write([]byte(articleInformation))
	hashOutput := h.Sum(nil)

	output := fmt.Sprintf("%x\n", hashOutput)

	output = output[:20]
	return output, nil
}
