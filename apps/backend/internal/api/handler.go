// Package api exposes the HTTP endpoints of the calculator service. It
// decodes/encodes JSON and maps service errors to HTTP status codes;
// the actual calculation logic lives in internal/service and
// internal/calculator.
package api

import (
	"encoding/json"
	"net/http"

	"calculator-backend/internal/service"
)

// Handler wires HTTP routes to a CalculatorService.
type Handler struct {
	svc *service.CalculatorService
}

// NewHandler creates a Handler backed by svc.
func NewHandler(svc *service.CalculatorService) *Handler {
	return &Handler{svc: svc}
}

// NewMux builds the HTTP routes for the API, backed by svc.
func NewMux(svc *service.CalculatorService) *http.ServeMux {
	h := NewHandler(svc)

	mux := http.NewServeMux()
	mux.HandleFunc("GET /api/health", h.handleHealth)
	mux.HandleFunc("POST /api/add", h.handleAdd)
	mux.HandleFunc("POST /api/subtract", h.handleSubtract)
	mux.HandleFunc("POST /api/multiply", h.handleMultiply)
	mux.HandleFunc("POST /api/divide", h.handleDivide)
	mux.HandleFunc("POST /api/power", h.handlePower)
	mux.HandleFunc("POST /api/sqrt", h.handleSqrt)
	mux.HandleFunc("POST /api/percentage", h.handlePercentage)
	return mux
}

func (h *Handler) handleHealth(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, map[string]string{"status": "ok"})
}

func (h *Handler) handleAdd(w http.ResponseWriter, r *http.Request) {
	var req binaryRequest
	if !decode(w, r, &req) {
		return
	}
	result, err := h.svc.Add(req.A, req.B)
	respond(w, result, err)
}

func (h *Handler) handleSubtract(w http.ResponseWriter, r *http.Request) {
	var req binaryRequest
	if !decode(w, r, &req) {
		return
	}
	result, err := h.svc.Subtract(req.A, req.B)
	respond(w, result, err)
}

func (h *Handler) handleMultiply(w http.ResponseWriter, r *http.Request) {
	var req binaryRequest
	if !decode(w, r, &req) {
		return
	}
	result, err := h.svc.Multiply(req.A, req.B)
	respond(w, result, err)
}

func (h *Handler) handleDivide(w http.ResponseWriter, r *http.Request) {
	var req binaryRequest
	if !decode(w, r, &req) {
		return
	}
	result, err := h.svc.Divide(req.A, req.B)
	respond(w, result, err)
}

func (h *Handler) handlePower(w http.ResponseWriter, r *http.Request) {
	var req powerRequest
	if !decode(w, r, &req) {
		return
	}
	result, err := h.svc.Power(req.Base, req.Exponent)
	respond(w, result, err)
}

func (h *Handler) handleSqrt(w http.ResponseWriter, r *http.Request) {
	var req sqrtRequest
	if !decode(w, r, &req) {
		return
	}
	result, err := h.svc.Sqrt(req.Value)
	respond(w, result, err)
}

func (h *Handler) handlePercentage(w http.ResponseWriter, r *http.Request) {
	var req percentageRequest
	if !decode(w, r, &req) {
		return
	}
	result, err := h.svc.Percentage(req.Value, req.Percent)
	respond(w, result, err)
}

// decode reads and validates the JSON request body into dst. On failure it
// writes a 400 response and returns false.
func decode(w http.ResponseWriter, r *http.Request, dst any) bool {
	if r.Body == nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return false
	}
	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()
	if err := dec.Decode(dst); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return false
	}
	return true
}

// respond writes a calculation result, or maps a service error to a 400
// response. All service-level errors (invalid input, division by zero,
// negative square root, undefined results) stem from client-supplied
// operands, so they're all reported as Bad Request.
func respond(w http.ResponseWriter, result float64, err error) {
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, calculateResponse{Result: result})
}

func writeError(w http.ResponseWriter, status int, message string) {
	writeJSON(w, status, errorResponse{Error: message})
}

func writeJSON(w http.ResponseWriter, status int, body any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(body)
}
