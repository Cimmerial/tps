//go:build !windows

package update

const FlagApplyStaged = "--tps-apply-staged"

func ApplyStagedWindowsBinaryIfNeeded() (needReload bool, err error) {
	return false, nil
}

func ScheduleWindowsStagedApplyAfterExit() {}
