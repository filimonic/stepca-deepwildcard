package deepwildcard

import "deepwildcard/internal/deepwildcard/validator"

type Config struct {
	ListenAddr      string           `yaml:"address"`
	ValidatorConfig validator.Config `yaml:"validator"`
}
