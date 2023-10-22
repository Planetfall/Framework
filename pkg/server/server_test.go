package server_test

import (
	"fmt"
	"testing"

	"github.com/planetfall/framework/pkg/config"
	"github.com/planetfall/framework/pkg/server"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

const methodEnvironment = "Environment"

type configMock struct {
	mock.Mock
}

func (c *configMock) Environment() config.Environment {
	args := c.Called()
	return args.Get(0).(config.Environment)
}

func TestNewServer_withDev(t *testing.T) {
	// given
	var cfgGiven = &configMock{}
	serviceGiven := "service-name"

	// when
	cfgGiven.On(methodEnvironment).Return(config.Development)
	s, err := server.NewServer(cfgGiven, serviceGiven)

	// then
	assert.Nil(t, err)
	assert.NotNil(t, s)

	assert.Contains(t, s.Logger.Prefix(), serviceGiven)
	cfgGiven.AssertNumberOfCalls(t, methodEnvironment, 1)
}

func TestRaise_withDev(t *testing.T) {
	// given
	var cfgGiven = &configMock{}
	serviceGiven := "service-name"
	errorGiven := fmt.Errorf("test error")
	messageGiven := "test error message"

	cfgGiven.On(methodEnvironment).Return(config.Development)
	s, err := server.NewServer(cfgGiven, serviceGiven)
	assert.Nil(t, err)
	assert.NotNil(t, s)

	// when
	s.Raise(messageGiven, errorGiven, nil)
	cfgGiven.AssertNumberOfCalls(t, methodEnvironment, 2)
}

func TestRaise_withDevAndErrorNil(t *testing.T) {
	// given
	var cfgGiven = &configMock{}
	serviceGiven := "service-name"
	var errorGiven error = nil
	messageGiven := "test error message"

	cfgGiven.On(methodEnvironment).Return(config.Development)
	s, err := server.NewServer(cfgGiven, serviceGiven)
	assert.Nil(t, err)
	assert.NotNil(t, s)

	// when
	s.Raise(messageGiven, errorGiven, nil)
	cfgGiven.AssertNumberOfCalls(t, methodEnvironment, 2)
}

func TestClose_withDev(t *testing.T) {
	// given
	var cfgGiven = &configMock{}
	serviceGiven := "service-name"

	cfgGiven.On(methodEnvironment).Return(config.Development)
	s, err := server.NewServer(cfgGiven, serviceGiven)
	assert.Nil(t, err)
	assert.NotNil(t, s)

	// when
	err = s.Close()
	assert.Nil(t, err)
	cfgGiven.AssertNumberOfCalls(t, methodEnvironment, 2)
}
