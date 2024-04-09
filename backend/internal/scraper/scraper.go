package scraper

import (
	"context"
	"log"

	colly "github.com/gocolly/colly/v2"
)

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
