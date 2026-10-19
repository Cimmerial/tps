//go:build !darwin

package ide

// AlreadyThere is only implemented on macOS (frontmost app + VS Code family window state).
func AlreadyThere(_ Choice, _ string) (bool, error) {
	return false, nil
}
