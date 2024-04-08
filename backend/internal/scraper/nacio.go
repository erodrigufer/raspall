package scraper

import (
	"context"
	"log"

	"github.com/gocolly/colly"
)

func scrapeNacioDigital(ctx context.Context, infoLog, errorLog *log.Logger) []Article {

	f := func(art *[]Article) colly.HTMLCallback {
		return func(element *colly.HTMLElement) {
			title := element.ChildText("h2 > a")
			url := element.ChildAttr("h2.titolnoticiallistat > a", "href")
			topic := element.ChildText("div.avantitol > span > a")
			// Sometimes the topic is not an 'a' element, then re-check
			// if a topic can be found with another query.
			if topic == "" {
				topic = element.ChildText("div.avantitol > span")
			}
			topics := make([]string, 0, 3)
			if title != "" && url != "" {
				if topic != "" {
					topics = append(topics, topic)
				}
				article := Article{Title: title, URL: url, Topics: topics}
				*art = append(*art, article)
			}
		}
	}

	q := collectorQuery{
		url:             "https://www.naciodigital.cat/",
		querySelector:   "div.noticia",
		queryCallbackFn: f,
	}

	articles := scrape(ctx, infoLog, errorLog, q)

	return articles
}

func GetNacioArticles(ctx context.Context, infoLog, errorLog *log.Logger) []Article {
	articles := scrapeNacioDigital(ctx, infoLog, errorLog)

	undesiredTopics := []string{"eleccions", "PSC", "Puigdemont", "Salvador Illa", "ERC", "amnistia", "Aragonès"}
	articles = filterByTopics(articles, undesiredTopics)
	return articles
}
