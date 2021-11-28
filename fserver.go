package main

import (
	"log"
	"net/http"

	"github.com/leeyzero/go-tools/utils"
)

var (
	gAddr    string
	gRootDir string
)

func init() {
	gAddr = utils.TryGetEnvString("ADDR", ":8080")
	gRootDir = utils.TryGetEnvString("ROOT_DIR", ".")
}

func main() {
	log.Printf("serve on ADDR=%v, ROOT_DIR:%v", gAddr, gRootDir)
	log.Fatal(http.ListenAndServe(gAddr, http.FileServer(http.Dir(gRootDir))))
}
