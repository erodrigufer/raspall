package scraper

import (
	"context"
	"log"

	colly "github.com/gocolly/colly/v2"
)

const LOBSTERS_URL = "https://lobste.rs/"

func scrapeLobsters(ctx context.Context, infoLog, errorLog *log.Logger) []Article {
	f := func(art *[]Article) colly.HTMLCallback {
		return func(element *colly.HTMLElement) {
			title := element.ChildText("span.link > a")
			url := element.ChildAttr("span.link > a", "href")
			tags := element.ChildTexts("span.tags > a.tag")
			if title != "" && url != "" {
				article := Article{Title: title, URL: url, Topics: tags}
				*art = append(*art, article)
			}
		}
	}

	q := collectorQuery{
		url:             LOBSTERS_URL,
		querySelector:   "div.details",
		queryCallbackFn: f,
	}

	articles := scrape(ctx, infoLog, errorLog, q)
	articles = fixMissingHostname(LOBSTERS_URL, articles)

	return articles
}

func GetLobstersArticles(ctx context.Context, infoLog, errorLog *log.Logger) []Article {
	articles := scrapeLobsters(ctx, infoLog, errorLog)

	undesiredTopics := []string{"rust", "graphics", "math", "science", "slides", "apl", "c", "c++", "dotnet", "fortran", "java", "haskell", "kotlin", "ml", "objectivec", "php", "python", "swift", "windows", "retrocomputing"}
	articles = filterByTopicsStrict(articles, undesiredTopics)
	return articles
}
