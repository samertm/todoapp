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
	"strconv"
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
			len(form["todo[name]"]) == 0 ||
			len(form["todo[description]"]) == 0 {
			// TODO log error
			fmt.Println("handleAddTask error")
			return
		}
		Session.Get <- form["session"][0]
		p := <-Session.Out
		t := engine.NewTask(form["todo[status]"][0],
			form["todo[name]"][0], form["todo[description]"][0])
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

func handlePerson(w http.ResponseWriter, req *http.Request) {
	if req.Method == "POST" {
		req.ParseForm()
		form := req.PostForm
		if len(form["session"]) == 0 {
			// TODO log error
			fmt.Println("handlePerson error")
			return
		}
		Session.Get <- form["session"][0]
		p := <-Session.Out
		data, err := json.Marshal(p)
		if err != nil {
			fmt.Println(err)
		}
		w.Header().Set("Content-Type", "application/json")
		io.WriteString(w, string(data))
	}
}

func handleSetUsername(w http.ResponseWriter, req *http.Request) {
	if req.Method == "POST" {
		req.ParseForm()
		form := req.PostForm
		if len(form["session"]) == 0 ||
			len(form["name"]) == 0 {
			// TODO log error
			fmt.Println("handleSetUsername error")
			return
		}
		Session.Set <- session.Set{form["session"][0], form["name"][0]}
	}
}

func handleTaskDelete(w http.ResponseWriter, req *http.Request) {
	if req.Method == "POST" {
		req.ParseForm()
		form := req.PostForm
		if len(form["session"]) == 0 ||
			len(form["id"]) == 0 {
			fmt.Println("handleTaskDelete error")
			return
		}
		Session.Get <- form["session"][0]
		p := <-Session.Out
		id, _ := strconv.Atoi(form["id"][0])
		for i, t := range p.Tasks {
			if t.Id == id {
				p.Tasks = append(p.Tasks[:i], p.Tasks[i+1:]...)
				break
			}
		}
	}
}

func handleTaskEdit(w http.ResponseWriter, req *http.Request) {
	if req.Method == "POST" {
		req.ParseForm()
		form := req.PostForm
		if len(form["session"]) == 0 ||
			len(form["task[id]"]) == 0 ||
			len(form["task[name]"]) == 0 ||
			len(form["task[status]"]) == 0 ||
			len(form["task[description]"]) == 0 {
			fmt.Println("handleTaskDelete error")
			return
		}
		Session.Get <- form["session"][0]
		p := <-Session.Out
		id, _ := strconv.Atoi(form["task[id]"][0])
		for i, _ := range p.Tasks {
			if p.Tasks[i].Id == id {
				p.Tasks[i].Name = form["task[name]"][0]
				p.Tasks[i].Status = form["task[status]"][0]
				p.Tasks[i].Description = form["task[description]"][0]
				break
			}
		}
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
	http.HandleFunc("/task/delete", handleTaskDelete)
	http.HandleFunc("/task/edit", handleTaskEdit)
	http.HandleFunc("/person", handlePerson)
	http.HandleFunc("/setusername", handleSetUsername)
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("./static/"))))
	go Session.Run()
	err := http.ListenAndServe(addr+port, nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
