package scraper

import (
	"net/url"
	"slices"
	"strings"
)

// filterByPaywall removes all articles that are behind a paywall, if the
// removePaywall boolean input parameter equals true.
func filterByPaywall(articles []Article, removePaywall bool) []Article {
	// Return whole article slice if paywall should not be removed.
	if !removePaywall {
		return articles
	}

	output := make([]Article, 0, 100)
	for _, article := range articles {
		if !article.Paywall {
			output = append(output, article)
		}
	}
	return output
}

func filterByTopics(articles []Article, undesiredTopics []string) []Article {
	output := make([]Article, 0, 100)

	for _, article := range articles {
		appendArticle := true
		// Create a single string with all articleTopics of the article.
		articleTopics := strings.Join(article.Topics, " ")
		articleTopics = strings.ToLower(articleTopics)

		for _, undesiredTopic := range undesiredTopics {
			undesiredTopic = strings.ToLower(undesiredTopic)
			if strings.Contains(articleTopics, undesiredTopic) {
				appendArticle = false
				break
			}
			// Check if the undesired topic is also found in the title of the
			// article, if so, do not append the article to the otuput.
			if strings.Contains(strings.ToLower(article.Title), undesiredTopic) {
				appendArticle = false
				break
			}
		}
		if appendArticle {
			output = append(output, article)
		}
	}
	return output
}

func filterByTopicsStrict(articles []Article, undesiredTopics []string) []Article {
	output := make([]Article, 0, 100)

	for _, article := range articles {
		if article.containsUndesiredTopic(undesiredTopics) {
			continue
		}
		output = append(output, article)
	}
	return output
}

func (article Article) containsUndesiredTopic(undesiredTopics []string) bool {
	for _, undesiredTopic := range undesiredTopics {
		if slices.Contains(article.Topics, undesiredTopic) {
			return true
		}
	}
	return false
}

func filterByUrlHostName(articles []Article, undesiredHostNames []string) []Article {
	output := make([]Article, 0, 100)

	for _, article := range articles {
		appendArticle := true

		// Parse URL of article.
		url, err := url.Parse(article.URL)
		if err != nil {
			continue
		}

		for _, undesiredHostName := range undesiredHostNames {
			if strings.Contains(url.Hostname(), undesiredHostName) {
				appendArticle = false
				break
			}
		}
		if appendArticle {
			output = append(output, article)
		}
	}
	return output
}

// isRelativeURL returns true if domain does not have a hostname defined.
func isRelativeURL(urlStr string) bool {
	parsedURL, err := url.Parse(urlStr)
	if err != nil {
		return false
	}

	return parsedURL.Hostname() == ""
}

func addSiteDomainIfRelativeURL(siteDomain, urlInput string) string {
	if isRelativeURL(urlInput) {
		parsedURL, err := url.Parse(urlInput)
		if err != nil {
			return ""
		}
		parsedDomain, err := url.Parse(siteDomain)
		if err != nil {
			return ""
		}
		parsedURL.Host = parsedDomain.Hostname()
		parsedURL.Scheme = parsedDomain.Scheme
		return parsedURL.String()
	}
	return urlInput
}

// fixMissingHostname adds the hostname of the siteDomain to all articles
// that only have a relative URL.
func fixMissingHostname(siteDomain string, articles []Article) []Article {
	fixedArticles := make([]Article, 0, len(articles))
	for _, article := range articles {
		article.URL = addSiteDomainIfRelativeURL(siteDomain, article.URL)
		fixedArticles = append(fixedArticles, article)
	}
	return fixedArticles
}
