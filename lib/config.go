package lib

import (
	"encoding/json"
	"log"
	"os"
)

type Config struct {
	TerraformFile string `json:"terraformfile"`
	PlanFile      string `json:"planfile"`
	Layers        string `json:layers"`
}

func (c *Config) LoadConfig(path string) (*Config, error) {

	_, err := os.Stat(path)
	if err != nil {
		log.Fatalf("Config file does not exist: %v", err)
		return nil, err
	}

	file, err := os.Open(path)
	if err != nil {
		log.Fatalf("Error opening config file: %v", err)
		return nil, err
	}
	defer file.Close()

	err = json.NewDecoder(file).Decode(c)
	if err != nil {
		log.Fatalf("Error decoding config file: %v", err)
		return nil, err
	}

	return c, nil
}

func NewConfig() (*Config, error) {
	config := &Config{}

	if os.Getenv("GOGH_CONF") == "" {
		log.Fatal("GOGH_CONF environment variable is not set.")
	}

	_, err := config.LoadConfig(os.Getenv("GOGH_CONF"))
	if err != nil {
		return nil, err
	}

	return config, nil
}
