package config

type Environment string

const (
	EnvDevelopment Environment = "development"
	EnvPreview     Environment = "preview"
	EnvStaging     Environment = "staging"
	EnvTest        Environment = "test"
	EnvProduction  Environment = "production"
)

func (e Environment) IsValid() bool {
	switch string(e) {
	case "development", "preview", "staging", "production", "test":
		return true
	default:
		return false
	}
}

func (e Environment) IsDevelopment() bool {
	return e == EnvDevelopment
}

func (e Environment) IsPreview() bool {
	return e == EnvPreview
}

func (e Environment) IsStaging() bool {
	return e == EnvStaging
}

func (e Environment) IsTest() bool {
	return e == EnvTest
}

func (e Environment) IsProduction() bool {
	return e == EnvProduction
}
