package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/leeyzero/go-tools/utils"
)

var (
	gAddr      string
	gTargetDir string
	gMaxMemory int64
)

func uploadHandler(w http.ResponseWriter, r *http.Request) {
	err := r.ParseMultipartForm(gMaxMemory)
	if err != nil {
		fmt.Fprintln(w, err)
		return
	}

	files := r.MultipartForm.File["files"]
	for i, _ := range files {
		file, err := files[i].Open()
		if err != nil {
			fmt.Fprintln(w, err)
			return
		}
		defer file.Close()

		dst, err := os.Create(strings.TrimRight(gTargetDir, "/") + "/" + files[i].Filename)
		if err != nil {
			fmt.Fprintln(w, err)
			return
		}
		defer dst.Close()

		_, err = io.Copy(dst, file)
		if err != nil {
			fmt.Fprintln(w, err)
			return
		}

		fmt.Fprintf(w, "upload %s success\n", files[i].Filename)
	}
}

func init() {
	gAddr = utils.TryGetEnvString("ADDR", ":8080")
	gTargetDir = utils.TryGetEnvString("TARGET_DIR", "/tmp")
	gMaxMemory = utils.TryGetEnvInt64("MAX_MEMORY", 50<<20)
}

func main() {
	log.Printf("serve on ADDR=%s TARGET_DIR=%s MAX_MEMORY=%d\n", gAddr, gTargetDir, gMaxMemory)
	http.HandleFunc("/upload", uploadHandler)
	log.Fatal(http.ListenAndServe(gAddr, nil))
}
