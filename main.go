package main

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"github.com/justinas/alice"
	_ "github.com/lib/pq"
	"golang.org/x/crypto/bcrypt"
)

type App struct {
	DB *sql.DB
}

type Credentials struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type UserResponse struct {
	XataID   string `json:"xata_id"`
	Username string `json:"username"`
}

type ErrorResponse struct {
	Message string `json:"message"`
}

type RouteResponse struct {
	Message string `json:"message"`
	ID      string `json:"id,omitempty"`
}

func main() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatalf("Error loading .env file")
	}
	connStr := os.Getenv("XATA_PSQL_URL")
	if len(connStr) == 0 {
		log.Fatal("XATA_PSQL_URL is not set")
	}

	DB, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal("Error opening database")
	}

	defer DB.Close()

	app := &App{DB: DB}

	log.Println("Starting server on port 8080")
	router := mux.NewRouter()
	log.Println("Setting up routes...")

	router.Handle("/register", alice.New(loggingMiddleware).ThenFunc(app.register)).Methods("POST")
	router.Handle("/login", alice.New(loggingMiddleware).ThenFunc(login)).Methods("POST")
	router.Handle("/project", alice.New(loggingMiddleware).ThenFunc(createProject)).Methods("POST")
	router.Handle("/project/{id}", alice.New(loggingMiddleware).ThenFunc(updateProject)).Methods("PUT")
	router.Handle("/project/{id}", alice.New(loggingMiddleware).ThenFunc(deleteProject)).Methods("DELETE")
	router.Handle("/projects", alice.New(loggingMiddleware).ThenFunc(getProjects)).Methods("GET")
	router.Handle("/projects/{id}", alice.New(loggingMiddleware).ThenFunc(getProject)).Methods("GET")
	log.Fatal(http.ListenAndServe(":8080", router))
}

func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("%s %s %s\n", r.RemoteAddr, r.Method, r.URL)
		next.ServeHTTP(w, r)
	})
}

func respondWithError(w http.ResponseWriter, code int, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(ErrorResponse{Message: message})
}

// User Registration
func (app *App) register(w http.ResponseWriter, r *http.Request) {
	var creds Credentials
	err := json.NewDecoder(r.Body).Decode(&creds)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(creds.Password), bcrypt.DefaultCost)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to hash password")
		return
	}

	var xataID string
	err = app.DB.QueryRow("INSERT INTO \"users\" (username, password) VALUES ($1, $2) RETURNING xata_id", creds.Username, string(hashedPassword)).Scan(&xataID)
	if err != nil {
		log.Printf("Error inserting user: %v", err)
		respondWithError(w, http.StatusInternalServerError, "Failed to insert user")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(UserResponse{XataID: xataID, Username: creds.Username})
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
