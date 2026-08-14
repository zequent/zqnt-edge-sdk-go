package edgeadapter

import (
	"testing"

	base "github.com/Zequent/zqnt-edge-sdk-go/gen/common/base/proto"
)

func TestSuccessHasNoErrors(t *testing.T) {
	resp := Success("tid-1")
	if resp.GetHasErrors() {
		t.Fatalf("expected HasErrors=false, got true")
	}
	if resp.GetMeta().GetTid() != "tid-1" {
		t.Fatalf("expected tid to round-trip, got %q", resp.GetMeta().GetTid())
	}
	if resp.GetEmpty() == nil {
		t.Fatalf("expected an empty payload, got %v", resp.GetResponse())
	}
}

func TestErrorReportsTheMessageAndDefaultsToServiceErrorCode(t *testing.T) {
	resp := Error("tid-2", "device unreachable")
	if !resp.GetHasErrors() {
		t.Fatalf("expected HasErrors=true, got false")
	}
	if resp.GetError().GetErrorMessage() != "device unreachable" {
		t.Fatalf("expected the message to round-trip, got %q", resp.GetError().GetErrorMessage())
	}
	if resp.GetError().GetErrorCode() != base.ErrorCode_ERROR_CODE_SERVICE {
		t.Fatalf("expected ERROR_CODE_SERVICE, got %v", resp.GetError().GetErrorCode())
	}
}

func TestErrorWithCodeUsesTheSuppliedCode(t *testing.T) {
	resp := ErrorWithCode("tid-3", "bad parameter", base.ErrorCode_ERROR_CODE_CLIENT)
	if resp.GetError().GetErrorCode() != base.ErrorCode_ERROR_CODE_CLIENT {
		t.Fatalf("expected ERROR_CODE_CLIENT, got %v", resp.GetError().GetErrorCode())
	}
}
