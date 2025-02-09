package go_gaku

import (
	"log"
	"net/http"
	"regexp"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

type Program struct {
	Title string
	SiteID string
	CornerID string
}

func GetPrgs(url string) []Program {
	const nhk_web_url = "https://www.nhk.or.jp/gogaku/"
	var prgs []Program

	if strings.Index(url, nhk_web_url) != 0 {
		log.Fatal("the argument should start with " + nhk_web_url)
	}
	
	res, err := http.Get(url)
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
		func(i int, prg *goquery.Selection) {
			title := prg.Find(".programtitle").First().Text()
			prg.Find("a").Each(
				func(i int, anchor *goquery.Selection) {
					href, exist := anchor.Attr("href")
					if exist {
						exp, _ := regexp.Compile(`radio/ondemand/detail\.html\?p=(.*?)_([^"]+)`)
						matched := exp.FindAllStringSubmatch(href, -1)
						if matched != nil {
							prgs = append(prgs,
								Program{
									title,
									matched[0][1],
									matched[0][2],
								})
						}
					}
			})

		})
		
	return prgs
}

// Local Variables:
// tab-width: 4
// End:
