package config

// Entry is a type that allows the config package to access configuration
// values.
type Entry struct {
	Flag         string // flag to parse from program argument
	DefaultValue string // default value in case no source provides a value
	Description  string // description used by the flag help command
	EnvKey       string // the environment variables that holds the value
}

// The fields for the config entry.
const (
	ConfigFlag         = "config"
	ConfigDefaultValue = "config/config.yaml"
	ConfigEnvKey       = "CONFIG"
)

// The fields for the runtime environment entry.
const (
	EnvironmentFlag         = "env"
	EnvironmentDefaultValue = productionFull
	EnvironmentEnvKey       = "ENV"
)

var (
	configFileEntry = Entry{
		Flag:         ConfigFlag,
		DefaultValue: ConfigDefaultValue,
		Description:  "the config file path",
		EnvKey:       ConfigEnvKey,
	}

	environmentEntry = Entry{
		Flag:         EnvironmentFlag,
		DefaultValue: EnvironmentDefaultValue,
		Description:  "the runtime environment",
		EnvKey:       EnvironmentEnvKey,
	}
)
