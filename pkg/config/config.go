package config

type AppConfig struct {
	GearmanServer string      `yaml:"gearman"`
	Jobs          []JobConfig `yaml:"jobs"`
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
	LogFile    string            `yaml:"log_gile"`
	WorkingDir string            `yaml:"working_dir"`
	Env        map[string]string `yaml:"env"`
	Timeout    string            `yaml:"timeout,omitempty"`
}

type HttpJobConfig struct {
	Url     string            `yaml:"url"`
	LogFile string            `yaml:"log_gile"`
	Headers map[string]string `yaml:"headers"`
	Method  string            `yaml:"method"`
	Timeout string            `yaml:"timeout,omitempty"`
}
