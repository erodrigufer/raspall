package scraper

import (
	"context"
	"log/slog"

	colly "github.com/gocolly/colly/v2"
)

const HN_URL = "https://news.ycombinator.com/"

func scrapeHackerNews(ctx context.Context, infoLog, errorLog *slog.Logger) []Article {
	f := func(art *[]Article) colly.HTMLCallback {
		return func(element *colly.HTMLElement) {
			titleTexts := element.ChildTexts("a")
			title := titleTexts[0]
			url := element.ChildAttr("a", "href")

			if title != "" && url != "" {
				article := Article{Title: title, URL: url}
				*art = append(*art, article)
			}
		}
	}

	q := collectorQuery{
		url:             HN_URL,
		querySelector:   "span.titleline",
		queryCallbackFn: f,
	}

	articles := scrape(ctx, infoLog, errorLog, q)
	articles = fixMissingHostname(HN_URL, articles)

	return articles
}

func GetHackerNewsArticles(ctx context.Context, infoLog, errorLog *slog.Logger) []Article {
	articles := scrapeHackerNews(ctx, infoLog, errorLog)

	undesiredTopics := []string{"hiring"}
	articles = filterByTopics(articles, undesiredTopics)
	return articles
}
