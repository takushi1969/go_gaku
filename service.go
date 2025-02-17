package main

import (
	"path/filepath"
)

func update() {
	removeEpsJSON()
	
	gogaku := updateMainPrf("https://www.nhk.or.jp/gogaku/english", false)
	for lang, prgs := range gogaku.Languages {
		for _, dir := range []string{LANGS_DIR, DL_DIR} {
			createDir(filepath.Join(dir, lang))
		}
		
		for _, prg := range prgs {
			if prg.DlFlag == false {
				continue
			}
			prg.downloadEps()
		}
	}
}


// Local Variables:
// tab-width: 4
// End:
