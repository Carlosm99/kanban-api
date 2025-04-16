package main

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/justinas/alice"
)

type RouteResponse struct {
	Message string `json:"message"`
	ID      string `json:"id,omitempty"`
}

func main() {
	log.Println("Starting server on port 8080")
	router := mux.NewRouter()
	log.Println("Setting up routes...")

	router.Handle("/register", alice.New(loggingMiddleware).ThenFunc(register)).Methods("POST")
	router.Handle("/login", alice.New(loggingMiddleware).ThenFunc(login)).Methods("POST")
	router.Handle("/project", alice.New(loggingMiddleware).ThenFunc(createProject)).Methods("POST")
	router.Handle("/project/{id}", alice.New(loggingMiddleware).ThenFunc(updateProject)).Methods("PUT")
	router.Handle("/project/{id}", alice.New(loggingMiddleware).ThenFunc(deleteProject)).Methods("DELETE")
	router.Handle("/projects", alice.New(loggingMiddleware).ThenFunc(getProjects)).Methods("GET")
	router.Handle("/projects/{id}", alice.New(loggingMiddleware).ThenFunc(getProject)).Methods("GET")
	log.Fatal(http.ListenAndServe(":8080", router))
}

func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request){
		log.Printf("%s %s %s\n", r.RemoteAddr, r.Method, r.URL)
		next.ServeHTTP(w, r)
	})
}

func register(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(RouteResponse{Message: "User registered"})
}

func login(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(RouteResponse{Message: "User logged in"})
}
func createProject(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(RouteResponse{Message: "Project created"})
}
func updateProject(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(RouteResponse{Message: "Project updated", ID: id})
}

func getProjects(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(RouteResponse{Message: "Projects fetched"})
}

func getProject(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(RouteResponse{Message: "Project fetched", ID: id})
}

func deleteProject(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(RouteResponse{Message: "Project deleted", ID: id})
}