// Package server provides an helper type that provides access to cloud
// features, error handling, logging and configuration management.
package server

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/planetfall/framework/pkg/config"
	"github.com/planetfall/framework/pkg/server/features"
)

// Server holds cloud features clients, a logger and the configuration.
type Server struct {
	Logger *log.Logger // the logger

	cfg config.Config // the configuration

	fp features.FeatureProvider // the provider for cloud features

}

// Raise logs the error and report it using the ErrorReporting cloud feature.
// The reporting is only available when in a Clouc environment.
func (s *Server) Raise(message string, err error, req *http.Request) {
	err = fmt.Errorf("%s: %v", message, err)
	s.Logger.Println(err)

	if s.cfg.Environment().OnCloud() {
		s.fp.Report(err, req)
	}
}

// Close terminates the server clients. If the server is not running on a
// Cloud environment, it does nothing and returns nil.
func (s *Server) Close() error {

	s.Logger.Printf("stopping the server")

	if s.cfg.Environment().OnCloud() {

		s.Logger.Printf("stopping onCloud features")

		if err := s.fp.Close(); err != nil {
			return fmt.Errorf("FeatureProvider.Close: %v", err)
		}
	}

	return nil
}

// NewServer creates a new server.
// It setup a dedicated logger using the serviceName parameter.
// A custom feature provider can be given. If 0, or more than one is given,
// it will fallback to the default provider.
// The default provider includes a metadata client, the error reporting and the
// secret manager.
func NewServer(
	cfg config.Config,
	serviceName string,
	featureProvider ...features.FeatureProvider,
) (*Server, error) {

	// setup logging
	environment := cfg.Environment()
	envPrefix := strings.ToUpper(environment.String())
	serviceNamePrefix := strings.ToLower(serviceName)
	logPrefix := fmt.Sprintf("[%s] - %s - ", envPrefix, serviceNamePrefix)
	logger := log.New(os.Stdout, logPrefix, log.Ldate|log.Ltime)

	// setup server features
	logger.Printf("setting up the server for %s", environment)

	var fp features.FeatureProvider
	if environment.OnCloud() {

		logger.Printf("starting onCloud features")

		fp = new(features.FeatureProviderImpl)
		if len(featureProvider) == 1 {
			fp = featureProvider[0]
		}

		onError := func(err error) {
			logger.Printf("could not log error: %v", err)
		}
		err := fp.New(serviceName, onError)
		if err != nil {
			return nil, fmt.Errorf("featureProvider.New: %v", err)
		}
	}

	return &Server{
		cfg:    cfg,
		Logger: logger,

		fp: fp,
	}, nil
}
