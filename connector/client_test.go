package connector

import (
	"context"
	"net"
	"testing"

	asset "github.com/Zequent/zqnt-edge-sdk-go/gen/common/asset/proto"
	base "github.com/Zequent/zqnt-edge-sdk-go/gen/common/base/proto"
	connectorpb "github.com/Zequent/zqnt-edge-sdk-go/gen/connector/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/protobuf/proto"
)

type fakeConnectorService struct {
	connectorpb.UnimplementedConnectorServiceServer
	lastRequestBase *base.RequestBase
	asset           *asset.AssetProtoDTO
	errorMessage    string
}

func (s *fakeConnectorService) GetAssetBySn(ctx context.Context, req *base.RequestBase) (*connectorpb.ConnectorResponse, error) {
	s.lastRequestBase = req
	if s.errorMessage != "" {
		hasErrors := true
		return &connectorpb.ConnectorResponse{HasErrors: &hasErrors,
			Response: &connectorpb.ConnectorResponse_Error{Error: &base.GlobalErrorMessage{ErrorMessage: s.errorMessage}}}, nil
	}
	return &connectorpb.ConnectorResponse{Response: &connectorpb.ConnectorResponse_Asset{Asset: s.asset}}, nil
}

func (s *fakeConnectorService) RegisterAsset(ctx context.Context, req *connectorpb.ConnectorRegisterAssetRequest) (*connectorpb.ConnectorResponse, error) {
	s.asset = req.GetAsset()
	return &connectorpb.ConnectorResponse{Response: &connectorpb.ConnectorResponse_Asset{Asset: req.GetAsset()}}, nil
}

func dialFake(t *testing.T, svc connectorpb.ConnectorServiceServer) *Client {
	t.Helper()
	lis, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("listen: %v", err)
	}
	server := grpc.NewServer()
	connectorpb.RegisterConnectorServiceServer(server, svc)
	go func() { _ = server.Serve(lis) }()
	t.Cleanup(server.Stop)

	conn, err := grpc.NewClient(lis.Addr().String(), grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		t.Fatalf("dial: %v", err)
	}
	t.Cleanup(func() { _ = conn.Close() })
	return New(conn)
}

func TestGetAssetBySnReturnsTheAssetAndCarriesTheSnThrough(t *testing.T) {
	fake := &fakeConnectorService{asset: &asset.AssetProtoDTO{Sn: proto.String("dock-1")}}
	client := dialFake(t, fake)

	got, err := client.GetAssetBySn(context.Background(), "dock-1")
	if err != nil {
		t.Fatalf("GetAssetBySn: %v", err)
	}
	if got.GetSn() != "dock-1" {
		t.Fatalf("expected sn=dock-1, got %q", got.GetSn())
	}
	if fake.lastRequestBase.GetSn() != "dock-1" {
		t.Fatalf("expected the request to carry sn=dock-1, got %q", fake.lastRequestBase.GetSn())
	}
}

func TestGetAssetBySnReturnsAnErrorWhenTheServerReportsOne(t *testing.T) {
	fake := &fakeConnectorService{errorMessage: "asset not found"}
	client := dialFake(t, fake)

	_, err := client.GetAssetBySn(context.Background(), "unknown-sn")
	if err == nil {
		t.Fatalf("expected an error, got none")
	}
}

func TestRegisterAssetSendsTheAssetThrough(t *testing.T) {
	fake := &fakeConnectorService{}
	client := dialFake(t, fake)

	got, err := client.RegisterAsset(context.Background(), &asset.AssetProtoDTO{Sn: proto.String("aircraft-1")})
	if err != nil {
		t.Fatalf("RegisterAsset: %v", err)
	}
	if got.GetSn() != "aircraft-1" {
		t.Fatalf("expected sn=aircraft-1, got %q", got.GetSn())
	}
}
