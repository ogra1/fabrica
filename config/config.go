package config

// Default settings
const (
	DefaultPort = "8000"
)

// Settings defines the application configuration
type Settings struct {
	Port string
}

// DefaultArgs checks the environment variables
func DefaultArgs() *Settings {
	return &Settings{
		Port: DefaultPort,
	}
}

// ReadParameters fetches the store config parameters
func ReadParameters() *Settings {
	config := DefaultArgs()
	return config
}
