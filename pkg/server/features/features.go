// Package features provides an interface to implement custom cloud features.
// This provider is used by the [server] package.
//
// This package also provides a default feature provider implementation
package features

import (
	"net/http"
)

// FeatureProvider provides specific cloud features.
// If manages the instanciation and closing the features.
type FeatureProvider interface {
	New(serviceName string, onError func(err error)) error
	Close() error

	Report(err error, req *http.Request)
}
