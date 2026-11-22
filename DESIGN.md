# DESIGN

- Module `github.com/e3tu/appwrite-base`; project folder `appwrite-go-kit/`; Go 1.26.11 (floor set by the SDK).
- Library-first: apps import the root package and capability subpackages; `cmd/appwrite-base` is the only binary.
- Self-hosted Appwrite only; no cloud defaults anywhere.
- Capabilities return the real Appwrite SDK objects; no wrapper types of our own.
- Own logic is limited to: config wiring, session extraction, `RequireUser`, public URL building.
- `Config` carries both `Endpoint` (internal, used by the SDK) and `PublicEndpoint` (used for browser-facing URLs); never merged.
- Capabilities are optional and independent via pointer blocks `Auth`, `Storage`, `DB`, `Functions`; a nil block means disabled.
- Disabled capability: accessor returns `(nil, false)`.
- No env reads inside library packages; only `cmd/appwrite-base/env.go` parses env and builds `Config`.
- No network I/O in constructors; `New` only stores dependencies and returns errors.
- `database` targets TablesDB (not the legacy Databases API); queries use the SDK `query` helpers.
- Gin is the only web framework dependency; no config or logging libraries.
- Latest-only support: one Appwrite line per release, no version range. Today that is Appwrite 2.7.x (`2.7.0`-`2.7.4` tested).
- SDK pinned to `github.com/appwrite/sdk-for-go/v7 v7.21.0`, which speaks response format `2.7.2`.
- Dropping an old Appwrite line is expected work: bump the SDK, fix breakage, update `version.go`, move on.
