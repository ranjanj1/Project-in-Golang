package grpc

import (
	"Project2/config"
	"fmt"
	"net"

	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	grpc_recovery "github.com/grpc-ecosystem/go-grpc-middleware/recovery"
	"google.golang.org/grpc"
	grpcTrace "gopkg.in/DataDog/dd-trace-go.v1/contrib/google.golang.org/grpc"
)

func ProvideListener(config config.ServerConfig) (net.Listener, error) {
	return net.Listen("tcp", fmt.Sprintf("%s:%d", config.Host, config.Port))
}

func ProvideGrpcServerOptions(config config.ServerConfig) (opts []grpc.ServerOption) {
	unaryInterceptors := grpc.UnaryInterceptor(grpc_middleware.ChainUnaryServer(
		grpcTrace.UnaryServerInterceptor(grpcTrace.WithServiceName(config.ServiceName)),
		grpc_recovery.UnaryServerInterceptor()))
	opts = append(opts, unaryInterceptors)
	return
}

func ProvideServer(opts ...grpc.ServerOption) *grpc.Server {
	srv := grpc.NewServer(opts...)
	return srv
}
