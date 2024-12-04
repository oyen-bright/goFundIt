package environment

type Environment int

const (
	Development Environment = iota
	Staging
	Production
)

func (e Environment) String() string {
	return [...]string{"dev", "stg", "prod"}[e]
}

func (e *Environment) New(environment string) {
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

func (e Environment) IsDevelopment() bool {
	return e == Development
}

func (e Environment) IsStaging() bool {
	return e == Staging
}

func (e Environment) IsProduction() bool {
	return e == Production
}
