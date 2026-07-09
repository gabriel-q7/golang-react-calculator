package api

// binaryRequest is the request body for operations taking two named
// operands: add, subtract, multiply, divide.
type binaryRequest struct {
	A float64 `json:"a"`
	B float64 `json:"b"`
}

// powerRequest is the request body for exponentiation.
type powerRequest struct {
	Base     float64 `json:"base"`
	Exponent float64 `json:"exponent"`
}

// sqrtRequest is the request body for square root.
type sqrtRequest struct {
	Value float64 `json:"value"`
}

// percentageRequest is the request body for percentage: percent% of value.
type percentageRequest struct {
	Value   float64 `json:"value"`
	Percent float64 `json:"percent"`
}

// calculateResponse is the response body for a successful calculation.
type calculateResponse struct {
	Result float64 `json:"result"`
}

// errorResponse is the response body for a failed request.
type errorResponse struct {
	Error string `json:"error"`
}
