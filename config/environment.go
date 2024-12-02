package config

type Environment int

const (
	Development Environment = iota
	Staging
	Production
)

func (e Environment) String() string {
	return [...]string{"dev", "stg", "prod"}[e]
}

func (e *Environment) init(environment string) {
	switch environment {
	case "":
		*e = Development
	case "stg":
		*e = Staging
	case "prod":
		*e = Production
	default:
		*e = Production
	}
}
