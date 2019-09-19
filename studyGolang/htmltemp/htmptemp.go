package htmltemp

import (
	"html/template"
	"net/http"
)

func process(w http.ResponseWriter, r *http.Request){
	t := template.Must(template.ParseFiles("temp.html"))
	m := "Helooooooo!!!<br>kjljohjohjgad" // <--- tmpl.html の {{ . }} に表示されるデータ
	t.Execute(w, m)
}

func Htmltempfunc(){
	server := http.Server{
		Addr: "127.0.0.1:8080",
	}
	http.HandleFunc("/process", process)
	server.ListenAndServe()
}