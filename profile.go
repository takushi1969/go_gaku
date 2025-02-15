package main

import (
	"encoding/json"
	"errors"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"time"
)

type Gogaku struct {
	Version float64
	Update time.Time
	Languages map[string][]Program
}

const PRF_VERSION = 1.0
var PRF_DIR string
var MAIN_PRF string
var LANGS_DIR string

func init() {
	PRF_DIR = filepath.Join(os.Getenv("HOME"),  ".go_gaku")
	MAIN_PRF = filepath.Join(PRF_DIR, "programs.json")
	LANGS_DIR = filepath.Join(PRF_DIR, "languages")

	for _, dir := range []string{PRF_DIR, LANGS_DIR} {
		CheckDir(dir)
	}
}


func ReadMainPrf() (*Gogaku) {
	var gogaku Gogaku

	val, err := os.ReadFile(MAIN_PRF)
	if err != nil {
		if errors.Is(err, fs.ErrNotExist) {
			return nil
		} else {
			log.Fatal(err)
		}
	}

	err = json.Unmarshal(val, &gogaku)
	if err != nil {
		log.Fatal(err)
	}

	if gogaku.Version != PRF_VERSION {
		log.Fatal("Unsupported Version")
	}

	return &gogaku
}

func WriteMainPrf(url string, force bool) {
	gogaku := ReadMainPrf()

	if force == false && gogaku != nil {
		if time.Now().Sub(gogaku.Update) < time.Duration(7 * 24 * time.Hour) {
			return
		}
	}
	
	re := regexp.MustCompile(`.*/(.*)`)
	lang := (re.FindStringSubmatch(url))[1]

	prgs := GetPrgs(url)
	if gogaku != nil {
		for _, old_prg := range gogaku.Languages[lang] {
			if old_prg.RecordFlag == false {
				continue
			}
			for n, new_prg := range prgs {
				if old_prg.Title == new_prg.Title {
					log.Println(n, old_prg, new_prg)
					prgs[n].RecordFlag = old_prg.RecordFlag
					break
				}
			}
		}
	} else {
		log.Println("route check")
		gogaku = &Gogaku{}
		gogaku.Languages = make(map[string][]Program)
	}

	log.Println(prgs[6])
	gogaku.Version = PRF_VERSION
	gogaku.Update = time.Now()
	gogaku.Languages[lang] = prgs

	v, _ := json.Marshal(gogaku)
	
	err := os.WriteFile(MAIN_PRF, v, 0o644)
	if err != nil {
	 	log.Fatal(err)
	}
}

// Local Variables:
// tab-width: 4
// End:
