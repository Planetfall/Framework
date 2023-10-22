package server

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	"cloud.google.com/go/compute/metadata"
	"cloud.google.com/go/errorreporting"
	secretmanager "cloud.google.com/go/secretmanager/apiv1"
	"github.com/planetfall/framework/pkg/config"
)

type Server struct {
	Logger *log.Logger

	cfg *config.Config

	metadataClient *metadata.Client
	secretManager  *secretmanager.Client
	errorReporting *errorreporting.Client
}

func (s *Server) Raise(message string, err error, req *http.Request) {
	err = fmt.Errorf("%s: %v", message, err)
	s.Logger.Println(err)

	if s.cfg.Environment().OnCloud() {
		s.errorReporting.Report(errorreporting.Entry{
			Error: err,
			Req:   req,
		})
	}
}

func (s *Server) Close() error {

	s.Logger.Printf("stopping the server")

	if s.cfg.Environment().OnCloud() {

		s.Logger.Printf("stopping onCloud features")

		if err := s.secretManager.Close(); err != nil {
			return fmt.Errorf("secretManager.Close: %v", err)
		}

		if err := s.errorReporting.Close(); err != nil {
			return fmt.Errorf("errorReporting.Close: %v", err)
		}
	}

	return nil
}

func NewServer(cfg *config.Config, serviceName string) (*Server, error) {
	// setup logging
	environment := cfg.Environment()
	envPrefix := strings.ToUpper(string(environment))
	serviceNamePrefix := strings.ToLower(serviceName)
	logPrefix := fmt.Sprintf("[%s] - %s - ", envPrefix, serviceNamePrefix)
	logger := log.New(os.Stdout, logPrefix, log.Ldate|log.Ltime|log.Lshortfile)

	// setup server features
	logger.Printf("setting up the server for %s", environment)

	var metadataClient *metadata.Client
	var secretManager *secretmanager.Client
	var errorReporting *errorreporting.Client

	if environment.OnCloud() {

		logger.Printf("starting onCloud features")
		ctx := context.Background()

		// metadata client
		metadataClient = metadata.NewClient(nil)
		projectId, err := metadataClient.ProjectID()
		if err != nil {
			return nil, fmt.Errorf("metadataClient.ProjectID: %v", err)
		}

		// secret manager
		secretManager, err = secretmanager.NewClient(ctx)
		if err != nil {
			return nil, fmt.Errorf("secretmanager.NewClient: %v", err)
		}

		// error reporting
		errorReporting, err = errorreporting.NewClient(ctx, projectId, errorreporting.Config{
			ServiceName: serviceName,
			OnError: func(err error) {
				logger.Printf("could not log error: %v", err)
			},
		})
		if err != nil {
			return nil, fmt.Errorf("errorreporting.NewClient: %v", err)
		}
	}

	return &Server{
		cfg:    cfg,
		Logger: logger,

		metadataClient: metadataClient,
		secretManager:  secretManager,
		errorReporting: errorReporting,
	}, nil
}
