package scraper

import (
	"context"
	"log"

	colly "github.com/gocolly/colly/v2"
)

func scrapeLobsters(ctx context.Context, infoLog, errorLog *log.Logger) []Article {
	f := func(art *[]Article) colly.HTMLCallback {
		return func(element *colly.HTMLElement) {
			title := element.ChildText("a")
			url := element.ChildAttr("a", "href")

			if title != "" && url != "" {
				article := Article{Title: title, URL: url}
				*art = append(*art, article)
			}
		}
	}

	q := collectorQuery{
		url:             "https://lobste.rs/",
		querySelector:   "span.link",
		queryCallbackFn: f,
	}

	articles := scrape(ctx, infoLog, errorLog, q)

	return articles
}

func GetLobstersArticles(ctx context.Context, infoLog, errorLog *log.Logger) []Article {
	articles := scrapeLobsters(ctx, infoLog, errorLog)
	return articles
}
