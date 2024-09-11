package types

type ShellJobPayload struct {
	EnvParams map[string]string `json:"env,omitempty"`
	LogFile   string            `json:"log_file,omitempty"`
}
