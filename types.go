package callback

type Data struct {
	Point    string    `json:"point"`
	Success  bool      `json:"success"`
	Response *Response `json:"response"`
	Error    *Error    `json:"error"`
}

type ErrorInterface interface {
	IsError() bool
}

type Error struct {
	Code     int    `json:"code"`
	Message  string `json:"message"`
	Critical bool   `json:"critical"`
}

func (r *Error) IsError() bool {
	return true
}

type ResponseInterface interface {
	IsResponse() bool
}

type Response struct {
	Data []byte `json:"data"`
}

func (r *Response) IsResponse() bool {
	return true
}
