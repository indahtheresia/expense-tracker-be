package dto

type Response struct {
	Success bool `json:"success"`
	Error   any  `json:"error,omitempty"`
	Data    any  `json:"data,omitempty"`
}

type ErrorType struct {
	ErrorStr string `json:"error"`
}

type CustomError struct {
	ErrorStr    string `json:"error,omitempty"`
	InternalErr string `json:"internal_error,omitempty"`
	Status      int    `json:"response_status,omitempty"`
}

func (err CustomError) Error() string {
	return err.InternalErr
}
