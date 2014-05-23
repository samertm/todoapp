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

func handleHome(w http.ResponseWriter, req *http.Request) {
	if req.Method == "GET" {
		t, err := template.ParseFiles("view/home.html")
		if err != nil {
			io.WriteString(w, "WHOOPS")
		}
		Session.update(req)
		if username, ok := Session["username"]; ok {
			io.WriteString(w, "YOUR USERNAME: " + username)
		} else {
			t.Execute(w, nil)
		}
	}
}

func (s session) update(req *http.Request) {
		cookies := req.Cookies()
		Session = make(map[string]string)
		for _, c := range cookies {
			if c.Value != "" {
				Session[c.Name] = c.Value
			}
		}
}

func handleLogin(w http.ResponseWriter, req *http.Request) {
	if req.Method == "POST" {
		req.ParseForm()
		form := req.PostForm
		if len(form["username"]) != 0 {
			http.SetCookie(w, &http.Cookie{Name: "username", Value: form["username"][0]})
			io.WriteString(w, form["username"][0])
		} else {
			io.WriteString(w, "ney")
		}
	}
}

func handleTodos(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	io.WriteString(w, "[{\"status\": \"done\", \"name\": \"poopies\"}, {\"status\": \"done\", \"name\": \"my second task\"}]");
}

func ListenAndServe(addr string) {
	port := ":4434"
	fmt.Print("Listening on " + addr + port + "\n")
	http.HandleFunc("/", handleHome)
	http.HandleFunc("/login", handleLogin)
	http.HandleFunc("/todo.json", handleTodos)
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("./static/"))))
	err := http.ListenAndServe(addr+port, nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
