package scraper

import (
	"context"
	"html"
	"log"

	"github.com/mmcdole/gofeed"
)

const nacioDigitalrssFeed = "https://naciodigital.cat/rss"

func scrapeNacioDigital(ctx context.Context, infoLog, errorLog *log.Logger) []Article {
	infoLog.Printf("Visiting: %s", nacioDigitalrssFeed)
	fp := gofeed.NewParser()
	feed, err := fp.ParseURLWithContext(nacioDigitalrssFeed, ctx)
	if err != nil {
		errorLog.Printf("an error occurred while scraping naciodigital: %s", err.Error())
		return []Article{}
	}
	articles := make([]Article, 0)
	for _, item := range feed.Items {
		title := html.UnescapeString(item.Title)
		article := Article{
			Title: title,
			URL:   item.Link,
		}
		articles = append(articles, article)
	}
	return articles
}

func GetNacioArticles(ctx context.Context, infoLog, errorLog *log.Logger) []Article {
	articles := scrapeNacioDigital(ctx, infoLog, errorLog)

	undesiredTopics := []string{"eleccions", "PSC", "Puigdemont", "Salvador Illa", "ERC", "amnistia", "Aragonès", "ANC", "PP", "PNB", "Koldo", "Illa", "BBVA"}
	articles = filterByTopics(articles, undesiredTopics)
	return articles
}
