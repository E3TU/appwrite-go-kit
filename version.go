package base

// Supported self-hosted Appwrite version matrix.
//
// The library targets one Appwrite line at a time: the latest tested release.
// There is no compatibility range. When a new Appwrite major ships, the SDK
// pin and these constants move together, and the previous line is dropped.
//
// These constants are informational for diagnostics and the verify
// subcommand; they do not gate behaviour.
const (
	// AppwriteMinVersion is the oldest self-hosted release tested.
	AppwriteMinVersion = "2.7.0"

	// AppwriteMaxVersion is the newest self-hosted release tested.
	AppwriteMaxVersion = "2.7.4"

	// SDKVersion is the Appwrite Go SDK module this package is built against.
	SDKVersion = "github.com/appwrite/sdk-for-go/v7 v7.21.0"

	// ResponseFormat is the X-Appwrite-Response-Format the SDK speaks. It
	// must match what the pinned SDK sends, so keep it in sync with
	// SDKVersion.
	ResponseFormat = "2.7.2"
)
