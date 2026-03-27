package response

type Envelope struct {
	Code      int         `json:"code"`
	Message   string      `json:"message"`
	Data      interface{} `json:"data,omitempty"`
	RequestID string      `json:"request_id,omitempty"`
}

func Success(data interface{}) Envelope {
	return Envelope{
		Code:    0,
		Message: "ok",
		Data:    data,
	}
}

func Error(code int, message string) Envelope {
	return Envelope{
		Code:    code,
		Message: message,
	}
}
