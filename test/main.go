package main

import (
	"encoding/json"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"strings"

	go_gaku "github.com/takushi1969/go_gogaku"
	ffmpeg "github.com/u2takey/ffmpeg-go"
)

func callFfmpeg(path string) {
	var eps []go_gaku.Episode

	val, err := os.ReadFile(path)
	if err != nil {
		log.Fatal(err)
	}

	err = json.Unmarshal(val, &eps)
	if err != nil {
		log.Fatal(err)
	}
	
	for _, ep := range eps {
		in_args := ffmpeg.KwArgs{"http_seekable": 0}
		
		output := ep.Date + "_" + ep.Title + ".mp3"
		out_args := ffmpeg.KwArgs{
			"c:a": "libmp3lame",
			"map": "a",
			"b": "128k",
		}
		ffmpeg.Input(ep.HLS, in_args).
			Output(output, out_args).
			OverWriteOutput().ErrorToStdOut().Run()		
	}
}

func record() {
	filepath.WalkDir(go_gaku.LANGS_DIR,
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
	
	//UpdateEps()
	record()
	
}

// Local Variables:
// tab-width: 4
// End:
