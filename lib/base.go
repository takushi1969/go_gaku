package main

import (
	"fmt"
	"log"
	"net/http"
	"regexp"

	"github.com/PuerkitoBio/goquery"
)

func get_programs() {
	res, err := http.Get("https://www.nhk.or.jp/gogaku/english/")
	if err != nil {
		log.Fatal(err)
	}

	defer res.Body.Close()

	if res.StatusCode != 200 {
		log.Fatalf("status code error: %d %s", res.StatusCode, res.Status)
	}

	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		log.Fatal(err)
	}

	doc.Find("#listRadio .programbox").Each(
		func(i int, program *goquery.Selection) {
			title := program.Find(".programtitle").First().Text()
			href := ""
			program.Find("a").Each(
				func(i int, anchor *goquery.Selection) {
					var exist bool
					href, exist = anchor.Attr("href")
					if exist {
						matched, _ := regexp.Match(
							`radio/ondemand/detail\.html\?p=([0-9_]+)`,
							[]byte(href))
						if matched {
							fmt.Println(title, href)
						}
					}
				})
		})
}

func main() {
	get_programs()
}


// Local Variables:
// tab-width: 4
// End:
