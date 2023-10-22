package server_test

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/planetfall/framework/pkg/config"
	"github.com/planetfall/framework/pkg/server"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

const methodEnvironment = "Environment"

// config
type configMock struct {
	mock.Mock
}

func (c *configMock) Environment() config.Environment {
	args := c.Called()
	return args.Get(0).(config.Environment)
}

// feature provider
type featureProviderMock struct {
	mock.Mock

	onError func(err error)
}

func (m *featureProviderMock) New(
	serviceName string, onError func(err error)) error {

	m.onError = onError

	args := m.Called(serviceName)
	return args.Error(0)
}

func (m *featureProviderMock) Close() error {
	args := m.Called()
	return args.Error(0)
}

func (m *featureProviderMock) Report(
	err error, req *http.Request) {

	m.Called(err)
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
	cfgGiven.AssertExpectations(t)
}

func TestNewServer_withPrd(t *testing.T) {
	// given
	var cfgGiven = &configMock{}
	serviceGiven := "service-name"
	fpGiven := &featureProviderMock{}

	// when
	cfgGiven.On(methodEnvironment).Return(config.Production)
	fpGiven.On("New", serviceGiven).Return(nil)
	s, err := server.NewServer(
		cfgGiven, serviceGiven, fpGiven)

	// then
	assert.Nil(t, err)
	assert.NotNil(t, s)

	cfgGiven.AssertExpectations(t)
	fpGiven.AssertExpectations(t)
}

func TestNewServer_withPrd_shouldFail(t *testing.T) {
	// given
	var cfgGiven = &configMock{}
	serviceGiven := "service-name"
	fpGiven := &featureProviderMock{}
	errorGiven := fmt.Errorf("new provider failed")

	// when
	cfgGiven.On(methodEnvironment).Return(config.Production)
	fpGiven.On("New", serviceGiven).Return(errorGiven)
	s, err := server.NewServer(
		cfgGiven, serviceGiven, fpGiven)

	// then
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), errorGiven.Error())
	assert.Nil(t, s)

	cfgGiven.AssertExpectations(t)
	fpGiven.AssertExpectations(t)
}

func TestRaise_withDev(t *testing.T) {
	// given
	var cfgGiven = &configMock{}
	serviceGiven := "service-name"

	errorGiven := fmt.Errorf("test error")
	messageGiven := "test error message"

	// when
	cfgGiven.On(methodEnvironment).Return(config.Development)
	s, err := server.NewServer(cfgGiven, serviceGiven)
	assert.Nil(t, err)
	assert.NotNil(t, s)

	s.Raise(messageGiven, errorGiven, nil)

	// then
	cfgGiven.AssertExpectations(t)
}

func TestRaise_withPrd(t *testing.T) {
	// given
	var cfgGiven = &configMock{}
	serviceGiven := "service-name"
	fpGiven := &featureProviderMock{}

	errorGiven := fmt.Errorf("test error")
	messageGiven := "test error message"

	// when
	cfgGiven.On(methodEnvironment).Return(config.Production)
	fpGiven.On("New", serviceGiven).Return(nil)
	fpGiven.On("Report", mock.Anything).Return()
	s, err := server.NewServer(
		cfgGiven, serviceGiven, fpGiven)

	assert.Nil(t, err)
	assert.NotNil(t, s)

	s.Raise(messageGiven, errorGiven, nil)

	// then
	fpGiven.AssertExpectations(t)
	cfgGiven.AssertExpectations(t)
}

func TestRaise_withDevAndErrorNil(t *testing.T) {
	// given
	var cfgGiven = &configMock{}
	serviceGiven := "service-name"
	var errorGiven error = nil
	messageGiven := "test error message"

	// when
	cfgGiven.On(methodEnvironment).Return(config.Development)
	s, err := server.NewServer(cfgGiven, serviceGiven)
	assert.Nil(t, err)
	assert.NotNil(t, s)

	s.Raise(messageGiven, errorGiven, nil)

	// then
	cfgGiven.AssertExpectations(t)
}

func TestClose_withDev(t *testing.T) {
	// given
	var cfgGiven = &configMock{}
	serviceGiven := "service-name"

	// when
	cfgGiven.On(methodEnvironment).Return(config.Development)
	s, err := server.NewServer(cfgGiven, serviceGiven)
	assert.Nil(t, err)
	assert.NotNil(t, s)

	err = s.Close()

	// then
	assert.Nil(t, err)
	cfgGiven.AssertExpectations(t)
}

func TestClose_withPrd(t *testing.T) {
	// given
	var cfgGiven = &configMock{}
	serviceGiven := "service-name"
	fpGiven := &featureProviderMock{}

	// when
	cfgGiven.On(methodEnvironment).Return(config.Production)
	fpGiven.On("New", serviceGiven).Return(nil)
	fpGiven.On("Close").Return(nil)
	s, err := server.NewServer(
		cfgGiven, serviceGiven, fpGiven)

	assert.Nil(t, err)
	assert.NotNil(t, s)

	err = s.Close()

	// then
	assert.Nil(t, err)
	fpGiven.AssertExpectations(t)
	cfgGiven.AssertExpectations(t)
}

func TestClose_withPrd_shouldFail(t *testing.T) {
	// given
	var cfgGiven = &configMock{}
	serviceGiven := "service-name"
	fpGiven := &featureProviderMock{}
	errorGiven := fmt.Errorf("error closing")

	// when
	cfgGiven.On(methodEnvironment).Return(config.Production)
	fpGiven.On("New", serviceGiven).Return(nil)
	fpGiven.On("Close").Return(errorGiven)
	s, err := server.NewServer(
		cfgGiven, serviceGiven, fpGiven)

	assert.Nil(t, err)
	assert.NotNil(t, s)

	err = s.Close()

	// then
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), errorGiven.Error())
	fpGiven.AssertExpectations(t)
	cfgGiven.AssertExpectations(t)
}
