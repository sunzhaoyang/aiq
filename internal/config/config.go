package config

// Config represents the application configuration
type Config struct {
	LLM LLMConfig `yaml:"llm"`
}

// LLMConfig represents LLM provider configuration
type LLMConfig struct {
	URL    string `yaml:"url"`
	APIKey string `yaml:"api_key"`
	Model  string `yaml:"model"`
}

// NewConfig creates a new empty configuration
func NewConfig() *Config {
	return &Config{
		LLM: LLMConfig{},
	}
}

// IsEmpty checks if the configuration is empty (first run)
func (c *Config) IsEmpty() bool {
	return c.LLM.URL == "" || c.LLM.APIKey == "" || c.LLM.Model == ""
}
