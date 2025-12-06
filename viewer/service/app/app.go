package app

import (
	"embed"
	"io/fs"
)

// 下划线开头的文件不能被直接嵌入，需要使用all
//
//go:embed all:dist/*
var embeddedFiles embed.FS

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
