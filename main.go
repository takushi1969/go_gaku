package main

import (
	"encoding/json"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/0xAX/notificator"
	ffmpeg "github.com/u2takey/ffmpeg-go"
	"github.com/bogem/id3v2"
)

func tagging(mp3file, album, title string) {
	log.Println(mp3file, album, title)
	
	// Open file and parse tag in it.
	tag, err := id3v2.Open(mp3file, id3v2.Options{Parse: true})
	if err != nil {
	log.Fatal("Error while opening mp3 file: ", err)
	}
	defer tag.Close()

	tag.SetDefaultEncoding(id3v2.EncodingUTF8)
	
	// Set simple text frames
	tag.AddTextFrame(
		id3v2.V23CommonIDs["Artist"],
		id3v2.EncodingUTF8,
		"NHK")
	tag.AddTextFrame(
		id3v2.V23CommonIDs["Album/Movie/Show title"],
		id3v2.EncodingUTF8,
		album)
	tag.AddTextFrame(
		id3v2.V23CommonIDs["Title"],
		id3v2.EncodingUTF8,
		title)

	// Write tag to file.
	if err = tag.Save(); err != nil {
		log.Fatal("Error while saving a tag: ", err)
	}
}

func callFfmpeg(path string) {
	var eps []Episode

	val, err := os.ReadFile(path)
	if err != nil {
		log.Fatal(err)
	}

	err = json.Unmarshal(val, &eps)
	if err != nil {
		log.Fatal(err)
	}
	
	for _, ep := range eps {
		output := ep.Date + "_" + ep.Title + ".mp3"

		if _, err = os.Stat(output); err == nil  {
			continue
		}

		tmpMp3, _ := os.CreateTemp("", "gogaku_*.mp3")
		tmpMp3.Close()
		
		in_args := ffmpeg.KwArgs{"http_seekable": 0}
		out_args := ffmpeg.KwArgs{
			"c:a": "libmp3lame",
			// "map": "a",
			"b:a": "128k",
			"id3v2_version": "0",
		}
		err = ffmpeg.Input(ep.HLS, in_args).
			Output(tmpMp3.Name(), out_args).
			OverWriteOutput().Run()
		if err != nil {
			log.Println(err)
			if _, err = os.Stat(tmpMp3.Name()); err == nil {
				os.Remove(tmpMp3.Name())
			}
			continue
		}
		os.Rename(tmpMp3.Name(), output)
		// https://askubuntu.com/questions/248811/how-can-i-fix-incorrect-mp3-duration
		// ffmpeg.Input(tmpMp3.Name()).
		// 	Output(output, ffmpeg.KwArgs{"acodec": "copy"}).
		// 	OverWriteOutput().Run()
		// os.Remove(tmpMp3.Name())
		tagging(
			output,
			(strings.Split(filepath.Base(path), "."))[0],
			ep.Date + "_" + ep.Title)
		// os.Rename(tmpMp3.Name(), output)

		notify := notificator.New(notificator.Options{
			DefaultIcon: "icon.png",
			AppName: "Gogaku Downloader",
		})
		notify.Push(
			"saved.",
			output,
			"icon.png",
			notificator.UR_NORMAL)

	}
}

func record() {
	filepath.WalkDir(LANGS_DIR,
		func (path string, d fs.DirEntry, err error) error {
			if err != nil {
				log.Fatal(err)
			}

			if strings.Index(path, ".json") > 0 {
				callFfmpeg(path)
			}
			return nil
		})
}

func main() {
	//WriteMainPrf("https://www.nhk.or.jp/gogaku/english", true)
	//WriteMainPrf("https://www.nhk.or.jp/gogaku/chinese", true)
	// gogaku := ReadMainPrf()
	
	// UpdateEps()
	record()
	
}

// Local Variables:
// tab-width: 4
// End:
