package types

type HttpJobPayload struct {
	UrlParams   map[string]string `json:"url,omitempty"`
	QueryParams map[string]string `json:"query,omitempty"`
	Body        interface{}       `json:"body,omitempty"`
	Headers     map[string]string `json:"headers,omitempty"`
}
