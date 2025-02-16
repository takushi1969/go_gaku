package main

import (
	"errors"
	"fmt"
	"io"
	"io/fs"
	"log"
	"net/http"
	"os"
	"regexp"
	"strings"
	"time"

	"encoding/json"
	"path/filepath"

	"github.com/0xAX/notificator"
	"github.com/bogem/id3v2"
	"github.com/tidwall/gjson"
	ffmpeg "github.com/u2takey/ffmpeg-go"
)	

type EpisodeInfo struct {
	Title string
	HLS string
	Date string
}

const eps_url = "https://www.nhk.or.jp/radio-api/app/v1/web/ondemand/series?site_id=%s&corner_site_id=%s"

var epsValidDuration time.Duration

func init() {
	epsValidDuration = time.Duration(24*60*time.Minute)
}
	
func removeEpsJSON() {
	filepath.WalkDir(
		LANGS_DIR,
		func (path string, d fs.DirEntry, err error) error {
			match, _ := regexp.MatchString(`(?i)\.json$`, path)
			if ! match {
				return nil
			}
			dinfo, err := d.Info()
			if err != nil {
				os.Remove(path)
				return nil
			}
			if time.Now().Sub(dinfo.ModTime()) < epsValidDuration {
				return nil
			}
			os.Remove(path)
			
			return nil
		})
}

func (prgInf ProgramInfo)downloadCover(profLangDir string) string {

	token := strings.Split(prgInf.CoverImg, "/")
	imgPath := filepath.Join(profLangDir, token[len(token)-1])
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

func (prgInf ProgramInfo)getJSON() string {
	res, err := http.Get(fmt.Sprintf(eps_url, prgInf.SiteID, prgInf.CornerID))
	if err != nil {
		log.Fatal(err)
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		log.Fatal(err)
	}
	
	return string(body)
}

func (prgInf ProgramInfo)getEps() ([]EpisodeInfo) {
	var eps []EpisodeInfo

	content := prgInf.getJSON()

	result := gjson.Get(content, "episodes")
	result.ForEach(func(key, value gjson.Result) bool {
		date :=	(strings.Split(
			(strings.Split(
				value.Get("aa_contents_id").String(), ";"))[4], "T"))[0]
		eps = append(eps,
			EpisodeInfo{
				value.Get("program_title").String(),
				value.Get("stream_url").String(),
				date,
			})
		return true
	})

	return eps
}

func (prgInf ProgramInfo)updateEps(profLangDir string) []EpisodeInfo {
	var epsInf []EpisodeInfo
	
	createDir(profLangDir)

	epsPath := filepath.Join(profLangDir, prgInf.Title + ".json")
	if _, err := os.Stat(epsPath); err != nil {
		if errors.Is(err, fs.ErrNotExist) {
			eps := prgInf.getEps()
			v, _ := json.Marshal(eps)
			err := os.WriteFile(epsPath, v, 0o644)
			if err != nil {
				log.Fatal(err)
			}
		}
	}

	val, err := os.ReadFile(epsPath)
	if err != nil {
		log.Fatal(err)
	}

	err = json.Unmarshal(val, &epsInf)
	if err != nil {
		log.Fatal(err)
	}

	return epsInf
}

func (epInf EpisodeInfo)mp3file() string {
	output := epInf.Date + "_" + epInf.Title + ".mp3"

	return output
}

func (epInf EpisodeInfo)notify(cover string) {
	notify := notificator.New(notificator.Options{
		DefaultIcon: cover,
		AppName: "Gogaku Downloader",
	})
	notify.Push(
		"saved.",
		epInf.mp3file(),
		"icon.png",
		notificator.UR_NORMAL)
}

func setCover(tag *id3v2.Tag, cover string) {

	var mimeType string
	
	switch filepath.Ext(cover) {
	case ".jpg":
		mimeType = "image/jpeg"
	case ".png":
		mimeType = "image/png"
	default:
		return
	}

	artwork, err:= os.ReadFile(cover)
	if err != nil {
		log.Fatal(err)
	}

	pic := id3v2.PictureFrame{
		Encoding: id3v2.EncodingUTF8,
		MimeType: mimeType,
		PictureType: id3v2.PTFrontCover,
		Description: "Front cover",
		Picture: artwork,
	}
	tag.AddAttachedPicture(pic)
}

func (epInf EpisodeInfo)writeTag(mp3file, album, cover string) {
	// Open file and parse tag in it.
	tag, err := id3v2.Open(mp3file, id3v2.Options{Parse: true})
	if err != nil {
	log.Fatal("Error while opening mp3 file: ", err)
	}
	defer tag.Close()

	tags := make(map[string]string)
	tags["Artist"] = "NHK"
	tags["Album/Movie/Show title"] = album
	tags["Title"] = epInf.Title
	tags["Genre"] = "Other"
	for k, v := range tags {
		tag.AddTextFrame(id3v2.V23CommonIDs[k],	id3v2.EncodingUTF8, v)
	}
	setCover(tag, cover)

	
	// Write tag to file.
	if err = tag.Save(); err != nil {
		log.Fatal("Error while saving a tag: ", err)
	}
}

func (epInf EpisodeInfo)download(album, cover string) {
	output := epInf.mp3file()
	_, err := os.Stat(output) 
	if err == nil  {
		return
	}

	tmpMp3, _ := os.CreateTemp("", "gogaku_*.mp3")
	tmpMp3.Close()
		
	in_args := ffmpeg.KwArgs{"http_seekable": 0}
	out_args := ffmpeg.KwArgs{
		"c:a": "libmp3lame",
		"b:a": "128k",
	}
	err = ffmpeg.Input(epInf.HLS, in_args).
		Output(tmpMp3.Name(), out_args).
		OverWriteOutput().Run()
	if err != nil {
		log.Println(fmt.Sprintf("fail to download %s", epInf.HLS))
		if _, err = os.Stat(tmpMp3.Name()); err == nil {
			os.Remove(tmpMp3.Name())
		}
		return
	}
	epInf.writeTag(tmpMp3.Name(), album, cover)
	os.Rename(tmpMp3.Name(), output)
}

func (prgInf ProgramInfo)downloadEps(profLangDir string) {
	coverImg := prgInf.downloadCover(profLangDir)
	epsInf := prgInf.updateEps(profLangDir)
	for _, epinf := range epsInf {
		epinf.download(prgInf.Title, coverImg)
	}
}

// Local Variables:
// tab-width: 4
// End:
