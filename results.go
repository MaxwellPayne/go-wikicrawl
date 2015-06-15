package wikicrawl

type CrawlResult struct {
	Trail []WikiPage
}

func (crawlResult *CrawlResult) StartPage() WikiPage {
	return crawlResult.Trail[0]
}
