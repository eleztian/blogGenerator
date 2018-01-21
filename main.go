package main

import (
	"github.com/eleztian/blog-generator/cli"
	"net/http"
)

func main() {
	cli.Run()
	http.Handle("/", http.FileServer(http.Dir("www/")))
	http.ListenAndServe(":9090", nil)
}
