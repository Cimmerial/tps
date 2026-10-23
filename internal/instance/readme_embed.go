package instance

import _ "embed"

//go:embed README.md
var embeddedReadme []byte

//go:embed SPECS.md
var embeddedSpecs []byte

// EmbeddedReadme returns the canonical TPS.md source bundled in release binaries.
func EmbeddedReadme() []byte {
	return embeddedReadme
}

// EmbeddedSpecs returns the canonical TPS_SPECS.md source bundled in release binaries.
func EmbeddedSpecs() []byte {
	return embeddedSpecs
}
