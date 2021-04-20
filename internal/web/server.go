package web

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"fetchrewards.com/points-api/internal/model"
)

type spendPointsRequest struct {
	Points int `json:"points"`
}

type pointService interface {
	AddPoints(ctx context.Context, transaction model.Transaction) error
	GetAccounts(ctx context.Context) ([]model.Account, error)
	SpendPoints(ctx context.Context, points int) ([]model.Transaction, error)
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
	mux := http.NewServeMux()
	mux.HandleFunc("/v1/points/add", s.addPointsHandler)
	mux.HandleFunc("/v1/payers", s.getPayersHandler)
	mux.HandleFunc("/v1/points/spend", s.spendPointsHandler)
	return loggingMiddleware(mux)
}

func (s *Server) spendPointsHandler(w http.ResponseWriter, req *http.Request) {
	switch req.Method {
	case http.MethodPost:
		spendPointsRequest := spendPointsRequest{}
		err := json.NewDecoder(req.Body).Decode(&spendPointsRequest)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		updatedTransactions, err := s.service.SpendPoints(req.Context(), spendPointsRequest.Points)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(updatedTransactions)
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

func (s *Server) getPayersHandler(w http.ResponseWriter, req *http.Request) {
	switch req.Method {
	case http.MethodGet:
		accounts, err := s.service.GetAccounts(req.Context())
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(accounts)
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

func (s *Server) addPointsHandler(w http.ResponseWriter, req *http.Request) {
	switch req.Method {
	case http.MethodPost:
		transaction := model.Transaction{}
		err := json.NewDecoder(req.Body).Decode(&transaction)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		s.service.AddPoints(req.Context(), transaction)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNoContent)
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}
