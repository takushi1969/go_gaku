package main

import (
	"log"
	"net/url"
	"net/http"
	"regexp"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

type Program struct {
	Title string
	SiteID string
	CornerID string
	RecordFlag bool
	CoverJPG string
}

func GetPrgs(lang_url string) []Program {
	const nhk_web_url = "https://www.nhk.or.jp/gogaku/"
	var prgs []Program

	if strings.Index(lang_url, nhk_web_url) != 0 {
		log.Fatal("the argument should start with " + nhk_web_url)
	}
	
	res, err := http.Get(lang_url)
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
			img_path, exist := prg.Find(".thumbnail img").First().Attr("src")
			img_url := ""
			if exist {
				if img_url, err = url.JoinPath(lang_url, img_path); err == nil {
					log.Println(img_url)
				}
			}
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
									Title: title,
									SiteID: matched[0][1],
									CornerID: matched[0][2],
									CoverJPG: img_url,
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
