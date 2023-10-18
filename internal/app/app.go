package app

import (
	"github.com/go-chi/chi"
	_ "github.com/lib/pq"
	"github.com/shii-cchi/todo-api/internal/server"
	"log"
)

func Run() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	r := chi.NewRouter()

	srv, err := server.NewServer(r)

	if err != nil {
		log.Fatal(err)
	}

	go func() {
		if err := srv.Run(); err != nil {
			log.Fatal(err)
		}
	}()

	select {}
}
