package main

import (
	"embed"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/nhlmg93/go-htmx-template/pkg/env"
	"github.com/nhlmg93/go-htmx-template/pkg/logging"
	"github.com/nhlmg93/go-htmx-template/pkg/router"
)

var (
	//go:embed all:templates/*
	templateFS embed.FS

	//go:embed css/output.css
	css embed.FS
)

func main() {
	handleSigTerms()
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
