package config

type AppConfig struct {
	GearmanServer string       `yaml:"gearman"`
	Server        ServerConfig `yaml:"server"`
	Jobs          []JobConfig  `yaml:"jobs"`
}

type ServerConfig struct {
	Addr string
}

type JobConfig struct {
	Concurrency int             `yaml:"concurrency"`
	Name        string          `yaml:"name"`
	Type        string          `yaml:"type"`
	ShellConfig *ShellJobConfig `yaml:"shell_config,omitempty"`
	HttpConfig  *HttpJobConfig  `yaml:"http_config,omitempty"`
	Timeout     string          `yaml:"timeout,omitempty"`
}

type ShellJobConfig struct {
	Command    string            `yaml:"command"`
	Args       []string          `yaml:"args"`
	LogFile    string            `yaml:"log_file"`
	WorkingDir string            `yaml:"working_dir"`
	Env        map[string]string `yaml:"env"`
	Timeout    string            `yaml:"timeout,omitempty"`
}

type HttpJobConfig struct {
	Url                string            `yaml:"url"`
	LogFile            string            `yaml:"log_file"`
	Headers            map[string]string `yaml:"headers"`
	Method             string            `yaml:"method"`
	Timeout            string            `yaml:"timeout,omitempty"`
	AllowedQueryParams []string          `yaml:"allowed_query,omitempty"`   // List of allowed query params to set
	AllowedHeaders     []string          `yaml:"allowed_headers,omitempty"` // List of allowed headers variable to use, if empty allow all
}
