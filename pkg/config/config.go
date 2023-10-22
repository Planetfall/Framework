// Package config implements utility types for parsing and storing into [viper]
// multiple configuration values.
//
// The config package extract config entries from a configuration file, program
// arguments and runtime environment. It extract the input entries in a specific
// order. Those entry values are then stored using the [viper] package.
package config

import (
	"fmt"

	flag "github.com/spf13/pflag"
	"github.com/spf13/viper"
)

type Config interface {
	Environment() Environment
}

// Config is the type that holds the current runtime environment.
type ConfigImpl struct {
	environment Environment // Current provided runtime environment.
}

// Environment provides the current Config environment value.
func (c *ConfigImpl) Environment() Environment {
	return c.environment
}

// initEnv binds environment variables into [viper] using provided entries.
func initEnv(entries []Entry) error {
	for _, entry := range entries {

		// BindEnv only returns an error when no key is provided.
		// Binding a value from the environment remains optional.
		if err := viper.BindEnv(entry.Flag, entry.EnvKey); err != nil {
			return fmt.Errorf("viper.BindEnv: %v", err)
		}
	}

	return nil
}

// initFlags setup the program flags. It initialize the flags for each entry.
// Then, the actual program arguments are parsed, and the values binded to
// [viper].
func initFlags(entries []Entry) error {
	for _, f := range entries {
		if flag.Lookup(f.Flag) != nil {
			continue
		}

		flag.String(f.Flag, f.DefaultValue, f.Description)
	}

	flag.Parse()
	if err := viper.BindPFlags(flag.CommandLine); err != nil {
		return fmt.Errorf("viper.BindPFlags: %v", err)
	}

	return nil
}

// setConfigFile takes a configFile path and reads it using [viper].
func setConfigFile(configFile string) error {

	viper.SetConfigFile(configFile)
	if err := viper.ReadInConfig(); err != nil {
		return fmt.Errorf("viper.ReadInConfig(%s): %v", configFile, err)
	}
	return nil
}

// setDefaultValues initialize all [viper] default values for all entries.
func setDefaultValues(entries []Entry) {
	for _, entry := range entries {
		viper.SetDefault(entry.Flag, entry.DefaultValue)
		viper.SetDefault(entry.EnvKey, entry.DefaultValue)
	}
}

// NewConfig takes a slice of entries, and setup the [viper] configuration map.
// It looks for those entries in various sources, in a specific order:
//
//  1. Add to the provided entries two default entries.
//  2. Set the default values in [viper] for all entries.
//  3. Look in the environment variables and bind the values into [viper].
//  4. Parse the program arguments and bind the values into [viper].
//  5. Read the config file and inject the values into [viper].
//
// The two default entries are currently:
//   - ENV, which indicates the program environment.
//   - CONFIG, which indicates the config file path.
func NewConfig(entries []Entry) (Config, error) {

	entries = append(entries, configFileEntry)
	entries = append(entries, environmentEntry)

	setDefaultValues(entries)

	if err := initEnv(entries); err != nil {
		return nil, fmt.Errorf("initEnv: %v", err)
	}

	// flags overrides the env
	if err := initFlags(entries); err != nil {
		return nil, fmt.Errorf("initFlags: %v", err)
	}

	// set config file
	configFilePath := viper.GetString(ConfigFlag)
	if err := setConfigFile(configFilePath); err != nil {
		return nil, fmt.Errorf("config.setConfigFile: %v", err)
	}

	// set environment
	environmentString := viper.GetString(EnvironmentFlag)
	environment, err := getEnvironment(environmentString)
	if err != nil {
		return nil, fmt.Errorf("config.getEnvironment: %v", err)
	}

	return &ConfigImpl{
		environment: environment,
	}, nil
}
