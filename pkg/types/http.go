package types

type HttpJobPayload struct {
	QueryParams map[string]string `json:"query,omitempty"`
	Body        string            `json:"body,omitempty"`
	Headers     map[string]string `json:"headers,omitempty"`
}
