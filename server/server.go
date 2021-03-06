package server

import (
	"encoding/json"
	"errors"
	"fmt"
	"html/template"
	"io"
	"log"
	"net/http"
	"net/url"
	"strconv"

	"github.com/samertm/todoapp/engine"
	"github.com/samertm/todoapp/server/session"
)

// warning: modifies req by calling req.ParseForm()
func parseForm(req *http.Request, values ...string) (form url.Values, err error) {
	req.ParseForm()
	form = req.PostForm
	err = checkForm(form, values...)
	return
}

func checkForm(data url.Values, values ...string) error {
	for _, s := range values {
		if len(data[s]) == 0 {
			return errors.New(s + " not passed")
		}
	}
	return nil
}

func handleHome(w http.ResponseWriter, req *http.Request) {
	if req.Method == "GET" {
		t, err := template.ParseFiles("view/home.html")
		if err != nil {
			io.WriteString(w, "WHOOPS")
		}
		t.Execute(w, nil)
	}
}

func handleAddTask(w http.ResponseWriter, req *http.Request) {
	if req.Method == "POST" {
		form, err := parseForm(req,
			"session",
			"todo[status]",
			"todo[name]",
			"todo[description]")
		if err != nil {
			log.Println(err)
			return
		}
		Session.Get <- form["session"][0]
		p := <-Session.Out
		t := engine.NewTask(form["todo[status]"][0],
			form["todo[name]"][0], form["todo[description]"][0])
		if err := checkForm(form, "parentid"); err == nil {
			// attaching a subtask
			i, err := strconv.Atoi(form["parentid"][0])
			if err != nil {
				log.Println(err)
				return
			}
			parentTask, err := engine.FindTask(p.Tasks, i)
			if err != nil {
				log.Println(err)
				return
			}
			parentTask.SubTasks = append(parentTask.SubTasks, &t)
		} else {
			p.Tasks = append(p.Tasks, &t)
		}
	}
}

func handleTasks(w http.ResponseWriter, req *http.Request) {
	if req.Method == "POST" {
		form, err := parseForm(req, "session")
		if err != nil {
			return
		}
		Session.Get <- form["session"][0]
		p := <-Session.Out
		data, err := json.Marshal(p.Tasks)
		if err != nil {
			log.Println(err)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		io.WriteString(w, string(data))
	}
}

func handlePerson(w http.ResponseWriter, req *http.Request) {
	if req.Method == "POST" {
		form, err := parseForm(req, "session")
		if err != nil {
			log.Println("handlePerson error")
			return
		}
		Session.Get <- form["session"][0]
		p := <-Session.Out
		data, err := json.Marshal(p)
		if err != nil {
			log.Println(err)
		}
		w.Header().Set("Content-Type", "application/json")
		io.WriteString(w, string(data))
	}
}

func handleSetUsername(w http.ResponseWriter, req *http.Request) {
	if req.Method == "POST" {
		form, err := parseForm(req, "session", "name")
		if err != nil {
			log.Println("handleSetUsername error")
			return
		}
		Session.Set <- session.Set{form["session"][0], form["name"][0]}
	}
}

func handleTaskDelete(w http.ResponseWriter, req *http.Request) {
	if req.Method == "POST" {
		form, err := parseForm(req, "session", "id")
		if err != nil {
			log.Println("handleTaskDelete error")
			return
		}
		Session.Get <- form["session"][0]
		p := <-Session.Out
		id, err := strconv.Atoi(form["id"][0])
		if err != nil {
			log.Println(err)
			return
		}
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
		form, err := parseForm(req,
			"session",
			"task[id]",
			"task[name]",
			"task[status]",
			"task[description]")
		if err != nil {
			log.Println("handleTaskDelete error")
			return
		}
		Session.Get <- form["session"][0]
		p := <-Session.Out
		id, err := strconv.Atoi(form["task[id]"][0])
		if err != nil {
			log.Println(err)
			return
		}
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

func handlePersonTimeEdit(w http.ResponseWriter, req *http.Request) {
	if req.Method == "POST" {
		form, err := parseForm(req, "session", "goalminutes")
		if err != nil {
			log.Println("handleTaskDelete error")
			return
		}
		Session.Get <- form["session"][0]
		p := <-Session.Out
		minutes, err := strconv.Atoi(form["goalminutes"][0])
		if err != nil {
			log.Println(err)
			return
		}
		p.GoalMinutes = minutes
	}
}

func setHandlers() {
	http.HandleFunc("/", handleHome)
	http.HandleFunc("/addtask", handleAddTask)
	http.HandleFunc("/tasks", handleTasks)
	http.HandleFunc("/task/delete", handleTaskDelete)
	http.HandleFunc("/task/edit", handleTaskEdit)
	http.HandleFunc("/person", handlePerson)
	http.HandleFunc("/person/time/edit", handlePersonTimeEdit)
	http.HandleFunc("/setusername", handleSetUsername)
	http.Handle("/static/",
		http.StripPrefix("/static/",
			http.FileServer(http.Dir("./static/"))))
}

var Session = session.New()

func ListenAndServe(addr string) {
	port := ":4434"
	fmt.Print("Listening on " + addr + port + "\n")
	setHandlers()
	go Session.Run()
	err := http.ListenAndServe(addr+port, nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
