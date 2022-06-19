package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gorilla/mux"
)

func main() {
	quit := make(chan os.Signal)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)

	router := mux.NewRouter()
	router.StrictSlash(true)

	store, err := New()
	if err != nil {
		log.Fatal(err)
	}

	server := Service{
		store: store,
	}

	router.HandleFunc("/config/", server.createConfigHandler).Methods("POST")
	router.HandleFunc("/config/{id}/{version}/", server.getConfigHandler).Methods("GET")
	router.HandleFunc("/config/{id}/{version}/", server.delConfigHandler).Methods("DELETE")
	//router.HandleFunc("/config/{label}/", server.filterConfigHandler).Methods("GET")

	router.HandleFunc("/group/", server.createGroupHandler).Methods("POST")
	router.HandleFunc("/group/{id}/{version}/", server.getGroupHandler).Methods("GET")
	router.HandleFunc("/group/{id}/{version}/", server.delGroupHandler).Methods("DELETE")
	router.HandleFunc("/group/{id}/{version}/{label}/", server.getGroupLabelHandler).Methods("GET")
	router.HandleFunc("/group/{id}/{version}/", server.updateGroupHandler).Methods("PUT")

	// start server
	srv := &http.Server{Addr: "0.0.0.0:8000", Handler: router}
	go func() {
		log.Println("server starting")
		if err := srv.ListenAndServe(); err != nil {
			if err != http.ErrServerClosed {
				log.Fatal(err)
			}
		}
	}()

	<-quit

	log.Println("service shutting down ...")

	// gracefully stop server
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal(err)
	}
	log.Println("server stopped")
}
