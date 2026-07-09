// Package api exposes the HTTP endpoints of the calculator service.
package api

import (
	"encoding/json"
	"errors"
	"net/http"

	"calculator-backend/internal/calculator"
)

type calculateRequest struct {
	A        float64             `json:"a"`
	B        float64             `json:"b"`
	Operator calculator.Operator `json:"operator"`
}

type calculateResponse struct {
	Result float64 `json:"result"`
}

type errorResponse struct {
	Error string `json:"error"`
}

// NewMux builds the HTTP routes for the API.
func NewMux() *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("/api/health", handleHealth)
	mux.HandleFunc("/api/calculate", handleCalculate)
	return mux
}

func handleHealth(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
}

func handleCalculate(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeError(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}

	var req calculateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	result, err := calculator.Calculate(req.A, req.B, req.Operator)
	if err != nil {
		status := http.StatusBadRequest
		var unsupported calculator.ErrUnsupportedOperator
		if errors.As(err, &unsupported) {
			status = http.StatusUnprocessableEntity
		}
		writeError(w, status, err.Error())
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(calculateResponse{Result: result})
}

func writeError(w http.ResponseWriter, status int, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(errorResponse{Error: message})
}
