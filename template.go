package goblin

import (
    "os"
	"path/filepath"
    "log"
)

func dirExists(dir string) bool {
	fi, err := os.Stat(dir)
	if err != nil {
		return false
	}
	return fi.IsDir()
}

func BuildTemplate(dir string, files string) error {
    log.Println("=> folder: " + dir)

    if !dirExists(dir) {
        log.Println("=> folder: " + dir + "not exists!")
        return nil
    }

    error := filepath.Walk(dir, func(path string, f os.FileInfo, err error) error {
        // TODO
        log.Println("=> folder: need build template")

        return nil
    })
    return error
}
