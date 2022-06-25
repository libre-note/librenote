package response

import "net/http"

type Response struct {
	Success  bool        `json:"success"`
	Message  string      `json:"message,omitempty"`
	Errors   interface{} `json:"errors,omitempty"`
	Count    *int        `json:"count,omitempty"`
	PageSize *int        `json:"page_size,omitempty"`
	Previous *int        `json:"previous,omitempty"`
	Next     *int        `json:"next,omitempty"`
	Current  *int        `json:"current,omitempty"`
	Results  interface{} `json:"result,omitempty"`
}

func RespondSuccess(msg string, results interface{}) (int, Response) {
	return http.StatusOK, Response{
		Success: true,
		Message: msg,
		Results: results,
	}
}
