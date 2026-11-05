package view

import (
	"embed"
	"io/fs"
)

//go:embed templates static
var Files embed.FS

func TemplatesFS() fs.FS {
	sub, err := fs.Sub(Files, "templates")
	if err != nil {
		panic(err)
	}
	return sub
}

func StaticFS() fs.FS {
	sub, err := fs.Sub(Files, "static")
	if err != nil {
		panic(err)
	}
	return sub
}
