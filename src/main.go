package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/go-chi/chi"
	"github.com/google/uuid"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"github.com/shii-cchi/todo-api/internal/database"
	"log"
	"net/http"
	"os"
)

type apiConfig struct {
	DB *database.Queries
}

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	godotenv.Load()

	port := os.Getenv("PORT")
	if port == "" {
		log.Fatal("PORT is not found in the environment")
	}

	dbURL := os.Getenv("DB_URL")
	if dbURL == "" {
		log.Fatal("DB_URL is not found in the environment")
	}

	conn, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatal("Can't connect to database")
	}

	apiCfg := apiConfig{
		DB: database.New(conn),
	}

	r := chi.NewRouter()

	r.Get("/", homeHandler)
	r.Mount("/users", apiCfg.userHandlers())
	r.Mount("/todo", apiCfg.todoHandlers())

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

func (apiCfg *apiConfig) todoHandlers() http.Handler {
	rg := chi.NewRouter()
	rg.Group(func(r chi.Router) {
		r.Get("/", apiCfg.handlerFetchTodos)
		r.Get("/{id}", apiCfg.handlerFetchTodo)
		r.Post("/", apiCfg.handlerCreateTodo)
		r.Put("/{id}", apiCfg.handlerUpdateTodo)
		r.Delete("/{id}", apiCfg.handlerDeleteTodo)
	})

	return rg
}

func (apiCfg *apiConfig) handlerFetchTodos(w http.ResponseWriter, r *http.Request) {
	todoList, err := apiCfg.DB.GetTodosList(r.Context())

	if err != nil {
		respondWithError(w, http.StatusNotFound, fmt.Sprintf("Coudn't get todos: %v", err))
		return
	}

	respondWithJSON(w, http.StatusOK, databaseTodoListtoTodoList(todoList))
}

func (apiCfg *apiConfig) handlerFetchTodo(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := uuid.Parse(idStr)

	if err != nil {
		respondWithError(w, http.StatusBadRequest, fmt.Sprintf("Coudn't parse id: %v", err))
		return
	}

	todo, err := apiCfg.DB.GetTodo(r.Context(), id)

	if err != nil {
		respondWithError(w, http.StatusNotFound, fmt.Sprintf("Coudn't get todos: %v", err))
		return
	}

	respondWithJSON(w, http.StatusOK, databaseTodotoTodo(todo))
}

func (apiCfg *apiConfig) handlerCreateTodo(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Title  string `json:"title"`
		Status string `json:"status"`
	}

	decoder := json.NewDecoder(r.Body)

	params := parameters{}
	err := decoder.Decode(&params)

	if err != nil {
		respondWithError(w, http.StatusBadRequest, fmt.Sprintf("Error parsing JSON: %v", err))
		return
	}

	apiKey := r.Header.Get("Authorization")

	if apiKey == "" {
		respondWithError(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	userId, err := apiCfg.DB.CheckApiKey(r.Context(), apiKey)

	if err != nil {
		respondWithError(w, http.StatusBadRequest, fmt.Sprintf("Coudn't check: %v", err))
		return
	}

	if userId == uuid.Nil {
		respondWithError(w, http.StatusUnauthorized, fmt.Sprintf("fffffffffff: %v", err))
		return
	}

	todo, err := apiCfg.DB.CreateTodo(r.Context(), database.CreateTodoParams{
		ID:     uuid.New(),
		Title:  params.Title,
		Status: params.Status,
		UserID: userId,
	})

	if err != nil {
		respondWithError(w, http.StatusBadRequest, fmt.Sprintf("Coudn't create todo: %v", err))
		return
	}

	respondWithJSON(w, http.StatusCreated, databaseTodotoTodo(todo))
}

func (apiCfg *apiConfig) handlerUpdateTodo(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := uuid.Parse(idStr)

	if err != nil {
		respondWithError(w, http.StatusBadRequest, fmt.Sprintf("Coudn't parse id: %v", err))
		return
	}

	todoOld, err := apiCfg.DB.GetTodo(r.Context(), id)

	if err != nil {
		respondWithError(w, http.StatusNotFound, fmt.Sprintf("Coudn't find todos: %v", err))
		return
	}

	type parameters struct {
		Title  string `json:"title"`
		Status string `json:"status"`
	}

	decoder := json.NewDecoder(r.Body)

	params := parameters{}
	err = decoder.Decode(&params)

	if err != nil {
		respondWithError(w, http.StatusBadRequest, fmt.Sprintf("Error parsing JSON: %v", err))
		return
	}

	todoNew, err := apiCfg.DB.UpdateTodo(r.Context(), database.UpdateTodoParams{
		ID:     id,
		Title:  params.Title,
		Status: params.Status,
	})

	if err != nil {
		respondWithError(w, http.StatusBadRequest, fmt.Sprintf("Coudn't update user: %v", err))
		return
	}

	if todoNew.Title == todoOld.Title && todoNew.Status == todoOld.Status {
		respondWithError(w, http.StatusNoContent, fmt.Sprintf("Coudn't update user: %v", err))
		return
	}

	respondWithJSON(w, http.StatusOK, databaseTodotoTodo(todoNew))
}

func (apiCfg *apiConfig) handlerDeleteTodo(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := uuid.Parse(idStr)

	if err != nil {
		respondWithError(w, http.StatusBadRequest, fmt.Sprintf("Coudn't parse id: %v", err))
		return
	}

	todo, err := apiCfg.DB.GetTodo(r.Context(), id)

	if err != nil {
		respondWithError(w, http.StatusNotFound, fmt.Sprintf("Coudn't find todos: %v", err))
		return
	}

	err = apiCfg.DB.DeleteTodo(r.Context(), id)

	if err != nil {
		respondWithError(w, http.StatusBadRequest, fmt.Sprintf("Coudn't delete: %v", err))
		return
	}

	respondWithJSON(w, http.StatusOK, databaseTodotoTodo(todo))
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

func respondWithError(w http.ResponseWriter, code int, msg string) {
	type errResponse struct {
		Error string `json:"error"`
	}

	respondWithJSON(w, code, errResponse{
		Error: msg,
	})
}

func (apiCfg *apiConfig) userHandlers() http.Handler {
	rg := chi.NewRouter()
	rg.Post("/", apiCfg.handlerCreateUser)
	return rg
}

func (apiCfg *apiConfig) handlerCreateUser(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		ID   uuid.UUID `json:"id"`
		Name string    `json:"name"`
	}

	decoder := json.NewDecoder(r.Body)

	params := parameters{}
	err := decoder.Decode(&params)

	if err != nil {
		respondWithError(w, http.StatusBadRequest, fmt.Sprintf("Error parsing JSON: %v", err))
		return
	}

	user, err := apiCfg.DB.CreateUser(r.Context(), database.CreateUserParams{
		ID:   uuid.New(),
		Name: params.Name,
	})

	if err != nil {
		respondWithError(w, http.StatusBadRequest, fmt.Sprintf("Coudn't create user: %v", err))
		return
	}

	respondWithJSON(w, http.StatusCreated, databaseUsertoUser(user))
}
