package server

import (
	"fmt"
	"html/template"
	"io"
	"log"
	"net/http"
	_ "time"
)

type session map[string]string

var Session session
var times int

func home(w http.ResponseWriter, req *http.Request) {
	if req.Method == "GET" {
		t, err := template.ParseFiles("view/home.html")
		if err != nil {
			io.WriteString(w, "WHOOPS")
		}
		cookies := req.Cookies()
		Session = make(map[string]string)
		for _, c := range cookies {
			if c.Value != "" {
				Session[c.Name] = c.Value
			}
		}
		if username, ok := Session["username"]; ok {
			// http.SetCookie(w, &http.Cookie{Name: "username", Value: "poop", MaxAge: -1})
			io.WriteString(w, "YOUR USERNAME: " + username)
			fmt.Println(times)
			times++
		} else {
			http.SetCookie(w, &http.Cookie{Name: "username", Value: "pee"})
			t.Execute(w, nil)
		}
	}
}

func ListenAndServe(addr string) {
	port := ":4434"
	fmt.Print("Listening on " + addr + port + "\n")
	http.HandleFunc("/", home)
	err := http.ListenAndServe(addr+port, nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
