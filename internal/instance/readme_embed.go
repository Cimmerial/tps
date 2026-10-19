package instance

import _ "embed"

//go:embed README.md
var embeddedReadme []byte

// EmbeddedReadme returns the canonical TPS.md source bundled in release binaries.
func EmbeddedReadme() []byte {
	return embeddedReadme
}
