package config

import "fmt"

// The runtime Environment enumeration.
// Environment values are meant to be set only in this package.
type Environment string

// onCloud says if the current environment is supposed to be a cloud environment
func (e Environment) OnCloud() bool {
	return e != Development
}

const (
	developmentFull  = "development"
	developmentShort = "dev"

	productionFull     = "production"
	productionShort    = "prod"
	productionShortAlt = "prd"
)

// The current supported runtime Environment values.
// Their values are used when Environment is cast to a string (can be useful for
// logging / debugging).
const (
	Development Environment = developmentShort
	Production  Environment = productionShort
)

// EnvironmentMapping stores the values that can be user-provided.
// It allows to retrieve an Environment type from an input string value.
var EnvironmentMapping = map[Environment][]string{
	Development: {
		developmentFull,
		developmentShort,
	},
	Production: {
		productionFull,
		productionShort,
		productionShortAlt,
	},
}

// Lookup in EnvironmentMapping to retrieve an Environment type from an input
// string value.
func getEnvironment(environmentValue string) (Environment, error) {
	for env, values := range EnvironmentMapping {
		for _, value := range values {
			if environmentValue == value {
				return env, nil
			}
		}
	}

	return "", fmt.Errorf("no suitable environment found for %s", environmentValue)
}
