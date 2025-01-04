package environment

import (
	"fmt"
	"strings"
)

type Environment int

type EnvironmentType string

const (
	Development Environment = iota
	Staging
	Production
)

const (
	DevEnv  EnvironmentType = "dev"
	StgEnv  EnvironmentType = "stg"
	ProdEnv EnvironmentType = "prod"
)

var (
	environmentStrings = map[Environment]EnvironmentType{
		Development: DevEnv,
		Staging:     StgEnv,
		Production:  ProdEnv,
	}

	environmentValues = map[EnvironmentType]Environment{
		DevEnv:  Development,
		StgEnv:  Staging,
		ProdEnv: Production,
	}
)

func (e Environment) String() string {
	if env, ok := environmentStrings[e]; ok {
		return string(env)
	}
	return string(ProdEnv)
}

func ParseEnvironment(env string) Environment {
	normalizedEnv := EnvironmentType(strings.ToLower(strings.TrimSpace(env)))
	if value, ok := environmentValues[normalizedEnv]; ok {
		return value
	}
	return Production
}

func (e Environment) Validate() error {
	if _, ok := environmentStrings[e]; !ok {
		return fmt.Errorf("invalid environment value: %d", e)
	}
	return nil
}

func (e Environment) IsDevelopment() bool {
	return e == Development
}

func (e Environment) IsStaging() bool {
	return e == Staging
}

func (e Environment) IsProduction() bool {
	return e == Production
}

type EnvironmentConfig struct {
	Current Environment
}

func NewEnvironmentConfig(envStr string) *EnvironmentConfig {
	return &EnvironmentConfig{
		Current: ParseEnvironment(envStr),
	}
}

func (ec *EnvironmentConfig) String() string {
	return ec.Current.String()
}

func (ec *EnvironmentConfig) IsValid() bool {
	return ec.Current.Validate() == nil
}
