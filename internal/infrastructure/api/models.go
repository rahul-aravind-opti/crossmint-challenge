package api

// CreatePolyanetRequest represents the request body for creating a Polyanet
type CreatePolyanetRequest struct {
	Row        int    `json:"row"`
	Column     int    `json:"column"`
	CandidateID string `json:"candidateId"`
}

// CreateSoloonRequest represents the request body for creating a Soloon
type CreateSoloonRequest struct {
	Row        int    `json:"row"`
	Column     int    `json:"column"`
	Color      string `json:"color"`
	CandidateID string `json:"candidateId"`
}

// CreateComethRequest represents the request body for creating a Cometh
type CreateComethRequest struct {
	Row        int    `json:"row"`
	Column     int    `json:"column"`
	Direction  string `json:"direction"`
	CandidateID string `json:"candidateId"`
}

// DeleteRequest represents the request body for deleting an object
type DeleteRequest struct {
	Row        int    `json:"row"`
	Column     int    `json:"column"`
	CandidateID string `json:"candidateId"`
}

// ErrorResponse represents an error response from the API
type ErrorResponse struct {
	Error   string `json:"error,omitempty"`
	Message string `json:"message,omitempty"`
	Status  int    `json:"status,omitempty"`
}
