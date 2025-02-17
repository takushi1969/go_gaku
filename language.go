package main

import (
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

type ProgramInfo struct {
	EpInfDir string
	EpInfFile string
	Title string
	SiteID string
	CornerID string
	DlFlag bool
	DlDir string
	CoverImg string
}

func getLang(langUrl string) string {
	// "https://www.nhk.or.jp/gogaku/english"
	token := strings.Split(langUrl, "/")
	n := 1
	if token[len(token)-n] == "" {
		n = 2
	} 
	return token[len(token)-n]
}

func getDlDir(lang, title string) string {
	return filepath.Join(DL_DIR, lang, title)
}

func getEpInfDir(lang string) string {
	return filepath.Join(LANGS_DIR, lang)
}

func getEpInfFile(lang, title string) string {
	return filepath.Join(getEpInfDir(lang), title+".json")
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

func parsePrg(anchor *goquery.Selection, langUrl, title, imgUrl string) *ProgramInfo {
	var prg *ProgramInfo
	
	href, exist := anchor.Attr("href")
	if ! exist {
		return prg
	}

	lang:= getLang(langUrl)
	
	exp, _ := regexp.Compile(`radio/ondemand/detail\.html\?p=(.*?)_([^"]+)`)
	matched := exp.FindAllStringSubmatch(href, -1)
	if matched != nil {
		prg = new(ProgramInfo)
		prg.EpInfDir = getEpInfDir(lang)
		prg.EpInfFile = getEpInfFile(lang, title)
		prg.Title = title
		prg.SiteID = matched[0][1]
		prg.CornerID = matched[0][2]
		prg.CoverImg = imgUrl
		prg.DlDir = getDlDir(lang, title)
	}

	return prg
}
	
func parsePrgs(doc *goquery.Document, langUrl string) []ProgramInfo {
	var prgs []ProgramInfo
	
	doc.Find("#listRadio .programbox").Each(
		func(i int, prg *goquery.Selection) {
			title := prg.Find(".programtitle").First().Text()
			imgUrl := getImgUrl(prg, langUrl)
			prg.Find("a").Each(
				func(i int, anchor *goquery.Selection) {
					prg := parsePrg(anchor, langUrl, title, imgUrl)
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

func (prgInf ProgramInfo)downloadCover() string {
	token := strings.Split(prgInf.CoverImg, "/")
	imgPath := filepath.Join(prgInf.EpInfDir, token[len(token)-1])
	if _, err := os.Stat(imgPath); err == nil {
		return imgPath
	}
		
    resp, err := http.Get(prgInf.CoverImg)
    if err != nil {
		log.Println(err)
		return ""
    }
    defer resp.Body.Close()

    out, err := os.Create(imgPath)
    if err != nil {
		log.Println(err)
		return ""
    }
    defer out.Close()

    _, err = io.Copy(out, resp.Body)
    if err != nil {
		log.Println(err)
		return ""
    }

	return imgPath
}

// Local Variables:
// tab-width: 4
// End:
