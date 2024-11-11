package scraper

import (
	"context"
	"fmt"
	"log"

	colly "github.com/gocolly/colly/v2"
)

func scrapeTheGuardian(ctx context.Context, infoLog, errorLog *log.Logger) []Article {
	f := func(art *[]Article) colly.HTMLCallback {
		return func(element *colly.HTMLElement) {
			title := element.ChildAttr("a.dcr-lv2v9o", "aria-label")
			url := element.ChildAttr("a.dcr-lv2v9o", "href")
			// URLs are relative paths in the webpage.
			url = fmt.Sprintf("https://www.theguardian.com%s", url)

			if title != "" && url != "" {
				article := Article{Title: title, URL: url}
				*art = append(*art, article)
			}
		}
	}

	q := collectorQuery{
		url:             "https://www.theguardian.com/europe",
		querySelector:   "div.dcr-f9aim1",
		queryCallbackFn: f,
	}

	articles := scrape(ctx, infoLog, errorLog, q)

	return articles
}

func GetTheGuardianArticles(ctx context.Context, infoLog, errorLog *log.Logger) []Article {
	articles := scrapeTheGuardian(ctx, infoLog, errorLog)

	return articles
}
