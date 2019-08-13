// THIS FILE IS AUTO GENERATED BY GK-CLI DO NOT EDIT!!
package grpc

import (
	endpoint "dogfooter-control/control/pkg/endpoint"
	pb "dogfooter-control/control/pkg/grpc/pb"
	grpc "github.com/go-kit/kit/transport/grpc"
)

// NewGRPCServer makes a set of endpoints available as a gRPC AddServer
type grpcServer struct {
	api  grpc.Handler
	root grpc.Handler
	file grpc.Handler
}

func NewGRPCServer(endpoints endpoint.Endpoints, options map[string][]grpc.ServerOption) pb.ControlServer {
	return &grpcServer{
		api:  makeApiHandler(endpoints, options["Api"]),
		file: makeFileHandler(endpoints, options["File"]),
		root: makeRootHandler(endpoints, options["Root"]),
	}
}