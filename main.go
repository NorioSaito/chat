package main

import (
	"flag"
	"html/template"
	"log"
	"net/http"
	"path/filepath"
	"sync"

	"github.com/stretchr/gomniauth"
	"github.com/stretchr/gomniauth/providers/google"
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
	// Gomniauthのセットアップ
	gomniauth.SetSecurityKey("@Force2868")
	gomniauth.WithProviders(google.New("1063896493054-duevn2cv6p7o066ltrv2ngo81ima742i.apps.googleusercontent.com", "FAC859Qe943W7jJjA8piIGNE", "http://localhost:8080/auth/callback/google"))
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
