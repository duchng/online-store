package environment

type Environment string

const (
	Local       Environment = "local"
	Development Environment = "development"
	Production  Environment = "production"
	Uat         Environment = "uat"
)

func (e Environment) IsProduction() bool {
	return e == Production
}

func (e Environment) IsLocal() bool {
	return e == Local
}

const EnvContextKey = "env"

var Current = Production
