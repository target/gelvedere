package model

// AdminConfig is a struct for the service configuration
type AdminConfig struct {
	GheKey    string            `json:"ghe_key"`
	GheSecret string            `json:"ghe_secret"`
	Port      string            `json:"port"`
	AdminSSH  string            `json:"admin_ssh_pubkey"`
	Size      string            `json:"size"`
	Image     string            `json:"image"`
	Region    string            `json:"region"`
	LogConfig map[string]string `json:"log_config"`
	EnvVars   map[string]string `json:"env_variables"`
}

// UserConfig is a struct for looking up team information
type UserConfig struct {
	Name    string `json:"name"`
	Admins  string `json:"admins"`
	Members string `json:"members"`
	Team    string `json:"team"`
}
