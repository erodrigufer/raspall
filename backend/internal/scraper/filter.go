package scraper

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
