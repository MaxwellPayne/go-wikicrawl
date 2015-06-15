package wikicrawl

import (
	"net/url"
)

type WikiPage struct {
  FormattedUrl *url.URL
}

func NewWikiPage(title string) *WikiPage {
	rawUrl := ApiRoot
	parsedUrl, _ := url.Parse(rawUrl)
	
	page := WikiPage{
		FormattedUrl: parsedUrl,
	}
	return &page
}

func NormalizeTitle(title string) string {
	return "hi"
}
