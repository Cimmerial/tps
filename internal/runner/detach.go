package runner

import (
	"strings"

	"github.com/Cimmerial/tps/internal/config"
)

func ShouldDetach(c config.Command) bool {
	if c.Detach || c.Capture != "" {
		return true
	}
	s := strings.ToLower(c.Shell)
	switch {
	case strings.Contains(s, "vite"),
		strings.Contains(s, "npm run dev"),
		strings.Contains(s, "pnpm dev"),
		strings.Contains(s, "yarn dev"),
		strings.Contains(s, "dev:mobile"),
		strings.Contains(s, "next dev"),
		strings.Contains(s, "webpack serve"):
		return true
	default:
		return false
	}
}
