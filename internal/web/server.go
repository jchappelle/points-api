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
	http.ListenAndServe(addr, s.setupHandlers())
}

func (s *Server) setupHandlers() http.Handler {
	mux := mux.NewRouter()
	mux.HandleFunc("/v1/users/{userID}/points/add", s.addPointsHandler)
	mux.HandleFunc("/v1/users/{userID}/payers", s.getPayersHandler)
	mux.HandleFunc("/v1/users/{userID}/points/spend", s.spendPointsHandler)
	return loggingMiddleware(mux)
}

func (s *Server) spendPointsHandler(w http.ResponseWriter, req *http.Request) {
	switch req.Method {
	case http.MethodPost:
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
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		// Try to spend the points
		updatedTransactions, err := s.service.SpendPoints(userID, spendPointsRequest.Points)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		// Return final response
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(updatedTransactions)
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

func (s *Server) getPayersHandler(w http.ResponseWriter, req *http.Request) {
	switch req.Method {
	case http.MethodGet:
		// Get and validate userID
		vars := mux.Vars(req)
		userID := vars["userID"]
		if userID == "" {
			http.Error(w, "userID is required", http.StatusBadRequest)
			return
		}

		accounts := s.service.GetAccounts(userID)

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(accounts)
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

func (s *Server) addPointsHandler(w http.ResponseWriter, req *http.Request) {
	switch req.Method {
	case http.MethodPost:
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
			http.Error(w, err.Error(), http.StatusBadRequest)
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
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}
