// Package connector is a minimal client for the subset of ConnectorService an edge adapter
// actually needs: registering itself and looking up/updating the asset it represents. It is
// deliberately narrow — customer applications should use the (separate) client-go-sdk's fuller
// Connector client instead.
package connector

import (
	"context"
	"fmt"
	"time"

	asset "github.com/Zequent/zqnt-edge-sdk-go/gen/common/asset/proto"
	base "github.com/Zequent/zqnt-edge-sdk-go/gen/common/base/proto"
	connectorpb "github.com/Zequent/zqnt-edge-sdk-go/gen/connector/proto"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// Client wraps ConnectorServiceClient for the edge-adapter-facing subset of its RPCs, converting
// the wire-level has_errors/error convention into an idiomatic Go error return.
type Client struct {
	grpc connectorpb.ConnectorServiceClient
}

// New wraps an existing gRPC connection (dial it yourself — a long-lived adapter typically keeps
// one connection open for its whole lifetime rather than dialing per call).
func New(conn grpc.ClientConnInterface) *Client {
	return &Client{grpc: connectorpb.NewConnectorServiceClient(conn)}
}

// GetAssetBySn looks up the asset registered under sn. Returns (nil, nil) if not found — not
// found is a normal, expected outcome (e.g. an adapter's first-ever contact with a new device),
// not an error condition.
func (c *Client) GetAssetBySn(ctx context.Context, sn string) (*asset.AssetProtoDTO, error) {
	resp, err := c.grpc.GetAssetBySn(ctx, &base.RequestBase{
		Tid: newTid(), Sn: sn, Timestamp: timestamppb.Now(),
	})
	if err != nil {
		return nil, fmt.Errorf("connector: GetAssetBySn(%s): %w", sn, err)
	}
	if resp.GetHasErrors() {
		return nil, fmt.Errorf("connector: GetAssetBySn(%s): %s", sn, resp.GetError().GetErrorMessage())
	}
	return resp.GetAsset(), nil
}

// RegisterAsset registers a new asset. Fails if one is already registered for the same SN — call
// UpdateAsset instead if you just need to refresh an existing registration.
func (c *Client) RegisterAsset(ctx context.Context, a *asset.AssetProtoDTO) (*asset.AssetProtoDTO, error) {
	resp, err := c.grpc.RegisterAsset(ctx, &connectorpb.ConnectorRegisterAssetRequest{
		Base:  &base.RequestBase{Tid: newTid(), Sn: a.GetSn(), Timestamp: timestamppb.Now()},
		Asset: a,
	})
	if err != nil {
		return nil, fmt.Errorf("connector: RegisterAsset(%s): %w", a.GetSn(), err)
	}
	if resp.GetHasErrors() {
		return nil, fmt.Errorf("connector: RegisterAsset(%s): %s", a.GetSn(), resp.GetError().GetErrorMessage())
	}
	return resp.GetAsset(), nil
}

// UpdateAsset replaces every mutable field of the asset identified by assetID with a. Use
// PatchAsset instead to update only specific fields.
func (c *Client) UpdateAsset(ctx context.Context, assetID string, a *asset.AssetProtoDTO) (*asset.AssetProtoDTO, error) {
	resp, err := c.grpc.UpdateAsset(ctx, &connectorpb.ConnectorUpdateAssetRequest{
		Base:    &base.RequestBase{Tid: newTid(), Sn: a.GetSn(), Timestamp: timestamppb.Now()},
		Asset:   a,
		AssetId: assetID,
	})
	if err != nil {
		return nil, fmt.Errorf("connector: UpdateAsset(%s): %w", assetID, err)
	}
	if resp.GetHasErrors() {
		return nil, fmt.Errorf("connector: UpdateAsset(%s): %s", assetID, resp.GetError().GetErrorMessage())
	}
	return resp.GetAsset(), nil
}

func newTid() string {
	return fmt.Sprintf("edge-go-sdk-%d", time.Now().UnixNano())
}
