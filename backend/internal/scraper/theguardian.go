package scraper

import (
	"context"
	"html"
	"log/slog"

	"github.com/mmcdole/gofeed"
)

const theGuardianRSSFeed = "https://theguardian.com/europe/rss"

func scrapeTheGuardian(ctx context.Context, infoLog, errorLog *slog.Logger) []Article {
	infoLog.Info("an RSS feed is being scraped", slog.String("url", theGuardianRSSFeed))
	fp := gofeed.NewParser()
	feed, err := fp.ParseURLWithContext(theGuardianRSSFeed, ctx)
	if err != nil {
		errorLog.Error("an error occurred while scraping", slog.String("url", theGuardianRSSFeed), slog.String("error_message", err.Error()))
		return []Article{}
	}
	articles := make([]Article, 0)
	for _, item := range feed.Items {
		title := html.UnescapeString(item.Title)
		topics := make([]string, 0)
		if len(item.Categories) >= 1 {
			topics = append(topics, item.Categories[0])
		}
		article := Article{
			Title:  title,
			URL:    item.Link,
			Topics: topics,
		}

		articles = append(articles, article)
	}
	return articles
}

func GetTheGuardianArticles(ctx context.Context, infoLog, errorLog *slog.Logger) []Article {
	articles := scrapeTheGuardian(ctx, infoLog, errorLog)

	undesiredTopics := []string{"Trump", "Musk", "Gaza", "Israel", "Lebanon", "Starmer"}
	articles = filterByTopics(articles, undesiredTopics)
	return articles
}
