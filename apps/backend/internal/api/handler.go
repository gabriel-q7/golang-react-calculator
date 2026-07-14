// Package api exposes the HTTP endpoints of the calculator service. It
// decodes/encodes JSON and maps errors to HTTP status codes; the actual
// calculation logic lives in internal/service and internal/operations,
// which this package only knows about through small interfaces
// (Executor, RateLimiter) — never their concrete types.
package api

import (
	"encoding/json"
	"errors"
	"net/http"

	"calculator-backend/internal/operations"
)

// Executor runs a named operation against a set of operand values. It's
// satisfied by *service.CalculatorService, but this package doesn't
// import the service package at all — it only depends on this shape.
type Executor interface {
	Execute(name string, values map[string]float64) (float64, error)
}

// NewMux builds the HTTP routes for the API: a health check, plus one
// POST route per operation in ops. Adding a new operation to the
// registry that produced ops is enough for it to get a working route —
// this function never needs to change.
//
// limiter, if non-nil, rate-limits every operation route (not the health
// check) keyed by client IP. Pass nil to disable rate limiting (tests
// that don't care about it do this, to keep request counts predictable).
func NewMux(exec Executor, ops []operations.Operation, limiter RateLimiter) *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /api/health", handleHealth)

	for _, op := range ops {
		op := op
		var handler http.Handler = makeOperationHandler(exec, op)
		if limiter != nil {
			handler = rateLimitMiddleware(limiter, clientIPKey, handler)
		}
		mux.Handle("POST /api/"+op.Name(), handler)
	}

	return mux
}

func handleHealth(w http.ResponseWriter, _ *http.Request) {
	writeJSON(w, http.StatusOK, map[string]string{"status": "ok"})
}

// makeOperationHandler builds the HTTP handler for a single operation:
// decode its named operands from the request body, execute it, respond.
// This is the one function shared by all 7 (and any future) operations —
// none of them need their own handler.
func makeOperationHandler(exec Executor, op operations.Operation) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		values, err := decodeOperands(r, op.Operands())
		if err != nil {
			writeError(w, http.StatusBadRequest, err.Error())
			return
		}

		result, err := exec.Execute(op.Name(), values)
		if err != nil {
			writeError(w, http.StatusBadRequest, err.Error())
			return
		}

		writeJSON(w, http.StatusOK, calculateResponse{Result: result})
	}
}

var errInvalidRequestBody = errors.New("invalid request body")

// decodeOperands reads a JSON object from r's body and checks it has
// exactly the keys operands lists — no more (rejecting unknown fields),
// no fewer (rejecting missing ones). This is what lets every operation
// share one handler despite each expecting different named fields.
func decodeOperands(r *http.Request, operands []string) (map[string]float64, error) {
	if r.Body == nil {
		return nil, errInvalidRequestBody
	}

	var raw map[string]float64
	if err := json.NewDecoder(r.Body).Decode(&raw); err != nil {
		return nil, errInvalidRequestBody
	}
	if len(raw) != len(operands) {
		return nil, errInvalidRequestBody
	}

	values := make(map[string]float64, len(operands))
	for _, name := range operands {
		v, ok := raw[name]
		if !ok {
			return nil, errInvalidRequestBody
		}
		values[name] = v
	}
	return values, nil
}

func writeError(w http.ResponseWriter, status int, message string) {
	writeJSON(w, status, errorResponse{Error: message})
}

func writeJSON(w http.ResponseWriter, status int, body any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(body)
}
