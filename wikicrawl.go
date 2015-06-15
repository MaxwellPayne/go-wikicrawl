package wikicrawl

import (
  "fmt"
	"strings"
	"regexp"
	//"net/url"
	"encoding/json"
	"golang.org/x/net/html"
	"github.com/franela/goreq"
	"github.com/PuerkitoBio/goquery"
)

const ApiRoot string = "https://en.wikipedia.org/w/api.php"
const WikiContentRoot string = "https://en.wikipedia.org/wiki/"

type WikiParams struct {
	Format string
	Action string
	Page string
	Prop string
}

func randWikiUrl(output chan *WikiPage) {
	// Keep retrying until find non-disambiguation page
	for {
		redirectResponse, _ := goreq.Request{
			Uri: WikiContentRoot + "Special:Random",
		}.Do()
		// Obtain the random url from the redirect's location header
		redirectUrl, _ := redirectResponse.Location()
		redirectUrlString := redirectUrl.String()
		if !strings.Contains(strings.ToLower(redirectUrlString), "disambiguation") {
			randomTitle := strings.Replace(redirectUrlString, WikiContentRoot, "", 1)
			output <- NewWikiPage(randomTitle)
		}
	}
}

func Crawl(startPage *WikiPage, resultsChannel chan CrawlResult) {
	const MAX_JUMPS int = 1
	//results := make([]WikiPage, 0, MAX_JUMPS)
	//visited := make(map[WikiPage]bool)

	var currentPage *WikiPage = startPage
	
	for i := 0; i < MAX_JUMPS; i++ {
		// Go to the next page
		params := WikiParams{
			Action: "parse",
			Format: "json",
			Page: currentPage.Title(),
			Prop: "text",
		}
		
		req := goreq.Request{
			Uri: "https://en.wikipedia.org/w/api.php",
			QueryString: params,
		}

		// get the next page
		res, err := req.Do()

		if err == nil {
			if rawJson, err2 := res.Body.ToString(); err2 == nil {
				// parse the next page's json
				var jsonMap map[string]interface{}
				json.Unmarshal([]byte(rawJson), &jsonMap)
				// extract html
				var pageHtml string = jsonMap["parse"].(map[string]interface{})["text"].(map[string]interface{})["*"].(string)
				if page, err3 := html.Parse(strings.NewReader(pageHtml)); err3 == nil {
					// tokenize html into goquery
					document := goquery.NewDocumentFromNode(page)
					anchors := document.Find("p").Find("a")
					// search for anchors between <p> tags
					for i := range anchors.Nodes {
						if href, exists := anchors.Eq(i).Attr("href"); exists {
							fmt.Println("found an href:", href)
							contentLinkRegex := regexp.MustCompile("^/wiki/([^:]+$)")
							if submatches := contentLinkRegex.FindStringSubmatch(href); submatches != nil {
								extractedTitle := submatches[1]
								fmt.Println("found a matching href:", extractedTitle)
								// TODO: log changes here
								break
							}
						}
					}
				}
			}
		}
	}
	resultsChannel <- CrawlResult{[]WikiPage{*startPage}}
}

// Dummy function that allows main() to be run from other programs
func Run() {
	main()
}

func main() {

	wikipageChannel := make(chan *WikiPage, 1)
	go randWikiUrl(wikipageChannel)
	startPage := <- wikipageChannel
	fmt.Println(startPage)
	fmt.Println(startPage.Title())

	dummyStartPage := NewWikiPage("akihabara")
	resultsChan := make(chan CrawlResult)
	go Crawl(dummyStartPage, resultsChan)
	fmt.Println("Done!", <- resultsChan)
}
