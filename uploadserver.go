package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
)

const (
	// 最大上传文件大小50MB
	MAX_UPLOAD_SIZE = 50 << 20

	// 文件保证位置
	SAVE_PATH = "/tmp/"
)

func uploadHandler(w http.ResponseWriter, r *http.Request) {
	err := r.ParseMultipartForm(MAX_UPLOAD_SIZE)
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

		dst, err := os.Create(SAVE_PATH + files[i].Filename)
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

func main() {
	addr := ":8080"
	if len(os.Args) > 1 {
		addr = os.Args[1]
	}

	http.HandleFunc("/upload", uploadHandler)
	log.Fatal(http.ListenAndServe(addr, nil))
}
