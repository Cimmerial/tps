//go:build !windows

package update

const FlagApplyStaged = "--tps-apply-staged"

func ApplyStagedWindowsBinaryIfNeeded() error {
	return nil
}
