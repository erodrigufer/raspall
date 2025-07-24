package scraper

import (
	"context"
	"log/slog"

	colly "github.com/gocolly/colly/v2"
)

func scrapeZeit(ctx context.Context, infoLog, errorLog *slog.Logger) []Article {
	f := func(art *[]Article) colly.HTMLCallback {
		return func(element *colly.HTMLElement) {
			title := element.ChildText("span.zon-teaser__title")
			url := element.ChildAttr("a", "href")
			// Check if Z+ symbol is present, i.e. if it is a paywall article.
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

func GetZeitArticles(ctx context.Context, infoLog, errorLog *slog.Logger, removePaywall bool) []Article {
	articles := scrapeZeit(ctx, infoLog, errorLog)
	articles = filterByPaywall(articles, removePaywall)
	articles = filterByUrlHostName(articles, []string{"premium.zeit.de", "verlag.zeit.de", "sudoku.zeit.de", "spiele.zeit.de", "wiwo.de", "freundederzeit.typeform.com", "zeitakademie.de"})
	return articles
}
