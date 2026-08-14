package edgeadapter

import (
	"context"
	"net"
	"testing"

	devicecontrol "github.com/Zequent/zqnt-edge-sdk-go/gen/devicecontrol/contracts/proto"
	edgepb "github.com/Zequent/zqnt-edge-sdk-go/gen/edge/sdk/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

// fakeAdapter overrides only TakeOff — every other RPC falls through to
// UnimplementedEdgeAdapterServiceServer's NOT_IMPLEMENTED default, exactly the contract adapters
// are meant to rely on.
type fakeAdapter struct {
	edgepb.UnimplementedEdgeAdapterServiceServer
}

func (a *fakeAdapter) TakeOff(ctx context.Context, req *devicecontrol.CoordinateCommandRequest) (*devicecontrol.CommandResponse, error) {
	return Success(req.GetBase().GetTid()), nil
}

func TestServerRoutesAnOverriddenCommandAndDefaultsEverythingElseToUnimplemented(t *testing.T) {
	lis, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("listen: %v", err)
	}
	server := NewServer(&fakeAdapter{})
	go func() { _ = server.Serve(lis) }()
	defer server.Stop()

	conn, err := grpc.NewClient(lis.Addr().String(), grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		t.Fatalf("dial: %v", err)
	}
	defer conn.Close()
	client := edgepb.NewEdgeAdapterServiceClient(conn)

	resp, err := client.TakeOff(context.Background(), &devicecontrol.CoordinateCommandRequest{})
	if err != nil {
		t.Fatalf("TakeOff: %v", err)
	}
	if resp.GetHasErrors() {
		t.Fatalf("expected the overridden TakeOff to succeed, got an error response")
	}

	_, err = client.OpenCover(context.Background(), &devicecontrol.EmptyCommandRequest{})
	if err == nil {
		t.Fatalf("expected OpenCover (never overridden) to fail as Unimplemented")
	}
}
