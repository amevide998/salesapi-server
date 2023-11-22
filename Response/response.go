package Response

type WebResponse[T any] struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
	Data    T      `json:"data"`
}

type WebErrorResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
	Error   string `json:"error"`
}

func NewWebErrorResponse(e string, msg ...string) *WebErrorResponse {
	defaultMsg := "something wrong"
	if len(msg) > 0 {
		defaultMsg = msg[0]
	}
	return &WebErrorResponse{
		Success: false,
		Message: defaultMsg,
		Error:   e,
	}
}

func NewWebResponse[T any](data T, msg ...string) *WebResponse[T] {
	defaultMsg := "success"
	if len(msg) > 0 {
		defaultMsg = msg[0]
	}
	return &WebResponse[T]{
		Success: true,
		Message: defaultMsg,
		Data:    data,
	}
}
