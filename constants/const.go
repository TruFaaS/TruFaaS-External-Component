package constants

type FaaSPlatform int

// Enum for available platforms
const (
	Fission   FaaSPlatform = iota + 1 // EnumIndex = 1
	OpenWhisk                         // EnumIndex = 2
	OpenFaaS                          // EnumIndex = 3
)

const ContentTypeJSON = "application/json"
const TreeStoreFileName = "tree.gob"

// headers
const (
	TrustVerificationHeader          = "x-trufaas-trust-verification"
	MACHeader                        = "x-trufaas-mac"
	ExternalComponentPublicKeyHeader = "x-trufaas-public-key"
	InvokerPublicKeyHeader           = "x-invoker-public-key"
)
