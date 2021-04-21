package web

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"fetchrewards.com/points-api/internal/model"
	"github.com/gorilla/mux"
)

type spendPointsRequest struct {
	Points int `json:"points"`
}

// pointService is an abstraction for the service layer methods the web server depends on
type pointService interface {
	AddPoints(userID string, transaction model.Transaction) error
	GetAccounts(userID string) []model.Account
	SpendPoints(userID string, points int) ([]model.Transaction, error)
}

// Server provides functionality for starting the server and routing web requests to the
// appropriate handlers. Handlers delegate business logic to a pointService interface
// for executing the business logic.
type Server struct {
	service pointService
}

// NewServer creates a new Server configured with the given pointService
func NewServer(service pointService) *Server {
	return &Server{
		service: service,
	}
}

// Start starts the web server on the given port
func (s *Server) Start(port int) {
	addr := fmt.Sprintf(":%d", port)
	log.Printf("Starting web server, listening at %s", addr)

	err := http.ListenAndServe(addr, s.setupHandlers())
	if err != nil {
		log.Fatal(err)
	}
}

func (s *Server) setupHandlers() http.Handler {
	router := mux.NewRouter()
	router.HandleFunc("/v1/users/{userID}/points/add", s.addPointsHandler).Methods("POST")
	router.HandleFunc("/v1/users/{userID}/payers", s.getPayersHandler).Methods("GET")
	router.HandleFunc("/v1/users/{userID}/points/spend", s.spendPointsHandler).Methods("POST")
	return loggingMiddleware(router)
}

func (s *Server) spendPointsHandler(w http.ResponseWriter, req *http.Request) {
	// Get and validate userID
	vars := mux.Vars(req)
	userID := vars["userID"]
	if userID == "" {
		http.Error(w, "userID is required", http.StatusBadRequest)
		return
	}

	// Marshal request into a struct
	spendPointsRequest := spendPointsRequest{}
	err := json.NewDecoder(req.Body).Decode(&spendPointsRequest)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnprocessableEntity)
		return
	}

	// Try to spend the points
	newTransactions, err := s.service.SpendPoints(userID, spendPointsRequest.Points)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Return final response
	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(newTransactions)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (s *Server) getPayersHandler(w http.ResponseWriter, req *http.Request) {
	// Get and validate userID
	vars := mux.Vars(req)
	userID := vars["userID"]
	if userID == "" {
		http.Error(w, "userID is required", http.StatusBadRequest)
		return
	}

	accounts := s.service.GetAccounts(userID)

	w.Header().Set("Content-Type", "application/json")
	err := json.NewEncoder(w).Encode(accounts)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (s *Server) addPointsHandler(w http.ResponseWriter, req *http.Request) {
	// Get and validate userID
	vars := mux.Vars(req)
	userID := vars["userID"]
	if userID == "" {
		http.Error(w, "userID is required", http.StatusBadRequest)
		return
	}

	// Marshal request into a struct
	transaction := model.Transaction{}
	err := json.NewDecoder(req.Body).Decode(&transaction)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnprocessableEntity)
		return
	}

	// Try to add the transaction
	err = s.service.AddPoints(userID, transaction)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusNoContent)
}
