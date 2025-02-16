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
	UpdateTime time.Time
	Languages map[string][]ProgramInfo
}

const PRF_VERSION = 1.0
var PRF_DIR string
var MAIN_PRF string
var LANGS_DIR string
var DL_DIR string

func init() {
	PRF_DIR = filepath.Join(os.Getenv("HOME"),  ".go_gaku")
	MAIN_PRF = filepath.Join(PRF_DIR, "programs.json")
	LANGS_DIR = filepath.Join(PRF_DIR, "languages")
	DL_DIR = filepath.Join(os.Getenv("HOME"),  "go_gaku")

	for _, dir := range []string{PRF_DIR, LANGS_DIR, DL_DIR} {
		createDir(dir)
	}
}


func readMainPrf() *Gogaku {
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
		return nil
	}

	return &gogaku
}

func updateMainPrf(gogakuUrl string, force bool) *Gogaku {
	gogaku := readMainPrf()

	if force == false&&
		gogaku != nil &&
		time.Now().Sub(gogaku.UpdateTime) < time.Duration(7 * 24 * time.Hour) {
		return gogaku
	}
	
	re := regexp.MustCompile(`.*/(.*)`)
	lang := (re.FindStringSubmatch(gogakuUrl))[1]

	prgs := getPrgs(gogakuUrl)
	if gogaku != nil {
		for _, oldPrg := range gogaku.Languages[lang] {
			if oldPrg.DlFlag == false {
				continue
			}
			for n, newPrg := range prgs {
				if oldPrg.Title == newPrg.Title {
					prgs[n].DlFlag = oldPrg.DlFlag
					break
				}
			}
		}
	} else {
		gogaku = new(Gogaku)
		gogaku.Languages = make(map[string][]ProgramInfo)
	}

	gogaku.Version = PRF_VERSION
	gogaku.UpdateTime = time.Now()
	gogaku.Languages[lang] = prgs
	v, err := json.Marshal(gogaku)
	if err != nil {
		log.Panic(err)
	}
	err = os.WriteFile(MAIN_PRF, v, 0o644)
	if err != nil {
	 	log.Fatal(err)
	}

	return gogaku
}

// Local Variables:
// tab-width: 4
// End:
