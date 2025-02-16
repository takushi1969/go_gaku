package main

import (
	"path/filepath"
)

func update() {
	removeEpsJSON()
	
	gogaku := updateMainPrf("https://www.nhk.or.jp/gogaku/english", false)
	for lang, prgs := range gogaku.Languages {
		langdir := filepath.Join(LANGS_DIR, lang)
		createDir(langdir)
		
		for _, prg := range prgs {
			if prg.DlFlag == false {
				continue
			}
			prg.downloadEps(langdir)
		}
	}
}

// Local Variables:
// tab-width: 4
// End:
