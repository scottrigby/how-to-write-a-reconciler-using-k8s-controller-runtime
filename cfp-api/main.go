package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/scottrigby/cfp-api/pkg/handlers"
)

const defaultPort = 50001

func main() {
	r := mux.NewRouter()

	RegisterSpeakerRoutes(r)
	RegisterProposaltRoutes(r)

	log.Printf("listening on port %v\n", defaultPort)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", defaultPort), r))
}

func RegisterSpeakerRoutes(router *mux.Router) {
	router.HandleFunc("/api/speakers", handlers.GetSpeakers).Methods("GET")
	router.HandleFunc("/api/speakers/{id}", handlers.GetSpeakerById).Methods("GET")
	router.HandleFunc("/api/speakers", handlers.CreateSpeaker).Methods("POST")
	router.HandleFunc("/api/speakers/{id}", handlers.UpdateSpeaker).Methods("PUT")
	router.HandleFunc("/api/speakers/{id}", handlers.DeleteSpeaker).Methods("DELETE")
}

func RegisterProposaltRoutes(router *mux.Router) {
	router.HandleFunc("/api/proposals", handlers.GetProposals).Methods("GET")
	router.HandleFunc("/api/proposals/{id}", handlers.GetProposalById).Methods("GET")
	router.HandleFunc("/api/proposals", handlers.CreateProposal).Methods("POST")
	router.HandleFunc("/api/proposals/{id}", handlers.UpdateProposal).Methods("PUT")
	router.HandleFunc("/api/proposals/{id}", handlers.DeleteProposal).Methods("DELETE")
}
