package scraper

import (
	"context"
	"log"

	colly "github.com/gocolly/colly/v2"
)

func scrapeHackerNews(ctx context.Context, infoLog, errorLog *log.Logger) []Article {

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
		url:             "https://news.ycombinator.com/",
		querySelector:   "span.titleline",
		queryCallbackFn: f,
	}

	articles := scrape(ctx, infoLog, errorLog, q)

	return articles
}

func GetHackerNewsArticles(ctx context.Context, infoLog, errorLog *log.Logger) []Article {

	articles := scrapeHackerNews(ctx, infoLog, errorLog)
	return articles

}
