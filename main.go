package main

import (
	"context"
	"database/sql"
	"embed"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	_ "modernc.org/sqlite"

	"github.com/nhlmg93/go-htmx-template/pkg/env"
	"github.com/nhlmg93/go-htmx-template/pkg/logging"
	"github.com/nhlmg93/go-htmx-template/pkg/router"
	_ "github.com/nhlmg93/go-htmx-template/pkg/sqlc"
)

var (
	//go:embed all:templates/*
	templateFS embed.FS

	//go:embed css/output.css
	css embed.FS
)

func main() {
	handleSigTerms()

	//TODO: move values too env file
	_, err := sql.Open("sqlite", "file:./build/dev.db")
	if err != nil {
		panic(err)
	}

	// TODO: dbmate for migration management.

	// Note: make generate -> make migrate
	//ctx := context.Background()
	//queries := sqlc.New(db)
	// create an author
	//	_, err = queries.CreateTodo(ctx, "walk the dog")
	//	if err != nil {
	//		panic(err)
	//	}
	// list all authors
	//	authors, err := queries.ListTodos(ctx)
	//	if err != nil {
	//		panic(err)
	//	}
	//	log.Println(authors)

	router.SetHtmlTemplates(&templateFS)

	// serve tailwind output
	router.Router.Handle("GET /css/output.css", http.FileServer(http.FS(css)))

	nextRequestID := func() string {
		return fmt.Sprintf("%d", time.Now().UnixNano())
	}
	logger := log.New(os.Stdout, "http: ", log.LstdFlags)
	middleware := logging.Tracing(nextRequestID)(logging.Logging(logger)(router.Router))

	port := env.GetEnvWithDefault("PORT", "8080")
	logger.Println("listening on http://localhost:" + port)
	if err := http.ListenAndServe(":"+port, middleware); err != nil {
		logger.Println("http.ListenAndServe: ", err)
		os.Exit(1)
	}
}
func handleSigTerms() {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		fmt.Println("recieved SIGTERM, exiting")
		os.Exit(1)
	}()
}
