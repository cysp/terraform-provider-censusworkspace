package client

//go:generate go run github.com/ogen-go/ogen/cmd/ogen -target . -package client -clean ../openapi.yml

const (
	// DefaultServerURL is the default URL of the server.
	DefaultServerURL = "https://app.getcensus.com"

	// DefaultUserAgent is the default user agent.
	DefaultUserAgent = "census-management-go/0.1"
)
