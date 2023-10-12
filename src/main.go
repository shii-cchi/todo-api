package main

import (
	"encoding/json"
	"fmt"
	"github.com/go-chi/chi"
	"github.com/joho/godotenv"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
)

type todo struct {
	ID     int    `json:"id"`
	Title  string `json:"title"`
	Status string `json:"status"`
}

var todoList []todo

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	godotenv.Load()

	port := os.Getenv("PORT")
	if port == "" {
		log.Fatal("PORT is not found in the environment")
	}

	r := chi.NewRouter()

	r.Get("/", homeHandler)
	r.Mount("/todo", todoHandlers())

	srv := &http.Server{
		Addr:    ":" + port,
		Handler: r,
	}

	log.Printf("Server starting on port %s", port)

	go func() {
		err := srv.ListenAndServe()
		if err != nil {
			log.Fatal(err)
		}
	}()

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
	respondWithJSON(w, http.StatusOK, todoList)
}

func createTodo(w http.ResponseWriter, r *http.Request) {
	var newTodo todo
	err := json.NewDecoder(r.Body).Decode(&newTodo)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	todoList = append(todoList, newTodo)

	respondWithJSON(w, http.StatusCreated, newTodo)
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

	indToUpdate := -1
	for i, elem := range todoList {
		if elem.ID == id {
			indToUpdate = i

			if newTodo.Title != "" {
				todoList[i].Title = newTodo.Title
			}

			if newTodo.Status != "" {
				todoList[i].Status = newTodo.Status
			}

			break
		}
	}

	if indToUpdate == -1 {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	respondWithJSON(w, http.StatusOK, todoList[indToUpdate])
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

	respondWithJSON(w, http.StatusOK, todoList[indToDelete])

	todoList = append(todoList[:indToDelete], todoList[indToDelete+1:]...)
}

func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	data, err := json.Marshal(payload)

	if err != nil {
		log.Printf("Failed to marshal JSON responce: %v", payload)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(data)
}
