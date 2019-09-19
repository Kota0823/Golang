package main

import (
	"html/template"
	"log"
	"net/http"
)
func htmlHandler0(w http.ResponseWriter, r *http.Request) {
	// テンプレートをパース
	t := template.Must(template.ParseFiles("temp.html"))
	str := "Sample Message<br>fwafadfad"

	// テンプレートを描画
	if err := t.ExecuteTemplate(w, "temp.html", str); err != nil {
		log.Fatal(err)
	}
}

func main() {
	http.HandleFunc("/page0", htmlHandler0)
	// サーバーを起動
	http.ListenAndServe(":8989", nil)
}