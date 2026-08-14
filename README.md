# zqnt-edge-sdk-go

Go SDK for building Zequent edge adapters — the Go counterpart to `edge-java-sdk` /
`edge-python-sdk`. Built locally in this session; **not yet pushed anywhere** — create the
`Zequent/zqnt-edge-sdk-go` GitHub repo, then `git remote add origin <url> && git push -u origin main`.

## Status: first slice, not yet at parity with the Java/Python Edge SDKs

This is a real, tested foundation — not a stub — but it covers a deliberately narrow slice:

**Included:**
- `gen/` — generated protobuf/gRPC Go code for the full canonical protocol (`utils/zqnt-utils/src/main/proto/*.proto`), including the complete 31-RPC `EdgeAdapterService` contract.
- `edgeadapter/` — `NewServer`/`Serve` to stand up an `EdgeAdapterService` gRPC server from your own implementation, plus `Success`/`Error`/`ErrorWithCode` response-building helpers. Standard Go gRPC convention applies: embed `edgepb.UnimplementedEdgeAdapterServiceServer` in your adapter type and only override the commands your hardware actually supports — every other RPC returns `Unimplemented` for free.
- `connector/` — a minimal `ConnectorClient` (`GetAssetBySn`, `RegisterAsset`, `UpdateAsset`) for an adapter to register itself and look up/update its own asset record.

**Not included yet (real gaps, not just "someday"):**
- No `LiveDataService` telemetry-streaming client (Java/Python's `TelemetryPublisher` equivalent).
- No `MissionAutonomyService` client for edge-side use.
- No Redis-based service-discovery helper (`EdgeAdapterConfig`/`EdgeAdapterRuntime` equivalent from the Python SDK) — this package only gives you the gRPC server + a bare Connector client, not the full adapter runtime bootstrap (logging, config loading, discovery registration).
- The generated `gen/` tree is vendored directly into this repo rather than pulled from a shared proto module (no `zqnt-utils-go` exists — see the parent conversation's plan notes) — future proto changes must be regenerated and copied in by hand until that's addressed.

## Generating the protobuf code

```bash
cd utils/zqnt-utils/src/main/proto
protoc -I . --go_out=<repo> --go-grpc_out=<repo> \
  $(for f in *.proto; do pkg=$(grep -oP '(?<=option go_package = ")[^"]+' "$f"); \
    echo "--go_opt=M${f}=github.com/Zequent/zqnt-edge-sdk-go/${pkg} --go-grpc_opt=M${f}=github.com/Zequent/zqnt-edge-sdk-go/${pkg}"; done) \
  *.proto
```
(then move the generated `github.com/Zequent/zqnt-edge-sdk-go/gen` directory up to `./gen` — protoc
places output under the full mapped import path.)

## Requirements

Go 1.24+. `go build ./...` / `go test ./...` both pass as of this commit.
