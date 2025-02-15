package main

import (
	"encoding/json"
	"log"
	"os"
	"path/filepath"
	"strings"
	"net/http"
	"io"
)

func removeEps() {
	dirs, err := filepath.Glob(
		filepath.Join(LANGS_DIR, "*"))
	if err != nil {
		log.Fatal(err)
	}
	for _, dir := range dirs {
		log.Printf("Remove all files under %s dir.\n",  dir)
		os.RemoveAll(dir)
	}
}

func DownloadFile(filepath string, url string) error {

    resp, err := http.Get(url)
    if err != nil {
        return err
    }
    defer resp.Body.Close()

    out, err := os.Create(filepath)
    if err != nil {
        return err
    }
    defer out.Close()

    _, err = io.Copy(out, resp.Body)
    return err
}

func UpdateEps() {
	gogaku := ReadMainPrf()

	removeEps()
	
	for lang, prgs := range gogaku.Languages {
		langdir := filepath.Join(LANGS_DIR, lang)
		CheckDir(langdir)
		for _, prg := range prgs {
			if ! prg.RecordFlag {
				continue
			}
			eps_json := filepath.Join(langdir, prg.Title + ".json")
			_, eps := GetEps(prg)
			v, _ := json.Marshal(eps)
			err := os.WriteFile(eps_json, v, 0o644)
			if err != nil {
				log.Fatal(err)
			}
			
			if prg.CoverJPG == "" {
				continue
			}
			cover_elms := strings.Split(prg.CoverJPG, "/")
			img_path := filepath.Join(langdir, cover_elms[len(cover_elms)-1])
			DownloadFile(img_path, prg.CoverJPG)
		}
	}
}
// Local Variables:
// tab-width: 4
// End:
