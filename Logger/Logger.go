package Logger

import (
	"io"
	"io/fs"
	"log"
	"nat/Config"
	"os"
	"sync"
)

var Logger *log.Logger
var once sync.Once

func init() {
	once.Do(func() {
		if Logger == nil {
			fd, err := os.OpenFile(Config.BASE_PATH+"/runtime.log", os.O_CREATE|os.O_APPEND, fs.ModePerm)
			if err != nil {
				log.Fatal(err.Error())
				return
			}
			Logger = log.New(io.MultiWriter(os.Stdout, fd), "", log.LstdFlags)
		}
	})
}
