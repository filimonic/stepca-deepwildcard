package validator

type ValidatorOption func(*Validator) error

func WithConfig(config *Config) ValidatorOption {
	return func(v *Validator) error {
		v.config = config
		return nil
	}
}
