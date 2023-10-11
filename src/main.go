package main

import (
	"encoding/json"
	"fmt"
	"github.com/go-chi/chi"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"
)

type todo struct {
	ID     int    `json:"id"`
	Title  string `json:"title"`
	Status string `json:"status"`
}

var todoList []todo

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	log.Println("Starting server...")

	r := chi.NewRouter()

	r.Get("/", homeHandler)
	r.Mount("/todo", todoHandlers())

	srv := &http.Server{
		Addr:         ":8080",
		Handler:      r,
		ReadTimeout:  60 * time.Second,
		WriteTimeout: 60 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	go func() {
		err := srv.ListenAndServe()
		checkErr(err)
	}()

	log.Println("server start listening on port 8080")
	select {}
}

func homeHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "my todo list")
}

func todoHandlers() http.Handler {
	rg := chi.NewRouter()
	rg.Group(func(r chi.Router) {
		r.Get("/", fetchTodos)
		r.Post("/", createTodo)
		r.Put("/{id}", updateTodo)
		r.Delete("/{id}", deleteTodo)
	})

	return rg
}

func fetchTodos(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(todoList)

	//fmt.Fprintf(w, "get todos '%v'", todoList)
}

func createTodo(w http.ResponseWriter, r *http.Request) {
	var newTodo todo
	err := json.NewDecoder(r.Body).Decode(&newTodo)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	todoList = append(todoList, newTodo)

	fmt.Fprintf(w, "post todo '%v'", newTodo)
}

func updateTodo(w http.ResponseWriter, r *http.Request) {
	idStr := strings.TrimSpace(chi.URLParam(r, "id"))
	id, err := strconv.Atoi(idStr)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	var newTodo todo
	err = json.NewDecoder(r.Body).Decode(&newTodo)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var exist bool
	for i, elem := range todoList {
		if elem.ID == id {
			if newTodo.Title != "" {
				todoList[i].Title = newTodo.Title
			}

			if newTodo.Status != "" {
				todoList[i].Status = newTodo.Status
			}

			exist = true
			break
		}
	}

	if !exist {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	fmt.Fprintf(w, "Element with id = %d has been updated", id)
}

func deleteTodo(w http.ResponseWriter, r *http.Request) {
	idStr := strings.TrimSpace(chi.URLParam(r, "id"))
	id, err := strconv.Atoi(idStr)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	indToDelete := -1
	for i, elem := range todoList {
		if elem.ID == id {
			indToDelete = i
			break
		}
	}

	if indToDelete == -1 {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	todoList = append(todoList[:indToDelete], todoList[indToDelete+1:]...)

	fmt.Fprintf(w, "Element with id = %d has been deleted", id)
}

func checkErr(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
