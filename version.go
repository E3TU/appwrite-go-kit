package base

// Supported self-hosted Appwrite version matrix.
//
// These are informational constants for diagnostics and the verify
// subcommand; they do not gate behaviour.
//
// TODO: confirm both bounds against the deployment. The SDK at the pinned
// version speaks ResponseFormat and offers no helper to negotiate versions.
const (
	// AppwriteMinVersion is the oldest self-hosted release supported.
	AppwriteMinVersion = "1.6.0"

	// AppwriteMaxVersion is the newest self-hosted release tested.
	AppwriteMaxVersion = "1.8.0"

	// ResponseFormat is the X-Appwrite-Response-Format the SDK speaks.
	// Keep in sync with the SDK pin in go.mod.
	ResponseFormat = "1.8.0"
)
