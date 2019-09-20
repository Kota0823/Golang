package main

import (
	"html/template"
	"log"
	"net/http"
)


func htmlHandler0(w http.ResponseWriter, r *http.Request) {
	
type Information struct {
	Index string
	En1   string
	En2   string
}

var tunnels = make(map[[16]byte]Information) //マップ型変数
	tunnels[[16]byte{0}] = Information{
		Index: "1",
		En1:   "192.168.100.1",
		En2:   "192.168.100.2",
	}
	tunnels[[16]byte{1}] = Information{
		Index: "2",
		En1:   "192.168.200.3",
		En2:   "192.168.200.4",
	}
	// テンプレートをパース
	t := template.Must(template.ParseFiles("temp.html"))

	if err := t.ExecuteTemplate(w, "temp.html", tunnels); err != nil {
		log.Fatal(err)
	}

}

func main() {
	http.HandleFunc("/page0", htmlHandler0)
	// サーバーを起動
	http.ListenAndServe(":8989", nil)
}
