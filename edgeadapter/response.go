// Package edgeadapter provides response-building helpers and gRPC server scaffolding for
// implementing zqnt's EdgeAdapterService — the device-facing command API every edge adapter
// (DJI, MAVLink, Betaflight, ...) implements, mirroring the Java/Python Edge SDKs' conventions.
package edgeadapter

import (
	"time"

	base "github.com/Zequent/zqnt-edge-sdk-go/gen/common/base/proto"
	devicecontrol "github.com/Zequent/zqnt-edge-sdk-go/gen/devicecontrol/contracts/proto"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// Now returns the current time as a protobuf Timestamp — every RequestBase/ResponseMeta needs one.
func Now() *timestamppb.Timestamp {
	return timestamppb.New(time.Now())
}

// Success builds an accepted, error-free CommandResponse with no payload — the common case for a
// fire-and-forget command (open cover, start charging, ...).
func Success(tid string) *devicecontrol.CommandResponse {
	hasErrors := false
	return &devicecontrol.CommandResponse{
		HasErrors: &hasErrors,
		Meta:      &base.ResponseMeta{Tid: tid, Timestamp: Now()},
		Response:  &devicecontrol.CommandResponse_Empty{Empty: &emptypb.Empty{}},
	}
}

// Error builds a CommandResponse reporting a command failure — the caller supplies the
// human-readable message; the error code defaults to ERROR_CODE_SERVICE (the adapter itself
// rejected/failed the command), matching what most adapters mean when they return an error here.
func Error(tid, message string) *devicecontrol.CommandResponse {
	return ErrorWithCode(tid, message, base.ErrorCode_ERROR_CODE_SERVICE)
}

// ErrorWithCode is Error with an explicit ErrorCode, for callers that need to distinguish (e.g.)
// a client-supplied parameter error from a device/service failure.
func ErrorWithCode(tid, message string, code base.ErrorCode) *devicecontrol.CommandResponse {
	hasErrors := true
	return &devicecontrol.CommandResponse{
		HasErrors: &hasErrors,
		Meta:      &base.ResponseMeta{Tid: tid, Timestamp: Now()},
		Response: &devicecontrol.CommandResponse_Error{Error: &base.GlobalErrorMessage{
			ErrorCode:    code,
			ErrorMessage: message,
			Timestamp:    Now(),
		}},
	}
}
