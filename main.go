package main

import (
	"flag"
	"html/template"
	"log"
	"net/http"
	"path/filepath"
	"sync"
)

type templateHandler struct {
	once     sync.Once
	filename string
	templ    *template.Template
}

//ServerHTTPはHTTPリクエストを処理します
func (t *templateHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	t.once.Do(func() {
		t.templ = template.Must(template.ParseFiles(filepath.Join("template", t.filename)))
	})
	t.templ.Execute(w, r)
}

func main() {
	var addr = flag.String("addr", ":8080 ", "アプリケーションのアドレス")
	flag.Parse() // フラグを解釈します。
	r := newRoom()
	// r.tracer = trace.New(os.Stdout)
	http.Handle("/chat", MustAuth(&templateHandler{filename: "chat.html"}))
	http.Handle("/login", &templateHandler{filename: "login.html"})
	http.HandleFunc("/auth/", loginHandler)
	http.Handle("/room", r)
	// チャットルームを開始します。
	go r.run()
	// webサーバを開始
	log.Println("Webサーバを起動します。ポート：", *addr)
	if err := http.ListenAndServe(*addr, nil); err != nil {
		log.Fatal("ListenAdServe:", err)
	}
}
