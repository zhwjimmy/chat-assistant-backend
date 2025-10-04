package elasticsearch

import "time"

// Config holds Elasticsearch configuration
type Config struct {
	Hosts    []string      `mapstructure:"hosts"`
	Username string        `mapstructure:"username"`
	Password string        `mapstructure:"password"`
	Timeout  time.Duration `mapstructure:"timeout"`
	Index    IndexConfig   `mapstructure:"index"`
}

// IndexConfig holds index-specific configuration
type IndexConfig struct {
	Conversations string `mapstructure:"conversations"`
	Messages      string `mapstructure:"messages"`
}

// DefaultConfig returns default Elasticsearch configuration
func DefaultConfig() *Config {
	return &Config{
		Hosts:    []string{"http://localhost:9200"},
		Username: "",
		Password: "",
		Timeout:  30 * time.Second,
		Index: IndexConfig{
			Conversations: "conversations",
			Messages:      "messages",
		},
	}
}
