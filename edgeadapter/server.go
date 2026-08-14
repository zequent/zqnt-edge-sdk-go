package edgeadapter

import (
	"context"
	"net"

	edgepb "github.com/Zequent/zqnt-edge-sdk-go/gen/edge/sdk/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

// NewServer builds a *grpc.Server with impl registered as the EdgeAdapterService, reflection
// enabled (so grpcurl/console tooling can introspect it). An adapter only needs to implement the
// commands its hardware actually supports — embed edgepb.UnimplementedEdgeAdapterServiceServer in
// your own type to get NOT_IMPLEMENTED for everything else for free, the same convention every
// other language's Edge SDK follows ("only the commands a device supports need to be overridden").
//
//	type MyAdapter struct {
//		edgepb.UnimplementedEdgeAdapterServiceServer
//	}
//	func (a *MyAdapter) TakeOff(ctx context.Context, req *devicecontrol.CoordinateCommandRequest) (*devicecontrol.CommandResponse, error) {
//		...
//	}
func NewServer(impl edgepb.EdgeAdapterServiceServer) *grpc.Server {
	server := grpc.NewServer()
	edgepb.RegisterEdgeAdapterServiceServer(server, impl)
	reflection.Register(server)
	return server
}

// Serve listens on addr (e.g. ":9001") and serves impl as the EdgeAdapterService, blocking until
// ctx is cancelled (triggering a graceful stop) or the listener fails.
func Serve(ctx context.Context, addr string, impl edgepb.EdgeAdapterServiceServer) error {
	lis, err := net.Listen("tcp", addr)
	if err != nil {
		return err
	}
	server := NewServer(impl)

	errCh := make(chan error, 1)
	go func() { errCh <- server.Serve(lis) }()

	select {
	case <-ctx.Done():
		server.GracefulStop()
		return ctx.Err()
	case err := <-errCh:
		return err
	}
}
