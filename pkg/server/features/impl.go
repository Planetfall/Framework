package features

import (
	"context"
	"fmt"
	"net/http"

	"cloud.google.com/go/compute/metadata"
	"cloud.google.com/go/errorreporting"
	secretmanager "cloud.google.com/go/secretmanager/apiv1"
)

type FeatureProviderImpl struct {
	metadataClient *metadata.Client       // the client access the Cloud project metadatas
	secretManager  *secretmanager.Client  // the client to access secrets
	errorReporting *errorreporting.Client // the client to report errors
}

// New initialize the provider features
func (f *FeatureProviderImpl) New(
	serviceName string, onError func(err error)) error {

	ctx := context.Background()

	// metadata client
	metadataClient := metadata.NewClient(nil)
	projectId, err := metadataClient.ProjectID()
	if err != nil {
		return fmt.Errorf("metadataClient.ProjectID: %v", err)
	}

	// secret manager
	secretManager, err := secretmanager.NewClient(ctx)
	if err != nil {
		return fmt.Errorf("secretmanager.NewClient: %v", err)
	}

	// error reporting
	errorReporting, err := errorreporting.NewClient(
		ctx, projectId, errorreporting.Config{
			ServiceName: serviceName,
			OnError:     onError,
		})
	if err != nil {
		return fmt.Errorf("errorreporting.NewClient: %v", err)
	}

	// set the feature in the provider
	f.metadataClient = metadataClient
	f.errorReporting = errorReporting
	f.secretManager = secretManager

	return nil
}

func (f *FeatureProviderImpl) Close() error {
	if err := f.secretManager.Close(); err != nil {
		return fmt.Errorf("secretManager.Close: %v", err)
	}

	if err := f.errorReporting.Close(); err != nil {
		return fmt.Errorf("errorReporting.Close: %v", err)
	}

	return nil
}

func (f *FeatureProviderImpl) Report(err error, req *http.Request) {
	f.errorReporting.Report(errorreporting.Entry{
		Error: err,
		Req:   req,
	})
}
