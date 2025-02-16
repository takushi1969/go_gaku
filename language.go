package main

import (
	"log"
	"net/url"
	"net/http"
	"regexp"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

type ProgramInfo struct {
	Title string
	SiteID string
	CornerID string
	DlFlag bool
	CoverImg string
}

func getImgUrl(prg *goquery.Selection, langUrl string) string {
	imgUrl, exist := prg.Find(".thumbnail img").First().Attr("src")
	if exist {
		var err error
		imgUrl, err = url.JoinPath(langUrl, imgUrl)
		if err != nil {
			log.Println(imgUrl)
			imgUrl = ""
		}
	}

	return imgUrl
}

func parsePrg(anchor *goquery.Selection, title, imgUrl string) *ProgramInfo {
	var prg *ProgramInfo
	
	href, exist := anchor.Attr("href")
	if ! exist {
		return prg
	}
	
	exp, _ := regexp.Compile(`radio/ondemand/detail\.html\?p=(.*?)_([^"]+)`)
	matched := exp.FindAllStringSubmatch(href, -1)
	if matched != nil {
		prg = new(ProgramInfo)
		prg.Title = title
		prg.SiteID = matched[0][1]
		prg.CornerID = matched[0][2]
		prg.CoverImg = imgUrl
	}

	return prg
}
	
func parsePrgs(doc *goquery.Document, langUrl string) []ProgramInfo {
	var prgs []ProgramInfo
	
	doc.Find("#listRadio .programbox").Each(
		func(i int, prg *goquery.Selection) {
			imgUrl := getImgUrl(prg, langUrl)
			title := prg.Find(".programtitle").First().Text()
			prg.Find("a").Each(
				func(i int, anchor *goquery.Selection) {
					prg := parsePrg(anchor, title, imgUrl)
					if prg != nil {
						prgs = append(prgs, *prg)
					}
				})
		})
	return prgs
}

func getPrgs(langUrl string) []ProgramInfo {
	const nhk_web_url = "https://www.nhk.or.jp/gogaku/"
	var prgs []ProgramInfo

	if strings.Index(langUrl, nhk_web_url) != 0 {
		log.Fatal("the argument should start with " + nhk_web_url)
	}
	
	res, err := http.Get(langUrl)
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

	prgs = parsePrgs(doc, langUrl)

	return prgs
}

// Local Variables:
// tab-width: 4
// End:
