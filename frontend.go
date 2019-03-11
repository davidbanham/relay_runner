package main

import (
	"github.com/gobuffalo/packr/v2"
)

var indexPage string
var cssPage string

func init() {
	box := packr.New("frontend", "./frontend")

	var err error
	indexPage, err = box.FindString("index.html")
	if err != nil {
		panic(err)
	}

	cssPage, err = box.FindString("css/main.css")
	if err != nil {
		panic(err)
	}
}
