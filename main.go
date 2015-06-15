package wikicrawl

import (
  "fmt"
	"strings"
	"net/url"
	"encoding/json"
	"io/ioutil"
	"golang.org/x/net/html"
	"github.com/franela/goreq"
	"github.com/PuerkitoBio/goquery"
)

const ApiRoot string = "https://en.wikipedia.org/w/api.php"


type WikiParams struct {
	Format string
	Action string
	Page string
	Prop string
}

func randWikiUrl(output chan *url.URL) {
	// Keep retrying until find non-disambiguation page
	for {
		redirectResponse, _ := goreq.Request{
			Uri: "https://en.wikipedia.org/wiki/Special:Random",
		}.Do()
		// Obtain the random url from the redirect's location header
		redirectUrl, _ := redirectResponse.Location()
		if !strings.Contains(strings.ToLower(redirectUrl.String()), "disambiguation") {
			output <- redirectUrl
		}
	}
}

func main() {

	urlChannel := make(chan *url.URL, 1)
	go randWikiUrl(urlChannel)
	startUrl := <- urlChannel
	fmt.Println("random is", startUrl)

	startPage := WikiPage{startUrl}
	fmt.Println(startPage)

	return

	params := WikiParams{
		Action: "parse",
		Format: "json",
		Page: "akihabara",
		Prop: "text",
	}
	
	
	
	req := goreq.Request{
		Uri: "https://en.wikipedia.org/w/api.php",
		QueryString: params,
	}

	res, err := req.Do()

	if err == nil {
		fmt.Println(res.StatusCode)
		if rawJson, err2 := res.Body.ToString(); err2 == nil {
			var jsonMap map[string]interface{}
			json.Unmarshal([]byte(rawJson), &jsonMap)
			var pageHtml string = jsonMap["parse"].(map[string]interface{})["text"].(map[string]interface{})["*"].(string)
			
			ioutil.WriteFile("out.html", []byte(pageHtml), 0666)
			if page, err3 := html.Parse(strings.NewReader(pageHtml)); err3 == nil {
				document := goquery.NewDocumentFromNode(page)
				anchors := document.Find("p").Find("a")
				for i := range anchors.Nodes {
					if href, exists := anchors.Eq(i).Attr("href"); exists {
						fmt.Println(href)
					}
					
				}
			}
		}
	}
}
