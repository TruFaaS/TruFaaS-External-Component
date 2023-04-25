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

const TrustVerificationHeader = "x-trust-verification"
const MACHeader = "x-MAC"
const ServerPublicKeyHeader = "x-server-public-key"
const ClientPublicKeyHeader = "x-client-public-key"
