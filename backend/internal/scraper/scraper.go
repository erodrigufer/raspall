package scraper

import (
	"context"
	"fmt"

	"github.com/gocolly/colly"
)

// TODO: use type.
type Article struct {
	Title string
	URL   string
}

func Scrape(ctx context.Context) []Article {
	collector := colly.NewCollector()

	url := "https://www.naciodigital.cat/"

	collector.OnRequest(func(r *colly.Request) {
		// TODO: remove fmt.Print functions.
		// print the url of that request
		fmt.Println("Visiting", r.URL)
	})
	collector.OnResponse(func(r *colly.Response) {
		fmt.Println("Got a response from", r.Request.URL)
	})
	collector.OnError(func(r *colly.Response, e error) {
		fmt.Println("Blimey, an error occurred!:", e)
	})
	// n := 1
	articles := make([]Article, 0, 200)
	collector.OnHTML("h2.titolnoticiallistat", func(element *colly.HTMLElement) {
		title := element.ChildText("a")
		url := element.ChildAttr("a", "href")
		if title != "" && url != "" {
			article := Article{Title: title, URL: url}
			articles = append(articles, article)
		}

		// fmt.Printf("[%d].\t %s | %s \n", n, title, url)
		// n = n + 1
	})
	collector.Visit(url)

	return articles
}
