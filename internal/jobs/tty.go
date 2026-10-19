package jobs

import (
	"fmt"
	"strings"
)

func normalizeTTY(s string) string {
	s = strings.TrimSpace(s)
	return strings.TrimPrefix(s, "/dev/")
}

func JobWindowMarker(projectID string, jobID int) string {
	return projectID + "-" + fmt.Sprint(jobID) + ".command"
}
