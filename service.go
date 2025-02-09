package go_gaku

import (
	"encoding/json"
	"log"
	"os"
	"path/filepath"
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

func UpdateEps() {
	gogaku := ReadMainPrf()

	removeEps()
	
	for lang, prgs := range gogaku.Languages {
		langdir := filepath.Join(LANGS_DIR, lang)
		CheckDir(langdir)
		for _, prg := range prgs {
			eps_json := filepath.Join(langdir, prg.Title + ".json")
			if _, ok := os.Stat(eps_json); ok != nil {
				title, eps := GetEps(prg)
				log.Println(title)
				v, _ := json.Marshal(eps)
				err := os.WriteFile(eps_json, v, 0o644)
				if err != nil {
					log.Fatal(err)
				}
			}
		}
	}
}
// Local Variables:
// tab-width: 4
// End:
