package wikicrawl

import (
	"net/url"
	"strings"
)

type WikiPage struct {
  FormattedUrl *url.URL
}

func NewWikiPage(title string) *WikiPage {
	rawUrl := WikiContentRoot + title
	parsedUrl, _ := url.Parse(rawUrl)
	
	page := WikiPage{
		FormattedUrl: parsedUrl,
	}
	return &page
}



func (wikiPage *WikiPage) Title() string {
	title, _ := url.QueryUnescape(wikiPage.FormattedUrl.String())
	title = strings.Replace(title, WikiContentRoot, "", 1)
	title = strings.Replace(title, "_", " ", -1)
	return title
}
