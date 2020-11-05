package main

import (
	"fmt"
	"net/http"
)

func main() {
	http.HandleFunc("/download/apk", downloadFile)
	fmt.Println("Server run on port 5001 for apk download")
	http.ListenAndServe(":5001", nil)
}

func downloadFile(w http.ResponseWriter, r *http.Request) {
	file := "delog201102.apk"

	// 設定此 Header 告訴瀏覽器下載檔案。 如果沒設定則會在新的 tab 開啟檔案。
	// w.Header().Set("Content-Disposition", "attachment; filename="+file)

	http.ServeFile(w, r, file)
}
