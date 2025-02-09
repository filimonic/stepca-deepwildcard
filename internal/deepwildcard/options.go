package deepwildcard

import (
	"fmt"
	"log"
	"os"

	"gopkg.in/yaml.v3"
)

type DeepHookOption func(*dhServer) error

func WithConfig(config *Config) DeepHookOption {
	return func(v *dhServer) error {
		v.config = config
		return nil
	}
}

func WithConfigFile(file string) DeepHookOption {
	return func(v *dhServer) error {
		data, err := os.ReadFile(file)
		if err != nil {
			return fmt.Errorf("failed to read file \"%s\": %w", file, err)
		}
		config := &Config{}
		err = yaml.Unmarshal(data, config)
		if err != nil {
			return fmt.Errorf("failed to parse file \"%s\": %w", file, err)
		}
		v.config = config
		return nil
	}
}

func WithLogger(logger *log.Logger) DeepHookOption {
	return func(v *dhServer) error {
		v.logger = logger
		return nil
	}
}
