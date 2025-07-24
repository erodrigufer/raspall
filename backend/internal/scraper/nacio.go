package scraper

import (
	"context"
	"html"
	"log/slog"

	"github.com/mmcdole/gofeed"
)

const nacioDigitalrssFeed = "https://naciodigital.cat/rss"

func scrapeNacioDigital(ctx context.Context, infoLog, errorLog *slog.Logger) []Article {
	infoLog.Info("an RSS feed is being scraped", slog.String("url", nacioDigitalrssFeed))
	fp := gofeed.NewParser()
	feed, err := fp.ParseURLWithContext(nacioDigitalrssFeed, ctx)
	if err != nil {
		errorLog.Error("an error occurred while scraping", slog.String("url", nacioDigitalrssFeed), slog.String("error_message", err.Error()))
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

func GetNacioArticles(ctx context.Context, infoLog, errorLog *slog.Logger) []Article {
	articles := scrapeNacioDigital(ctx, infoLog, errorLog)

	undesiredTopics := []string{"eleccions", "PSC", "Puigdemont", "Salvador Illa", "ERC", "amnistia", "Aragonès", "ANC", "PP", "PNB", "Koldo", "Illa", "BBVA"}
	articles = filterByTopics(articles, undesiredTopics)
	return articles
}
