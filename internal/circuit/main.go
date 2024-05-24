package circuit

import "embed"

//go:embed verification_key.json
var VerificationKey embed.FS

const VerificationKeyFileName = "verification_key.json"
