package spanner_demo_20240801

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"

	"cloud.google.com/go/spanner"
)

func Ignition(ctx context.Context) {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
		log.Printf("Defaulting to port %s", port)
	}

	log.Printf("Listening on port %s", port)
	static := http.StripPrefix("/static", http.FileServer(http.Dir("static")))
	http.Handle("/static/", static)

	dbName := fmt.Sprintf("projects/%s/instances/%s/databases/%s", "gcpug-public-spanner", "merpay-sponsored-instance", "sinmetal")
	dspc := spanner.DefaultSessionPoolConfig
	dspc.MinOpened = 3
	spa, err := spanner.NewClientWithConfig(ctx, dbName,
		spanner.ClientConfig{
			SessionPoolConfig: dspc,
		})
	if err != nil {
		log.Fatalln(err)
	}

	s, err := NewService(ctx, spa)
	if err != nil {
		log.Fatalln(err)
	}
	messageHandler := MessageHandler{
		s: s,
	}
	http.HandleFunc("/api/search", messageHandler.SearchHandler)
	http.HandleFunc("/api/postMessage", messageHandler.PostMessageHandler)
	http.HandleFunc("/", StaticContentsHandler)

	if err := http.ListenAndServe(fmt.Sprintf(":%s", port), http.DefaultServeMux); err != nil {
		log.Printf("failed ListenAndServe err=%+v", err)
	}
}
