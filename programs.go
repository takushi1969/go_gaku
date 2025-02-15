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

func getImgUrl(prg *goquery.Selection, lang_url string) string {
	img_url, exist := prg.Find(".thumbnail img").First().Attr("src")
	if exist {
		img_url, err := url.JoinPath(lang_url, img_url)
		if err == nil {
			log.Println(img_url)
			img_url = ""
		}
	}

	return img_url
}

func parsePrg(anchor *goquery.Selection, title, img_url string) *Program {
	var prg *Program
	
	href, exist := anchor.Attr("href")
	if ! exist {
		return prg
	}
	
	exp, _ := regexp.Compile(`radio/ondemand/detail\.html\?p=(.*?)_([^"]+)`)
	matched := exp.FindAllStringSubmatch(href, -1)
	if matched != nil {
		prg = new(Program)
		prg.Title = title
		prg.SiteID = matched[0][1]
		prg.CornerID = matched[0][2]
		prg.CoverJPG = img_url
	}

	return prg
}
	
func parsePrgs(doc *goquery.Document, lang_url string) []Program {
	var prgs []Program
	
	doc.Find("#listRadio .programbox").Each(
		func(i int, prg *goquery.Selection) {
			img_url := getImgUrl(prg, lang_url)
			title := prg.Find(".programtitle").First().Text()
			prg.Find("a").Each(
				func(i int, anchor *goquery.Selection) {
					prg := parsePrg(anchor, title, img_url)
					if prg != nil {
						prgs = append(prgs, *prg)
					}
				})
		})
	return prgs
}

func getPrgs(lang_url string) []Program {
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

	prgs = parsePrgs(doc, lang_url)

	return prgs
}

// Local Variables:
// tab-width: 4
// End:
