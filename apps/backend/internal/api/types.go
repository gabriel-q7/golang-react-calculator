package api

// calculateResponse is the response body for a successful calculation.
type calculateResponse struct {
	Result float64 `json:"result"`
}

// errorResponse is the response body for a failed request, used for
// every error this API returns (decode failures, domain errors, rate
// limiting) so clients always see the same shape.
type errorResponse struct {
	Error string `json:"error"`
}
