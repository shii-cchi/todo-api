package server

import (
	"database/sql"
	"errors"
	"github.com/go-chi/chi"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"github.com/shii-cchi/todo-api/internal/database"
	dh "github.com/shii-cchi/todo-api/internal/delivery/http"
	"log"
	"net/http"
	"os"
)

type Server struct {
	httpServer  *http.Server
	httpHandler *dh.Handler
	queries     *database.Queries
}

func NewServer(r chi.Router) (*Server, error) {
	err := godotenv.Load(".env")
	if err != nil {
		return nil, err
	}

	port := os.Getenv("PORT")
	if port == "" {
		return nil, errors.New("PORT is not found in the environment")
	}

	dbURL := os.Getenv("DB_URL")
	if dbURL == "" {
		return nil, errors.New("DB_URL is not found in the environment")
	}

	conn, err := sql.Open("postgres", dbURL)
	if err != nil {
		return nil, err
	}

	queries := database.New(conn)

	log.Printf("Server starting on port %s", port)

	handler := dh.New(queries)
	handler.RegisterHTTPEndpoints(r)

	return &Server{
		httpServer: &http.Server{
			Addr:    ":" + port,
			Handler: r,
		},
		httpHandler: handler,
	}, nil
}

func (s *Server) Run() error {
	return s.httpServer.ListenAndServe()
}
