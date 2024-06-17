package types

type HttpJobPayload struct {
	QueryParams map[string]string `json:"query,omitempty"`
	Body        interface{}       `json:"body,omitempty"`
	Headers     map[string]string `json:"headers,omitempty"`
}
