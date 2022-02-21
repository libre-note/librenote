package response

type Response struct {
	Success  bool        `json:"success"`
	Message  string      `json:"message,omitempty"`
	Count    *int        `json:"count,omitempty"`
	PageSize *int        `json:"page_size,omitempty"`
	Previous *int        `json:"previous,omitempty"`
	Next     *int        `json:"next,omitempty"`
	Current  *int        `json:"current,omitempty"`
	Results  interface{} `json:"result,omitempty"`
}
