package xds

import (
	"fmt"
	"net"

	"github.com/envoyproxy/go-control-plane/envoy/api/v2"
	"github.com/envoyproxy/go-control-plane/envoy/api/v2/core"
	envoyv2 "github.com/envoyproxy/go-control-plane/envoy/service/discovery/v2"
	envoycache "github.com/envoyproxy/go-control-plane/pkg/cache"
	xds "github.com/envoyproxy/go-control-plane/pkg/server"
	"github.com/grpc-ecosystem/go-grpc-middleware"
	"github.com/grpc-ecosystem/go-grpc-middleware/logging/zap"
	"github.com/grpc-ecosystem/go-grpc-middleware/tags"
	"github.com/solo-io/gloo/pkg/log"
	"go.uber.org/zap"
	"google.golang.org/grpc"
)

// For now we're only running one envoy instance
const NodeKey = string("gloo-envoy")

type hasher struct{}

func (h hasher) ID(node *core.Node) string {
	return NodeKey
}

type logger struct{}

func (*logger) Infof(format string, args ...interface{}) {
	log.GreyPrintf(format, args...)
}
func (*logger) Errorf(format string, args ...interface{}) {
	log.Warnf(format, args...)
}

func RunXDS(port int) (envoycache.SnapshotCache, *grpc.Server, error) {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		return nil, nil, fmt.Errorf("failed to listen: %v", err)
	}
	envoyCache := envoycache.NewSnapshotCache(true, hasher{}, &logger{})
	grpcServer := grpc.NewServer(grpc.StreamInterceptor(
		grpc_middleware.ChainStreamServer(
			grpc_ctxtags.StreamServerInterceptor(),
			grpc_zap.StreamServerInterceptor(zap.NewNop()),
			func(srv interface{}, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
				log.Debugf("xds: %v", info.FullMethod)
				return handler(srv, ss)
			},
		)),
	)
	xdsServer := xds.NewServer(envoyCache, nil)
	envoyv2.RegisterAggregatedDiscoveryServiceServer(grpcServer, xdsServer)
	v2.RegisterEndpointDiscoveryServiceServer(grpcServer, xdsServer)
	v2.RegisterClusterDiscoveryServiceServer(grpcServer, xdsServer)
	v2.RegisterRouteDiscoveryServiceServer(grpcServer, xdsServer)
	v2.RegisterListenerDiscoveryServiceServer(grpcServer, xdsServer)

	go func() {
		log.Debugf("xDS server listening on %s", port)
		if err = grpcServer.Serve(lis); err != nil {
			log.Fatalf("failed to serve grpc: %v", err)
		}
	}()
	return envoyCache, grpcServer, nil
}
