package main

import (
	"flag"
	"log"
	"net/http"
)

var (
	port      string
	directory string
)

func main() {
	flag.Parse()

	log.Printf("fserver listen on port:%v, root:%v", port, directory)
	log.Fatal(http.ListenAndServe(":"+port, http.FileServer(http.Dir(directory))))
}

func init() {
	flag.StringVar(&port, "p", "8080", "fserver listen port")
	flag.StringVar(&directory, "d", ".", "fserver root directory")
}
