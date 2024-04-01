package scraper

import (
	"context"
	"log"

	"github.com/gocolly/colly"
)

func ScrapeZeit(ctx context.Context, infoLog, errorLog *log.Logger) []Article {

	f := func(art *[]Article) colly.HTMLCallback {
		return func(element *colly.HTMLElement) {
			title := element.ChildText("span.zon-teaser__title")
			url := element.ChildAttr("a", "href")
			// Check if Z+ symbol is present.
			paywallStr := element.ChildAttr("svg.zplus-logo", "aria-label")
			var paywall bool
			if paywallStr != "" {
				paywall = true
			}
			if title != "" && url != "" {
				article := Article{Title: title, URL: url, Paywall: paywall}
				*art = append(*art, article)
			}

		}
	}

	q := collectorQuery{
		url:             "https://www.zeit.de/",
		querySelector:   "div.zon-teaser__container",
		queryCallbackFn: f,
	}

	articles := scrape(ctx, infoLog, errorLog, q)

	return articles
}
