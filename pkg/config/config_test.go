package config_test

import (
	"testing"

	"os"

	"github.com/planetfall/framework/pkg/config"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
)

var osArgsSaved = os.Args
var configFileTest = "testdata/config.yaml"

func addCommandLineArguments(arguments ...string) {
	osArgsSaved = os.Args
	os.Args = append(os.Args, arguments...)
}

func resetCommandLine() {
	os.Args = osArgsSaved
}

func initEntries() []config.Entry {
	return []config.Entry{
		{
			Flag:         "flag",
			DefaultValue: "default",
			Description:  "description",
			EnvKey:       "KEY",
		},
	}
}

func TestNewConfig_emptyEntries(t *testing.T) {
	// given
	entries := initEntries()
	resetCommandLine()

	// when
	_, err := config.NewConfig(entries)

	// then
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "viper.ReadInConfig")

	environmentExpected := config.EnvironmentDefaultValue
	environmentActual := viper.GetString(config.EnvironmentFlag)
	assert.Equal(t, environmentExpected, environmentActual)

	configExpected := config.ConfigDefaultValue
	configActual := viper.GetString(config.ConfigFlag)
	assert.Equal(t, configExpected, configActual)
}

func TestNewConfig_configWithFlag(t *testing.T) {
	// given
	entries := initEntries()
	configGiven := configFileTest
	resetCommandLine()
	addCommandLineArguments("--config", configGiven)

	// when
	_, err := config.NewConfig(entries)

	// then
	assert.Nil(t, err)

	configExpected := configGiven
	configActual := viper.GetString(config.ConfigFlag)
	assert.Equal(t, configExpected, configActual)
}

func TestNewConfig_configWithEnv(t *testing.T) {
	// given
	entries := initEntries()
	resetCommandLine()
	configGiven := configFileTest
	os.Setenv(config.ConfigEnvKey, configGiven)

	// when
	_, err := config.NewConfig(entries)

	// then
	assert.Nil(t, err)

	configExpected := configGiven
	configActual := viper.GetString(config.ConfigFlag)
	assert.Equal(t, configExpected, configActual)
}

func TestNewConfig_configWithFlagAndEnv(t *testing.T) {
	// given
	entries := initEntries()

	configEnvGiven := "config.yaml"
	os.Setenv(config.ConfigEnvKey, configEnvGiven)

	configFlagGiven := configFileTest
	resetCommandLine()
	addCommandLineArguments("--config", configFlagGiven)

	// when
	_, err := config.NewConfig(entries)

	// then
	assert.Nil(t, err)

	configExpected := configFlagGiven
	configActual := viper.GetString(config.ConfigFlag)
	assert.Equal(t, configExpected, configActual)
}

func TestNewConfig_environmentWithFlag(t *testing.T) {

	// test for all mapped environment values
	for environment, environmentValues := range config.EnvironmentMapping {
		for _, environmentValue := range environmentValues {

			// given
			environmentValueGiven := environmentValue
			t.Logf("testing with environment value %s", environmentValueGiven)

			entries := initEntries()
			resetCommandLine()
			addCommandLineArguments(
				"--env", environmentValueGiven,
				"--config", configFileTest)

			// when
			c, err := config.NewConfig(entries)

			// then
			assert.Nil(t, err)

			// checking value in viper
			environmentValueActual := viper.GetString(config.EnvironmentFlag)
			environmentValueExpected := environmentValue
			assert.Equal(t, environmentValueExpected, environmentValueActual)

			// checking env in config object
			environmentActual := c.Environment()
			environmentExpected := environment
			assert.Equal(t, environmentActual, environmentExpected)
		}
	}
}

func FuzzNewConfig_environmentWithFlagInvalid(f *testing.F) {
	environmentValuesInvalid := []string{"devv", "prd ", "", " "}
	environmentValues := make([]string, 0)
	environmentValues = append(environmentValues, config.EnvironmentMapping[config.Development]...)
	environmentValues = append(environmentValues, config.EnvironmentMapping[config.Production]...)

	for _, tc := range environmentValuesInvalid {
		f.Add(tc)
	}

	f.Fuzz(func(t *testing.T, environmentValueGiven string) {
		// given
		entries := initEntries()
		resetCommandLine()
		addCommandLineArguments(
			"--env", environmentValueGiven,
			"--config", configFileTest)

		// when
		_, err := config.NewConfig(entries)

		// then

		// in case fuzz generates a correct environment value
		for _, environmentValue := range environmentValues {
			if environmentValue == environmentValueGiven {
				assert.Nil(t, err)
				return
			}
		}

		if err == nil {
			t.Errorf("Expected err to be not nil")
		}
		assert.NotNil(t, err)
		assert.Contains(t, err.Error(), "no suitable environment found for")
	})
}

func TestNewConfig_entryWithEnv(t *testing.T) {

	viper.Reset()
	// given
	entries := initEntries()
	clientIdFlag := "client-id"
	clientIdEnv := "CLIENT_ID"
	entries = append(entries, config.Entry{
		Flag:         clientIdFlag,
		DefaultValue: "",
		Description:  "",
		EnvKey:       clientIdEnv,
	})

	resetCommandLine()
	addCommandLineArguments(
		"--config", configFileTest)

	clientIdEnvGiven := "cliend_id_env"
	os.Setenv(clientIdEnv, clientIdEnvGiven)

	// when
	_, err := config.NewConfig(entries)

	// then
	assert.Nil(t, err)

	clientIdExpected := clientIdEnvGiven
	clientIdActual := viper.GetString(clientIdFlag)
	assert.Equal(t, clientIdExpected, clientIdActual)
}

func TestNewConfig_entryWithFlagAndEnv(t *testing.T) {
	// given
	entries := initEntries()
	clientIdFlag := "client-id"
	clientIdEnv := "CLIENT_ID"
	entries = append(entries, config.Entry{
		Flag:         clientIdFlag,
		DefaultValue: "",
		Description:  "",
		EnvKey:       clientIdEnv,
	})

	clientIdFlagGiven := "client_id_flag"
	resetCommandLine()
	addCommandLineArguments(
		"--config", configFileTest,
		"--client-id", clientIdFlagGiven)

	clientIdEnvGiven := "cliend_id_env"
	os.Setenv(clientIdEnv, clientIdEnvGiven)

	// when
	_, err := config.NewConfig(entries)

	// then
	assert.Nil(t, err)

	clientIdExpected := clientIdFlagGiven
	clientIdActual := viper.GetString(clientIdFlag)
	assert.Equal(t, clientIdExpected, clientIdActual)
}
