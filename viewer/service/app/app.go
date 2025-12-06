package app

import (
	"io/fs"
)

func InitSingleViewApp(port ...int) error {
	initServices()
	p := 0
	if len(port) > 0 {
		p = port[0]
	} else {
		p = 10420
	}
	return initRouter(embeddedFiles, p)
}

func InitViewCommandApp(files fs.FS, port int) error {
	initServices()
	return initRouter(files, port)
}
