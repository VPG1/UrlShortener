package handlers

type ResponseError struct {
	Error string `json:"error"`
}

func NewResponseError(error string) ResponseError {
	return ResponseError{Error: error}
}

type SuccessResponse struct {
	Status string `json:"status"`
}

func NewSuccessResponse(status string) SuccessResponse {
	return SuccessResponse{Status: status}
}
