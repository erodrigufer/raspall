package scraper

import (
	"context"
	"log"

	"github.com/gocolly/colly"
)

func ScrapeNacioDigital(ctx context.Context, infoLog, errorLog *log.Logger) []Article {

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
		url:             "https://www.naciodigital.cat/",
		querySelector:   "h2.titolnoticiallistat",
		queryCallbackFn: f,
	}

	articles := scrape(ctx, infoLog, errorLog, q)

	return articles
}
