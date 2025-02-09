package go_gaku

import (
	"errors"
	"io/fs"
	"log"
	"os"
)

func CheckDir(dir string) { 
	if _, err := os.Stat(dir); err != nil {
		if errors.Is(err, fs.ErrNotExist) {
			err = os.Mkdir(dir, 0o755)
			if err != nil {
				log.Fatal(err)
			}
		} else {
			log.Fatal(err)
		}
	}
}

// Local Variables:
// tab-width: 4
// End:
