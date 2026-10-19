//go:build !darwin

package jobs

func spawnCloseWatcherScript(projectID string, job Job) error { return nil }
