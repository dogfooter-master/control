// THIS FILE IS AUTO GENERATED BY GK-CLI DO NOT EDIT!!
package endpoint

import (
	service "dogfooter-control/control/pkg/service"
	endpoint "github.com/go-kit/kit/endpoint"
)

// Endpoints collects all of the endpoints that compose a profile service. It's
// meant to be used as a helper struct, to collect all of the endpoints into a
// single parameter.
type Endpoints struct {
	ApiEndpoint  endpoint.Endpoint
	RootEndpoint endpoint.Endpoint
	FileEndpoint endpoint.Endpoint
}

// New returns a Endpoints struct that wraps the provided service, and wires in all of the
// expected endpoint middlewares
func New(s service.ControlService, mdw map[string][]endpoint.Middleware) Endpoints {
	eps := Endpoints{
		ApiEndpoint:  MakeApiEndpoint(s),
		FileEndpoint: MakeFileEndpoint(s),
		RootEndpoint: MakeRootEndpoint(s),
	}
	for _, m := range mdw["Api"] {
		eps.ApiEndpoint = m(eps.ApiEndpoint)
	}
	for _, m := range mdw["Root"] {
		eps.RootEndpoint = m(eps.RootEndpoint)
	}
	for _, m := range mdw["File"] {
		eps.FileEndpoint = m(eps.FileEndpoint)
	}
	return eps
}
