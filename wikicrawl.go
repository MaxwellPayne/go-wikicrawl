package wikicrawl

import (
  "fmt"
	"strings"
	//"net/url"
	"encoding/json"
	"io/ioutil"
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

func Run() {
	main()
}

func main() {

	wikipageChannel := make(chan *WikiPage, 1)
	go randWikiUrl(wikipageChannel)
	startPage := <- wikipageChannel
	fmt.Println(startPage)
	fmt.Println(startPage.Title())

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
