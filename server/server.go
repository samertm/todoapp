package server

import (
	"encoding/json"
	"fmt"
	"html/template"
	"io"
	"log"
	"net/http"
	_ "time"

	"github.com/samertm/todoapp/engine"
	"github.com/samertm/todoapp/server/session"
)

func handleHome(w http.ResponseWriter, req *http.Request) {
	if req.Method == "GET" {
		t, err := template.ParseFiles("view/home.html")
		if err != nil {
			io.WriteString(w, "WHOOPS")
		}
		t.Execute(w, nil)
	}
}

func handleLogin(w http.ResponseWriter, req *http.Request) {
	// if req.Method == "POST" {
	// 	req.ParseForm()
	// 	form := req.PostForm
	// 	if len(form["username"]) != 0 {
	// 		http.SetCookie(w, &http.Cookie{Name: "username", Value: form["username"][0]})
	// 		io.WriteString(w, form["username"][0])
	// 	} else {
	// 		io.WriteString(w, "ney")
	// 	}
	// }
}

func handleAddTask(w http.ResponseWriter, req *http.Request) {
	if req.Method == "POST" {
		req.ParseForm()
		form := req.PostForm
		if len(form["session"]) == 0 ||
			len(form["todo[status]"]) == 0 ||
			len(form["todo[name]"]) == 0 {
			// TODO log error
			fmt.Println("submission error")
			return
		}
		Session.Get <- form["session"][0]
		p := <-Session.Out
		t := engine.NewTask(form["todo[status]"][0],
			form["todo[name]"][0])
		p.Tasks = append(p.Tasks, t)
	}
}

func handleTasks(w http.ResponseWriter, req *http.Request) {
	if req.Method == "POST" {
		req.ParseForm()
		form := req.PostForm
		if len(form["session"]) == 0 {
			// TODO log error
			return
		}
		Session.Get <- form["session"][0]
		p := <-Session.Out
		data, err := json.Marshal(p.Tasks)
		if err != nil {
			fmt.Println(err)
		}
		w.Header().Set("Content-Type", "application/json")
		io.WriteString(w, string(data))
	}
}

var Session = session.New()

func ListenAndServe(addr string) {
	port := ":4434"
	fmt.Print("Listening on " + addr + port + "\n")
	http.HandleFunc("/", handleHome)
	http.HandleFunc("/login", handleLogin)
	http.HandleFunc("/addtask", handleAddTask)
	http.HandleFunc("/tasks", handleTasks)
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("./static/"))))
	go Session.Run()
	err := http.ListenAndServe(addr+port, nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
