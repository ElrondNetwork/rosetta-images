package version

const (
	// RosettaVersion is the version of the Rosetta interface the implementation adheres to.
	RosettaVersion = "v1.4.12"

	// RosettaMiddlewareVersion is the version of this package (application)
	RosettaMiddlewareVersion = "v0.2.5"

	// NodeVersion is the canonical version of the node runtime
	// TODO: We should fetch this from node/status.
	NodeVersion = "rc/2022-july"
)
