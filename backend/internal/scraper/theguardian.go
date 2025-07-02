package scraper

import (
	"context"
	"html"
	"log"

	"github.com/mmcdole/gofeed"
)

const theGuardianRSSFeed = "https://theguardian.com/europe/rss"

func scrapeTheGuardian(ctx context.Context, infoLog, errorLog *log.Logger) []Article {
	infoLog.Printf("Visiting: %s", theGuardianRSSFeed)
	fp := gofeed.NewParser()
	feed, err := fp.ParseURLWithContext(theGuardianRSSFeed, ctx)
	if err != nil {
		errorLog.Printf("an error occurred while scraping The Guardian: %s", err.Error())
		return []Article{}
	}
	articles := make([]Article, 0)
	for _, item := range feed.Items {
		title := html.UnescapeString(item.Title)
		article := Article{
			Title:  title,
			URL:    item.Link,
			Topics: item.Categories[:1],
		}

		articles = append(articles, article)
	}
	return articles
}

func GetTheGuardianArticles(ctx context.Context, infoLog, errorLog *log.Logger) []Article {
	articles := scrapeTheGuardian(ctx, infoLog, errorLog)

	undesiredTopics := []string{"Trump", "Musk", "Gaza", "Israel", "Lebanon", "Starmer"}
	articles = filterByTopics(articles, undesiredTopics)
	return articles
}
