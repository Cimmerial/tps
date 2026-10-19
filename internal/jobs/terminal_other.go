//go:build !darwin

package jobs

func CloseTerminalTab(ttyDev, windowMarker string) {}
